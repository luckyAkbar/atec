package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// QuestionnaireUsecase usecase for questionnaire
type QuestionnaireUsecase struct {
	packageRepo repository.PackageRepoIface
	childRepo   repository.ChildRepositoryIface
	resultRepo  repository.ResultRepositoryIface
	font        *truetype.Font
}

// QuestionnaireUsecaseIface interface
type QuestionnaireUsecaseIface interface {
	HandleSubmitQuestionnaire(ctx context.Context, input SubmitQuestionnaireInput) (*SubmitQuestionnaireOutput, error)
	HandleDownloadQuestionnaireResult(ctx context.Context, input DownloadQuestionnaireResultInput) (*DownloadQuestionnaireResultOutput, error)
	HandleSearchQuestionnaireResult(ctx context.Context, input SearchQuestionnaireResultInput) ([]SearchQuestionnaireResultOutput, error)
	HandleGetUserHistory(ctx context.Context, input GetUserHistoryInput) ([]GetUserHistoryOutput, error)
	HandleInitializeATECQuestionnaire(ctx context.Context, input InitializeATECQuestionnaireInput) (*InitializeATECQuestionnaireOutput, error)
}

// NewQuestionnaireUsecase create new QuestionnaireUsecase instance
func NewQuestionnaireUsecase(
	packageRepo *repository.PackageRepo, childRepo *repository.ChildRepository,
	resultRepo *repository.ResultRepository, font *truetype.Font,
) *QuestionnaireUsecase {
	return &QuestionnaireUsecase{
		packageRepo: packageRepo,
		childRepo:   childRepo,
		resultRepo:  resultRepo,
		font:        font,
	}
}

// SubmitQuestionnaireInput input
type SubmitQuestionnaireInput struct {
	PackageID uuid.UUID          `validate:"required" json:"package_id"`
	ChildID   uuid.UUID          `json:"child_id"`
	Answers   model.AnswerDetail `validate:"required" json:"answers"`
}

// validate SubmitQuestionnaireInput struct and also ensure all the questions are answered
func (sqi SubmitQuestionnaireInput) validate() error {
	if err := common.Validator.Struct(sqi); err != nil {
		return err
	}

	return ensureAllQuestionAnswered(model.DefaultATECTemplate.SubTest, sqi.Answers)
}

func ensureAllQuestionAnswered(subtest model.SubTest, answers model.AnswerDetail) error {
	for id, subtest := range subtest {
		group, ok := answers[id]
		if !ok {
			return fmt.Errorf("group %d %s is missing answers", id+1, subtest.Name)
		}

		if len(group) != subtest.QuestionCount {
			return fmt.Errorf("group %d %s is expecting %d answers, but got %d", id+1, subtest.Name, subtest.QuestionCount, len(group))
		}
	}

	return nil
}

func performGrading(questionnaire model.Questionnaire, answers model.AnswerDetail) (*model.ResultDetail, error) {
	resultDetail := model.ResultDetail{}

	for subTestID, checklistGroup := range questionnaire {
		answers, ok := answers[subTestID]
		if !ok {
			return nil, fmt.Errorf("missing answers for subtest id %d %s", subTestID+1, checklistGroup.CustomName)
		}

		totalScore := 0
		possibleAnswers := checklistGroup.Options

		for _, answer := range answers {
			found := false

			for _, opt := range possibleAnswers {
				if opt.ID == answer {
					found = true

					totalScore += opt.Score

					break
				}
			}

			if !found {
				return nil, fmt.Errorf("answer with id: %d is not a valid option", answer)
			}
		}

		resultDetail[subTestID] = model.SubtestGrade{
			Name:  checklistGroup.CustomName,
			Grade: totalScore,
		}
	}

	return &resultDetail, nil
}

// SubmitQuestionnaireOutput output
type SubmitQuestionnaireOutput struct {
	ResultID   uuid.UUID                `json:"result_id"`
	PackageID  uuid.UUID                `json:"package_id"`
	Answers    model.AnswerDetail       `json:"answers"`
	Result     model.ResultDetail       `json:"result"`
	Indication model.IndicationCategory `json:"indication"`
	ChildID    uuid.UUID                `json:"child_id"`
	CreatedBy  uuid.UUID                `json:"created_by"`
	CreatedAt  time.Time                `json:"created_at"`
}

