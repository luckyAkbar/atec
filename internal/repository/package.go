package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/db"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PackageRepo repository for packages
type PackageRepo struct {
	db          *gorm.DB
	cacheKeeper db.CacheKeeperIface
}

// PackageRepoIface interface for PackageRepo
type PackageRepoIface interface {
	Create(ctx context.Context, input CreatePackageInput, txControllers ...*gorm.DB) (*model.Package, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error)
	Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput, txControllers ...*gorm.DB) (*model.Package, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, input SearchPackageInput) ([]model.Package, error)
	FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error)
	FindAllActivePackages(ctx context.Context) ([]model.Package, error)
}

// NewPackageRepo create new package repo instance
func NewPackageRepo(db *gorm.DB, cacheKeeper *db.CacheKeeper) *PackageRepo {
	return &PackageRepo{
		db:          db,
		cacheKeeper: cacheKeeper,
	}
}

// CreatePackageInput input
type CreatePackageInput struct {
	UserID                  uuid.UUID
	PackageName             string
	Questionnaire           model.Questionnaire
	IndicationCategories    model.IndicationCategories
	ImageResultAttributeKey model.ImageResultAttributeKey
}

// Create insert new record of packages to the database
func (r *PackageRepo) Create(ctx context.Context, input CreatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
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

	if err := tx.WithContext(ctx).Create(pack).Error; err != nil {
		return nil, err
	}

	return pack, nil
}

// FindByID find package by id
func (r *PackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	pack := &model.Package{}

	err := r.db.WithContext(ctx).Take(pack, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return pack, nil
	}
}

// UpdatePackageInput input
type UpdatePackageInput struct {
	ActiveStatus *bool

	Questionnaire *model.Questionnaire
	PackageName   string
}

// ToUpdateFields convert the update params to gorm dynamic update fields
func (upi UpdatePackageInput) ToUpdateFields() map[string]interface{} {
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

	return fields
}

// Update update package record by its id
func (r *PackageRepo) Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
	tx := r.db
	if len(txControllers) > 0 {
		tx = txControllers[0]
	}

	pack := &model.Package{}

	err := tx.WithContext(ctx).Model(pack).
		Clauses(clause.Returning{}).Where("id = ?", id).
		Updates(input.ToUpdateFields()).Error

	if err != nil {
		return nil, err
	}

	if input.ActiveStatus != nil {
		if err := r.refreshAllActivePackagesCache(ctx, []uuid.UUID{id}, true); err != nil {
			logrus.WithError(err).Warn("failure to update active packages cache after update (report only)")
		}
	}

	return pack, nil
}

// Delete use soft delete to delete a package record by its id
func (r *PackageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.Package{ID: id}).Error
	if err != nil {
		return err
	}

	err = r.refreshAllActivePackagesCache(ctx, []uuid.UUID{id}, false)
	if err != nil {
		logrus.WithError(err).Warn("failure to update active packages cache after deletion (report only)")
	}

	return nil
}

// SearchPackageInput input to search package. any fields typed with a pointer means it is optional
type SearchPackageInput struct {
	IsActive *bool
	Limit    int
}

func (spi SearchPackageInput) toSearchFields(cursor *gorm.DB) *gorm.DB {
	if spi.IsActive != nil {
		cursor = cursor.Where("is_active = ?", *spi.IsActive)
	}

	if spi.Limit > 0 {
		cursor = cursor.Limit(spi.Limit)
	}

	return cursor
}

// Search search package based on provided parameters
func (r *PackageRepo) Search(ctx context.Context, input SearchPackageInput) ([]model.Package, error) {
	packages := []model.Package{}

	conn := r.db.WithContext(ctx)
	cursor := input.toSearchFields(conn)

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

	packages, err := r.Search(ctx, SearchPackageInput{
		IsActive: &active,
		Limit:    limit,
	})

	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		err := r.cacheKeeper.SetNil(ctx, string(AllActivePackageCacheKey))
		if err != nil {
			logrus.WithError(err).Warn("failed to set nil cache for all active packages")
		}

		return nil, ErrNotFound
	}

	if err := r.cacheKeeper.SetJSON(ctx, string(AllActivePackageCacheKey), packages, config.CacheExpiryDuration().AllActivePackage); err != nil {
		logrus.WithError(err).Warn("failed to cache all active packages")

		return packages, err
	}

	return packages, nil
}
