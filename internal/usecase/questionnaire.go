package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

// QuestionnaireUsecase usecase for questionnaire
type QuestionnaireUsecase struct {
	packageRepo repository.PackageRepoIface
	childRepo   repository.ChildRepositoryIface
	resultRepo  repository.ResultRepositoryIface
}

// QuestionnaireUsecaseIface interface
type QuestionnaireUsecaseIface interface {
	HandleSubmitQuestionnaire(ctx context.Context, input SubmitQuestionnaireInput) (*SubmitQuestionnaireOutput, error)
}

// NewQuestionnaireUsecase create new QuestionnaireUsecase instance
func NewQuestionnaireUsecase(
	packageRepo *repository.PackageRepo, childRepo *repository.ChildRepository,
	resultRepo *repository.ResultRepository,
) *QuestionnaireUsecase {
	return &QuestionnaireUsecase{
		packageRepo: packageRepo,
		childRepo:   childRepo,
		resultRepo:  resultRepo,
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
	ResultID  uuid.UUID          `json:"result_id"`
	PackageID uuid.UUID          `json:"package_id"`
	Answers   model.AnswerDetail `json:"answers"`
	Result    model.ResultDetail `json:"result"`
	ChildID   uuid.UUID          `json:"child_id"`
	CreatedBy uuid.UUID          `json:"created_by"`
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
			ResultID:  result.ID,
			PackageID: input.PackageID,
			Answers:   result.Answer,
			Result:    result.Result,
			CreatedBy: result.CreatedBy,
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
		ResultID:  result.ID,
		PackageID: result.PackageID,
		Answers:   result.Answer,
		Result:    result.Result,
		ChildID:   result.ChildID,
		CreatedBy: result.CreatedBy,
	}, nil
}