// HandleSubmitQuestionnaire will handle the submission of a questionnaire result.
func (u *QuestionnaireUsecase) HandleSubmitQuestionnaire(ctx context.Context, input SubmitQuestionnaireInput) (*SubmitQuestionnaireOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	pack, err := u.packageRepo.FindByID(ctx, input.PackageID)
	switch err {
	default:
		logger.WithError(err).Error("failed to fetch package detail from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	grade, err := performGrading(pack.Questionnaire, input.Answers)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	// chosen not to save this to database because it not necessary, and it is fairly cheap to compute
	// and also not many places need this. If on the future decide to save this to database, this is the place
	// to start
	indication := pack.IndicationCategories.GetIndicationCategoryByScore(grade.CountTotalScore())

	requester := model.GetUserFromCtx(ctx)

	if input.ChildID == uuid.Nil {
		createInput := repository.CreateResultInput{
			PackageID: input.PackageID,
			Answer:    input.Answers,
			Result:    *grade,
		}

		if requester != nil {
			createInput.CreatedBy = requester.ID
		}

		result, err := u.resultRepo.Create(ctx, createInput)

		if err != nil {
			logger.WithError(err).Error("failed to write questionnaire result to database")

			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		}

		return &SubmitQuestionnaireOutput{
			ResultID:   result.ID,
			PackageID:  input.PackageID,
			Answers:    result.Answer,
			Result:     result.Result,
			Indication: indication,
			CreatedBy:  result.CreatedBy,
			CreatedAt:  result.CreatedAt,
		}, nil
	}

	child, err := u.childRepo.FindByID(ctx, input.ChildID)
	switch err {
	default:
		logger.WithError(err).Error("failed to fetch child detail from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "filling questionnaire for a child requires valid authorization",
		}
	}

	if requester.Role != model.RolesAdmin && requester.ID != child.ParentUserID {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "filling questionnaire for a child must be done by either the parents or admin",
		}
	}

	result, err := u.resultRepo.Create(ctx, repository.CreateResultInput{
		PackageID: input.PackageID,
		Answer:    input.Answers,
		Result:    *grade,
		ChildID:   input.ChildID,
		CreatedBy: requester.ID,
	})

	if err != nil {
		logger.WithError(err).Error("failed to write questionnaire result to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &SubmitQuestionnaireOutput{
		ResultID:   result.ID,
		PackageID:  result.PackageID,
		Answers:    result.Answer,
		Result:     result.Result,
		Indication: indication,
		ChildID:    result.ChildID,
		CreatedBy:  result.CreatedBy,
		CreatedAt:  result.CreatedAt,
	}, nil
}

// SearchQuestionnaireResultInput search questionnaire result
type SearchQuestionnaireResultInput struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Limit     int `validate:"min=1,max=100"`
	Offset    int `validate:"min=0"`
}

func (sqi SearchQuestionnaireResultInput) validate() error {
	return common.Validator.Struct(sqi)
}

// SearchQuestionnaireResultOutput search questionnaire result
type SearchQuestionnaireResultOutput struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Answer    model.AnswerDetail
	Result    model.ResultDetail
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

// HandleSearchQuestionnaireResult search questionnaire results from database based on given search param
// this function will return a list of questionnaire result based on the search param
func (u *QuestionnaireUsecase) HandleSearchQuestionnaireResult(
	ctx context.Context, input SearchQuestionnaireResultInput,
) ([]SearchQuestionnaireResultOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	results, err := u.resultRepo.Search(ctx, repository.SearchResultInput{
		ID:        input.ID,
		PackageID: input.PackageID,
		ChildID:   input.ChildID,
		CreatedBy: input.CreatedBy,
		Limit:     input.Limit,
		Offset:    input.Offset,
	})

	switch err {
	default:
		logger.WithError(err).Error("failed to search questionnaire result from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	output := []SearchQuestionnaireResultOutput{}

	for _, res := range results {
		output = append(output, SearchQuestionnaireResultOutput{
			ID:        res.ID,
			PackageID: res.PackageID,
			ChildID:   res.ChildID,
			CreatedBy: res.CreatedBy,
			Answer:    res.Answer,
			Result:    res.Result,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
			DeletedAt: res.DeletedAt,
		})
	}

	return output, nil
}

// GetUserHistoryInput input
type GetUserHistoryInput struct {
	Limit  int `validate:"min=1,max=100"`
	Offset int `validate:"min=0"`
}

func (guhi GetUserHistoryInput) validate() error {
	return common.Validator.Struct(guhi)
}

// GetUserHistoryOutput output
type GetUserHistoryOutput struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Answer    model.AnswerDetail
	Result    model.ResultDetail
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

// HandleGetUserHistory will return a list of questionnaire result made by the requester
// or questionnaire submitted for the requester's children
func (u *QuestionnaireUsecase) HandleGetUserHistory(ctx context.Context, input GetUserHistoryInput) ([]GetUserHistoryOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(ctx))

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "getting user history requires valid authorization",
		}
	}

	results, err := u.resultRepo.FindAllUserHistory(ctx, repository.FindAllUserHistoryInput{
		UserID: requester.ID,
		Limit:  input.Limit,
		Offset: input.Offset,
	})

	switch err {
	default:
		logger.WithError(err).Error("failed to get user history from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	output := []GetUserHistoryOutput{}

	for _, res := range results {
		output = append(output, GetUserHistoryOutput{
			ID:        res.ID,
			PackageID: res.PackageID,
			ChildID:   res.ChildID,
			CreatedBy: res.CreatedBy,
			Answer:    res.Answer,
			Result:    res.Result,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
			DeletedAt: res.DeletedAt,
		})
	}

	return output, nil
}

// DownloadQuestionnaireResultInput input
type DownloadQuestionnaireResultInput struct {
	ResultID uuid.UUID `validate:"required"`
}

func (dqri DownloadQuestionnaireResultInput) validate() error {
	return common.Validator.Struct(dqri)
}

// DownloadQuestionnaireResultOutput output
type DownloadQuestionnaireResultOutput struct {
	ContentType string
	Buffer      bytes.Buffer
}

// HandleDownloadQuestionnaireResult handler to handle downloading questionnaire result
// as an image
func (u *QuestionnaireUsecase) HandleDownloadQuestionnaireResult(
	ctx context.Context, input DownloadQuestionnaireResultInput,
) (*DownloadQuestionnaireResultOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	result, err := u.resultRepo.FindByID(ctx, input.ResultID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find result from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	if result.CreatedBy != uuid.Nil {
		requester := model.GetUserFromCtx(ctx)

		if requester == nil {
			return nil, UsecaseError{
				ErrType: ErrUnauthorized,
				Message: "accessing this result requires authorization",
			}
		}

		if requester.Role != model.RolesAdmin && requester.ID != result.CreatedBy {
			return nil, UsecaseError{
				ErrType: ErrUnauthorized,
				Message: "only owner and admin can access this result",
			}
		}
	}

	pack, err := u.packageRepo.FindByID(ctx, result.PackageID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find package from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	imgGenerator := newImageGenerator(u.font, imageGenerationOpts{
		Result:                  result.Result,
		TestID:                  result.PackageID,
		indicationCategories:    pack.IndicationCategories,
		imageResultAttributeKey: pack.ImageResultAttributeKey,
	})

	imageResult := imgGenerator.GenerateJPEG()

	return &DownloadQuestionnaireResultOutput{
		ContentType: imageResult.ContentType,
		Buffer:      imageResult.Buffer,
	}, nil
}

type imageGenerator struct {
	Result model.ResultDetail
	TestID uuid.UUID

	indicationCategories    model.IndicationCategories
	imageResultAttributeKey model.ImageResultAttributeKey
	rgba                    *image.RGBA
	ttp                     []string // ttp stands for text to print
	width                   int
	height                  int
	titleDrawer             *font.Drawer
	textDrawer              *font.Drawer
	sampleDrawer            *font.Drawer
	font                    *truetype.Font
	dpi                     float64
	textSize                float64
	titleSize               float64
	spacing                 float64
}

type imageGenerationOpts struct {
	Title                   string
	Result                  model.ResultDetail
	TestID                  uuid.UUID
	indicationCategories    model.IndicationCategories
	imageResultAttributeKey model.ImageResultAttributeKey
}

type imageResult struct {
	ContentType string
	Buffer      bytes.Buffer
}

//nolint:mnd
func newImageGenerator(f *truetype.Font, opts imageGenerationOpts) *imageGenerator {
	dpi := float64(208)
	size := float64(12)
	spacing := float64(1.5)
	titleSize := float64(18)

	initialTitleDrawer := &font.Drawer{
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    titleSize,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}),
	}

	initialTextDrawer := &font.Drawer{
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}

	imgGenerator := &imageGenerator{
		Result:                  opts.Result,
		TestID:                  opts.TestID,
		indicationCategories:    opts.indicationCategories,
		imageResultAttributeKey: opts.imageResultAttributeKey,
		sampleDrawer:            initialTextDrawer,
		spacing:                 spacing,
		font:                    f,
		textSize:                size,
		titleSize:               titleSize,
		dpi:                     dpi,
	}

	imgGenerator.generateTTP()
	imgGenerator.countOptimumImageWidth(initialTextDrawer, initialTitleDrawer)
	imgGenerator.countOptimumImageHeight()

	rgba := image.NewRGBA(image.Rect(0, 0, imgGenerator.width, imgGenerator.height))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.Point{}, draw.Src)
	imgGenerator.rgba = rgba
	imgGenerator.generateTextDrawer()
	imgGenerator.generateTitleDrawer()

	return imgGenerator
}

// GenerateJPEG will generate jpeg image for the test result
//
//nolint:mnd
func (ig *imageGenerator) GenerateJPEG() *imageResult {
	y := 10 + int(math.Ceil(ig.textSize*ig.dpi/72))
	dy := int(math.Ceil(ig.textSize * ig.spacing * ig.dpi / 72))
	ig.textDrawer.Dot = fixed.Point26_6{
		X: (fixed.I(ig.width) - ig.textDrawer.MeasureString(ig.imageResultAttributeKey.Title)) / 2,
		Y: fixed.I(y),
	}

	ty := 10 + int(math.Ceil(ig.titleSize*ig.dpi/72))
	tdy := int(math.Ceil(ig.titleSize * ig.spacing * ig.dpi / 72))

	tx := (fixed.I(ig.width) - ig.titleDrawer.MeasureString(ig.imageResultAttributeKey.Title)) / 2

	ig.titleDrawer.Dot = fixed.Point26_6{
		X: tx,
		Y: fixed.I(ty),
	}

	ig.titleDrawer.DrawString(ig.imageResultAttributeKey.Title)

	y += tdy

	for _, s := range ig.ttp {
		center := (fixed.I(ig.width) - ig.textDrawer.MeasureString(s)) / 2
		ig.textDrawer.Dot = fixed.P(center.Ceil(), y)
		ig.textDrawer.DrawString(s)

		y += dy
	}

	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, ig.rgba, nil); err != nil {
		logrus.WithError(err).Error("failed to encode image")
	}

	return &imageResult{
		ContentType: "image/jpeg",
		Buffer:      imgBuf,
	}
}

