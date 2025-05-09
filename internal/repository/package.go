package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/db"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PackageRepo repository for packages
type PackageRepo struct {
	db          *gorm.DB
	cacheKeeper db.CacheKeeperIface
}

// NewPackageRepo create new package repo instance
func NewPackageRepo(db *gorm.DB, cacheKeeper db.CacheKeeperIface) *PackageRepo {
	return &PackageRepo{
		db:          db,
		cacheKeeper: cacheKeeper,
	}
}

// Create insert new record of packages to the database
func (r *PackageRepo) Create(ctx context.Context, input usecase.RepoCreatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
	tx := r.db
	if len(txControllers) > 0 {
		tx = txControllers[0]
	}

	pack := &model.Package{
		CreatedBy:               input.UserID,
		Questionnaire:           input.Questionnaire,
		Name:                    input.PackageName,
		IndicationCategories:    input.IndicationCategories,
		ImageResultAttributeKey: input.ImageResultAttributeKey,
	}

	err := tx.WithContext(ctx).Clauses(clause.Returning{}).Create(pack).Error
	if err != nil {
		return nil, err
	}

	cacheKey := CacheKeyForPackage(*pack)

	packageLock, err := r.cacheKeeper.AcquireLock(cacheKey)
	if err != nil {
		logrus.WithError(err).Warn("failed to acquire lock to cache package, just reporting")

		return pack, nil
	}

	defer func() {
		_, err := packageLock.Unlock()
		if err != nil {
			logrus.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	err = r.cacheKeeper.SetJSON(ctx, cacheKey, pack, config.CacheExpiryDuration().Package)
	if err != nil {
		logrus.WithError(err).Warn("failed to cache package, just reporting")
	}

	return pack, nil
}

// FindByID find package by id
func (r *PackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	cacheKey := CacheKeyForPackage(model.Package{
		ID: id,
	})

	val, mutex, err := r.cacheKeeper.GetOrLock(ctx, cacheKey)
	switch err {
	default:
		return nil, fmt.Errorf("unexpected error when trying to call GetOrLock: %w", err)
	case db.ErrCacheNil:
		return nil, ErrNotFound
	case db.ErrLockWaitingTooLong:
		return nil, ErrTimeout
	case nil:
		break
	}

	pack := &model.Package{}

	if mutex == nil && val != "" {
		if err := json.Unmarshal([]byte(val), pack); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache value: %w", err)
		}

		return pack, nil
	}

	defer func() {
		_, err := mutex.Unlock()
		if err != nil {
			logrus.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	err = r.db.WithContext(ctx).Take(pack, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		err = r.cacheKeeper.SetNil(ctx, cacheKey)
		if err != nil {
			logrus.WithError(err).Warn("failed to set nil cache for package, just reporting")
		}

		return nil, ErrNotFound
	case nil:
		break
	}

	err = r.cacheKeeper.SetJSON(ctx, cacheKey, pack, config.CacheExpiryDuration().Package)
	if err != nil {
		logrus.WithError(err).Warn("failed to cache package, just reporting")
	}

	return pack, nil
}

// updatePackageInputToUpdateFields convert the update params to gorm dynamic update fields
func updatePackageInputToUpdateFields(upi usecase.RepoUpdatePackageInput) map[string]interface{} {
	fields := map[string]interface{}{}

	if upi.ActiveStatus != nil {
		fields["is_active"] = *upi.ActiveStatus
	}

	if upi.PackageName != "" {
		fields["name"] = upi.PackageName
	}

	if upi.Questionnaire != nil {
		fields["questionnaire"] = upi.Questionnaire
	}

	if upi.LockStatus != nil {
		fields["is_locked"] = *upi.LockStatus
	}

	return fields
}

// Update update package record by its id
//
//nolint:lll
func (r *PackageRepo) Update(ctx context.Context, id uuid.UUID, input usecase.RepoUpdatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
	tx := r.db
	if len(txControllers) > 0 {
		tx = txControllers[0]
	}

	cacheKey := CacheKeyForPackage(model.Package{
		ID: id,
	})

	packageLock, err := r.cacheKeeper.AcquireLock(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("unable to update package because failed to acquire the lock: %w", err)
	}

	defer func() {
		_, err := packageLock.Unlock()
		if err != nil {
			logrus.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	// to ensure consistency by deleting the stale data from cache
	if err := r.cacheKeeper.Del(ctx, cacheKey); err != nil {
		return nil, fmt.Errorf("abort package update because failed to delete from cache: %w", err)
	}

	pack := &model.Package{}

	err = tx.WithContext(ctx).Model(pack).
		Clauses(clause.Returning{}).Where("id = ?", id).
		Updates(updatePackageInputToUpdateFields(input)).Error

	if err != nil {
		return nil, err
	}

	err = r.cacheKeeper.SetJSON(ctx, cacheKey, pack, config.CacheExpiryDuration().Package)
	if err != nil {
		logrus.WithError(err).Warn("failed to cache package, just reporting")
	}

	// handle edge case where when activate a package, force refresh the active packages cache
	forceRefresh := false
	if input.ActiveStatus != nil {
		forceRefresh = true
	}

	if err := r.refreshAllActivePackagesCache(ctx, []uuid.UUID{id}, forceRefresh); err != nil {
		logrus.WithError(err).Warn("failure to update active packages cache after update (report only)")
	}

	return pack, nil
}

// Delete use soft delete to delete a package record by its id
func (r *PackageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	packageCacheKey := CacheKeyForPackage(model.Package{ID: id})

	packageLock, err := r.cacheKeeper.AcquireLock(packageCacheKey)
	if err != nil {
		return fmt.Errorf("unable to delete package because failed to acquire the lock: %w", err)
	}

	defer func() {
		_, err := packageLock.Unlock()
		if err != nil {
			logrus.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	forceRefreshActivePackagesCache := false

	cacheVal, err := r.cacheKeeper.Get(ctx, packageCacheKey)
	switch err {
	default:
		return fmt.Errorf("unexpected error when trying to call Get: %w", err)
	case nil:
		pack := &model.Package{}
		if err := json.Unmarshal([]byte(cacheVal), pack); err != nil {
			return fmt.Errorf("failed to unmarshal cache value: %w", err)
		}

		if pack.IsActive {
			forceRefreshActivePackagesCache = true
		}
	case db.ErrCacheKeyNotFound, db.ErrCacheNil:
		break
	}

	if err := r.cacheKeeper.Del(ctx, packageCacheKey); err != nil {
		return fmt.Errorf("abort package deletion because failed to delete from cache: %w", err)
	}

	err = r.db.WithContext(ctx).Delete(&model.Package{ID: id}).Error
	if err != nil {
		return err
	}

	if err := r.cacheKeeper.SetNil(ctx, packageCacheKey); err != nil {
		logrus.WithField("package_id", id).WithError(err).Warn("failed to set nil cache for package, just reporting")
	}

	err = r.refreshAllActivePackagesCache(ctx, []uuid.UUID{id}, forceRefreshActivePackagesCache)
	if err != nil {
		logrus.WithError(err).Warn("failure to update active packages cache after deletion (report only)")
	}

	return nil
}

func searchPackageInputToSearchFields(cursor *gorm.DB, spi usecase.RepoSearchPackageInput) *gorm.DB {
	if spi.IsActive != nil {
		cursor = cursor.Where("is_active = ?", *spi.IsActive)
	}

	if spi.Limit > 0 {
		cursor = cursor.Limit(spi.Limit)
	}

	return cursor
}

// Search search package based on provided parameters
func (r *PackageRepo) Search(ctx context.Context, input usecase.RepoSearchPackageInput) ([]model.Package, error) {
	packages := []model.Package{}

	conn := r.db.WithContext(ctx)
	cursor := searchPackageInputToSearchFields(conn, input)

	if err := cursor.Find(&packages).Error; err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, ErrNotFound
	}

	return packages, nil
}

// FindAllActivePackages get all active packages by wrapping search function within it.
// This was made to ease the caching process. You see that the search function if in itself is cached
// doesn't make any sense because we could potentially lose realtime data.
// But, when implemented this way, we can safely cache the result of this function and only need to
// invalidate the cache when a new package is activated or deactivated.
func (r *PackageRepo) FindAllActivePackages(ctx context.Context) ([]model.Package, error) {
	logger := logrus.WithContext(ctx).WithField("function", "FindAllActivePackages")

	val, mutex, err := r.cacheKeeper.GetOrLock(ctx, string(AllActivePackageCacheKey))
	switch err {
	default:
		logger.WithError(err).Debug("got unexpected error when trying to call GetOrLock")

		return nil, err
	case db.ErrCacheNil:
		return nil, ErrNotFound
	case db.ErrLockWaitingTooLong:
		return nil, ErrTimeout
	case nil:
		break
	}

	packages := []model.Package{}

	if mutex == nil && val != "" {
		if err := json.Unmarshal([]byte(val), &packages); err != nil {
			return nil, err
		}

		return packages, nil
	}

	defer func() {
		_, err := mutex.Unlock()
		if err != nil {
			logger.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	return r.findAndSetAllActivePackagesToCache(ctx)
}

// FindOldestActiveAndLockedPackage get the oldest active and locked package
func (r *PackageRepo) FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error) {
	pack := &model.Package{}

	err := r.db.WithContext(ctx).Where("is_active = ? AND is_locked = ?", true, true).Order("created_at ASC").Take(pack).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return pack, nil
	}
}

func (r *PackageRepo) refreshAllActivePackagesCache(ctx context.Context, updatedPackageIDs []uuid.UUID, force bool) error {
	logger := logrus.WithContext(ctx).WithField("function", "refreshAllActivePackagesCache")
	needRefresh := false

	if force {
		needRefresh = true
	} else {
		activePackageCache, err := r.FindAllActivePackages(ctx)
		if err != nil {
			logger.WithError(err).Warn("failed to find all active packages from cache, thus unable to refresh it (if it is needed)")

			return err
		}

		for _, id := range updatedPackageIDs {
			for _, pack := range activePackageCache {
				if pack.ID == id {
					needRefresh = true

					break
				}
			}
		}
	}

	if !needRefresh {
		return nil
	}

	logger.Debug("refreshing all active packages cache")

	mutex, err := r.cacheKeeper.AcquireLock(string(AllActivePackageCacheKey))
	if err != nil {
		logger.WithError(err).Warn("failed to acquire lock to refresh all active packages cache")

		return err
	}

	defer func() {
		_, err := mutex.Unlock()
		if err != nil {
			logger.WithError(err).Warn("failed to unlock cache mutex")
		}
	}()

	_, err = r.findAndSetAllActivePackagesToCache(ctx)
	if err != nil {
		logger.WithError(err).Warn("failed to refresh all active packages cache")

		return err
	}

	return nil
}

func (r *PackageRepo) findAndSetAllActivePackagesToCache(ctx context.Context) ([]model.Package, error) {
	// limit is basically hardcoded here, because i see no reason to make it configurable
	// if you want to make it configurable, you can change it by adding a new parameter to this function
	limit := 100
	active := true

	packages, err := r.Search(ctx, usecase.RepoSearchPackageInput{
		IsActive: &active,
		Limit:    limit,
	})

	switch err {
	default:
		return nil, err
	case ErrNotFound:
		err := r.cacheKeeper.SetNil(ctx, string(AllActivePackageCacheKey))
		if err != nil {
			logrus.WithError(err).Warn("failed to set nil cache for all active packages")
		}

		return nil, ErrNotFound
	case nil:
		break
	}

	if err := r.cacheKeeper.SetJSON(ctx, string(AllActivePackageCacheKey), packages, config.CacheExpiryDuration().AllActivePackage); err != nil {
		logrus.WithError(err).Warn("failed to cache all active packages")

		return packages, err
	}

	return packages, nil
}