const optimumTextLength = 65

func (ig *imageGenerator) generateTTP() {
	total := 0

	for _, r := range ig.Result {
		ig.appendTTP(fmt.Sprintf("%s: %d", r.Name, r.Grade))
		total += r.Grade
	}

	indication := ig.indicationCategories.GetIndicationCategoryByScore(total)

	ig.appendTTP(ig.imageResultAttributeKey.Total + " : " + strconv.Itoa(total))
	ig.appendTTP(ig.imageResultAttributeKey.Indication + " : " + indication.Detail)
	ig.appendTTP(ig.imageResultAttributeKey.ResultID + " : " + ig.TestID.String())
	ig.appendTTP(ig.imageResultAttributeKey.SubmittedAt + " : " + time.Now().Format(time.DateTime))
}

func (ig *imageGenerator) appendTTP(s string) {
	arrStr := ig.ensureSafeLongText(s, ig.sampleDrawer)
	ig.ttp = append(ig.ttp, arrStr...)
}

// maxImageWidth are made to limit the maximum image output width
const maxImageWitdh = 1080

//nolint:mnd
func (ig *imageGenerator) countOptimumImageWidth(initialTextDrawer, initialTitleDrawer *font.Drawer) {
	maxWidth := initialTitleDrawer.MeasureString(ig.imageResultAttributeKey.Title)

	for _, t := range ig.ttp {
		ms := initialTextDrawer.MeasureString(t)
		if ms > maxWidth {
			maxWidth = ms
		}
	}

	if maxWidth >= maxImageWitdh {
		ig.width = maxImageWitdh
	} else {
		ig.width = maxWidth.Ceil() + 5*maxWidth.Ceil()/100
	}
}

// ensureSafeLongText will try to check if writing s will cause text overflow
// or if the text length more than optimumTextLength.
// If text deemed to long, will call wordWrapper to wrap the text and prevent overflow
func (ig *imageGenerator) ensureSafeLongText(s string, drawer *font.Drawer) []string {
	width := drawer.MeasureString(s)
	if width.Ceil() >= maxImageWitdh || len(s) >= optimumTextLength {
		return wordWrapper(s)
	}

	return []string{s}
}

func wordWrapper(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{s}
	}

	result := []string{}
	ss := strings.Fields(s)

	var temp string

	for i, word := range ss {
		appended := false
		temp += word + " "

		if len(temp) >= optimumTextLength {
			result = append(result, temp)
			temp = ""

			appended = true
		}

		// when it's already the last word and the temp is not empty
		// just append it to the result
		if i == len(ss)-1 && !appended {
			result = append(result, temp)
		}
	}

	return result
}

//nolint:mnd
func (ig *imageGenerator) countOptimumImageHeight() {
	y := 10 + int(math.Ceil(ig.textSize*ig.dpi/72))
	tdy := int(math.Ceil(ig.titleSize * ig.spacing * ig.dpi / 72))
	y += tdy

	incrementor := int(math.Ceil(ig.textSize * ig.spacing * ig.dpi / 72))
	y += incrementor * len(ig.ttp)

	ig.height = y
}

func (ig *imageGenerator) generateTextDrawer() {
	ig.textDrawer = &font.Drawer{
		Dst: ig.rgba,
		Src: image.Black,
		Face: truetype.NewFace(ig.font, &truetype.Options{
			Size:    ig.textSize,
			DPI:     ig.dpi,
			Hinting: font.HintingNone,
		}),
	}
}

func (ig *imageGenerator) generateTitleDrawer() {
	ig.titleDrawer = &font.Drawer{
		Dst: ig.rgba,
		Src: image.Black,
		Face: truetype.NewFace(ig.font, &truetype.Options{
			Size:    ig.titleSize,
			DPI:     ig.dpi,
			Hinting: font.HintingFull,
		}),
	}
}

// InitializeATECQuestionnaireInput input
type InitializeATECQuestionnaireInput struct {
	PackageID uuid.UUID
}

func (iaqi InitializeATECQuestionnaireInput) useDefaultQuestionnaire() bool {
	return iaqi.PackageID == uuid.Nil
}

// InitializeATECQuestionnaireOutput output
type InitializeATECQuestionnaireOutput struct {
	ID            uuid.UUID
	Questionnaire model.Questionnaire
	Name          string
}

// HandleInitializeATECQuestionnaire get an atec questionaire based on provided input.PackageID
// if not, default to the oldest active and locked package
func (u *QuestionnaireUsecase) HandleInitializeATECQuestionnaire(ctx context.Context, input InitializeATECQuestionnaireInput) (
	*InitializeATECQuestionnaireOutput, error,
) {
	if input.useDefaultQuestionnaire() {
		return u.getDefaultATECQuestionnaire(ctx)
	}

	pack, err := u.packageRepo.FindByID(ctx, input.PackageID)

	switch err {
	default:
		logrus.WithField("input", helper.Dump(input)).WithError(err).Error("failed to find package by id")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		return &InitializeATECQuestionnaireOutput{
			ID:            pack.ID,
			Questionnaire: pack.Questionnaire,
			Name:          pack.Name,
		}, nil
	}
}

func (u *QuestionnaireUsecase) getDefaultATECQuestionnaire(ctx context.Context) (*InitializeATECQuestionnaireOutput, error) {
	pack, err := u.packageRepo.FindOldestActiveAndLockedPackage(ctx)

	switch err {
	default:
		logrus.WithError(err).Error("failed to get oldest active locked package")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		return &InitializeATECQuestionnaireOutput{
			ID:            pack.ID,
			Questionnaire: pack.Questionnaire,
			Name:          pack.Name,
		}, nil
	}
}
