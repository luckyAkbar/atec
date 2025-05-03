package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	mockUsecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuestionnaireUsecase_HandleSubmitQuestionnaire(t *testing.T) {
	ctx := context.Background()
	parentID := uuid.New()
	childID := uuid.New()
	packageID := uuid.New()
	truth := true

	validAnswersZeroed := model.AnswerDetail{
		0: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
		},
		1: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
		},
		2: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
		},
		3: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
			21: 0,
			22: 0,
			23: 0,
			24: 0,
			25: 0,
		},
	}

	answersMissingAGroup := model.AnswerDetail{
		0: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
		}, // missing number 1 group
		2: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
		},
		3: {
			1:  0,
			2:  0,
			3:  3,
			4:  3,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
			21: 0,
			22: 0,
			23: 0,
			24: 0,
			25: 0,
		},
	}

	answersHasMoreQuestions := model.AnswerDetail{
		0: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0, // here is the extra
		},
		1: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
		},
		2: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
		},
		3: {
			1:  0,
			2:  0,
			3:  3,
			4:  3,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
			21: 0,
			22: 0,
			23: 0,
			24: 0,
			25: 0,
		},
	}

	answersValueExceedPossibleValues := model.AnswerDetail{
		0: {
			1:  10, // this exceed max possible values
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
		},
		1: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
		},
		2: {
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
		},
		3: {
			1:  0,
			2:  0,
			3:  3,
			4:  3,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
			11: 0,
			12: 0,
			13: 0,
			14: 0,
			15: 0,
			16: 0,
			17: 0,
			18: 0,
			19: 0,
			20: 0,
			21: 0,
			22: 0,
			23: 0,
			24: 0,
			25: 0,
		},
	}

	questionnaire := model.Questionnaire{
		0: {
			CustomName: "Kemampuan Bicara/Berbahasa",
			Questions: []string{
				"Mengetahui namanya sendiri",
				"Berespon pada \"Tidak\" atau \"Stop\"",
				"Dapat mengikuti perintah",
				"Dapat menggunakan 1 kata (Tidak!, Makan, Air, dll)",
				"Dapat menggunakan 2 kata sekaligus bersamaan (Tidak mau!, Pergi pulang, dll)",
				"Dapat menggunakan 3 kata sekaligus bersamaan (Mau minum susu, dll)",
				"Mengetahui 10 kata atau lebih",
				"Dapat membuat kalimat yang berisi 4 kata atau lebih",
				"Mampu menjelaskan apa yang dia inginkan",
				"Mampu menanyakan pertanyaan yang bermakna",
				"Isi pembicaraan cenderung relevan/bermakna",
				"Sering menggunakan kalimat-kalimat yang berurutan",
				"Bisa mengikuti pembicaraan dengan cukup baik",
				"Memiliki kemampuan bicara/berbahasa yang sesuai dengan seusianya",
			},
			Options: []model.AnswerOption{
				{ID: 2, Description: "Tidak Benar", Score: 2},
				{ID: 1, Description: "Agak Benar", Score: 1},
				{ID: 0, Description: "Sangat Benar", Score: 0},
			},
		},
		1: {
			CustomName: "Kemampuan Bersosialisasi",
			Questions: []string{
				"Terlihat seperti berada dalam \"tempurung\" – Anda tidak bisa menjangkau dia",
				"Mengabaikan orang lain",
				"Ketika dipanggil, hanya sedikit atau malah tidak memperhatikan",
				"Tidak kooperatif dan menolak",
				"Tidak ada kontak mata",
				"Lebih suka menyendiri",
				"Tidak menunjukkan rasa kasih sayang",
				"Tidak mampu menyapa orang tua",
				"Menghindari kontak dengan orang lain",
				"Tidak mampu menirukan orang lain",
				"Tidak suka dipegang atau dipeluk",
				"Tidak mau berbagi atau menunjukkan",
				"Tidak bisa melambaikan tangan \"Da..Dahh\"",
				"Sering tidak setuju / menolak (not compliant)",
				"Tantrum, marah-marah",
				"Tidak mempunyai teman",
				"Jarang tersenyum",
				"Tidak peka terhadap perasaan orang lain",
				"Acuh tak acuh ketika disukai orang lain",
				"Acuh tak acuh ketika ditinggal pergi oleh orang tuanya",
			},
			Options: []model.AnswerOption{
				{ID: 0, Description: "Tidak Cocok", Score: 0},
				{ID: 1, Description: "Agak Cocok", Score: 1},
				{ID: 2, Description: "Sangat Cocok", Score: 2},
			},
		},
		2: {
			CustomName: "Kesadaran Sensori/Kognitif",
			Questions: []string{
				"Merespon saat dipanggil namanya",
				"Merespon saat dipuji",
				"Melihat pada orang dan binatang",
				"Melihat pada gambar (dan TV)",
				"Menggambar, mewarnai dan melakukan kesenian",
				"Bermain dengan mainannya secara sesuai",
				"Menggunakan ekspresi wajah yang sesuai",
				"Memahami cerita yang ditayangkan di TV",
				"Memahami penjelasan",
				"Sadar akan lingkungannya",
				"Sadar akan bahaya",
				"Mampu berimajinasi",
				"Memulai aktivitas",
				"Mampu berpakaian sendiri",
				"Memiliki rasa penasaran dan ketertarikan",
				"Suka tantangan, senang mengeksplorasi",
				"Tampak selaras, tidak tampak ‘kosong’",
				"Mampu mengikuti pandangan ke arah semua orang memandang.",
			},
			Options: []model.AnswerOption{
				{ID: 0, Description: "Sangat Cocok", Score: 0},
				{ID: 1, Description: "Agak Cocok", Score: 1},
				{ID: 2, Description: "Tidak Cocok", Score: 2},
			},
		},
		3: {
			CustomName: "Kesehatan Umum, Fisik dan Perilaku",
			Questions: []string{
				"Mengompol saat tidur",
				"Mengompol di celana/popok",
				"Buang air besar di celana/popok",
				"Diare",
				"Konstipasi / Sembelit",
				"Gangguan Tidur",
				"Makan terlalu banyak / terlalu sedikit",
				"Pilihan makanan yang diinginkan sangat terbatas (extremely limited diet, picky eater)",
				"Hiperaktif",
				"Letargi, lemah, lesu",
				"Memukul atau melukai diri sendiri",
				"Memukul atau melukai orang lain",
				"Destruktif",
				"Sensitif terhadap suara",
				"Cemas / penuh ketakutan",
				"Tidak senang/ mudah rewel/ menangis",
				"Kejang",
				"Bicara secara obsesif",
				"Kaku terhadap rutinitas",
				"Berteriak / menjerit-jerit",
				"Menuntut hal atau cara yang sama berulang-ulang",
				"Sering gelisah / agitasi",
				"Tidak peka terhadap nyeri",
				"Terfokus atau sulit dialihkan dari objek atau topik tertentu",
				"Gerakan repetitive (stimming, menggoyang-goyangkan bagian badan)",
			},
			Options: []model.AnswerOption{
				{ID: 3, Description: "Sangat Bermasalah", Score: 3},
				{ID: 2, Description: "Cukup Bermasalah", Score: 2},
				{ID: 1, Description: "Sedikit Bermasalah", Score: 1},
				{ID: 0, Description: "Tidak bermasalah", Score: 0},
			},
		},
	}

	inactivePackage := &model.Package{
		ID:       packageID,
		IsActive: false,
	}

	selectedLockedPackage := &model.Package{
		ID:            packageID,
		Questionnaire: questionnaire,
		IsActive:      true,
		IsLocked:      true,
	}

	selectedUnlockedPackage := &model.Package{
		ID:            packageID,
		Questionnaire: questionnaire,
		IsActive:      true,
		IsLocked:      false,
	}

	zeroedResultDetail := model.ResultDetail{
		0: model.SubtestGrade{
			Name:  selectedLockedPackage.Questionnaire[0].CustomName,
			Grade: 0,
		},
		1: model.SubtestGrade{
			Name:  selectedLockedPackage.Questionnaire[1].CustomName,
			Grade: 0,
		},
		2: model.SubtestGrade{
			Name:  selectedLockedPackage.Questionnaire[2].CustomName,
			Grade: 0,
		},
		3: model.SubtestGrade{
			Name:  selectedLockedPackage.Questionnaire[3].CustomName,
			Grade: 0,
		},
	}

	resultZeroed := &model.Result{
		ID:     uuid.New(),
		Result: zeroedResultDetail,
	}

	resultZeroedSubmittedByParent := &model.Result{
		ID:        uuid.New(),
		Result:    zeroedResultDetail,
		CreatedBy: parentID,
	}

	parent := model.AuthUser{
		ID:   parentID,
		Role: model.RolesParent,
	}

	therapist := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesTherapist,
	}
	therapistCtx := model.SetUserToCtx(ctx, therapist)

	resultZeroedSubmittedByTherapist := &model.Result{
		ID:        uuid.New(),
		Result:    zeroedResultDetail,
		CreatedBy: therapist.ID,
	}

	child := &model.Child{
		ID:           childID,
		ParentUserID: parentID,
	}

	parentCtx := model.SetUserToCtx(ctx, parent)

	randomParent := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesParent,
	}
	randomParentCtx := model.SetUserToCtx(ctx, randomParent)

	mockPackageRepo := mockUsecase.NewPackageRepo(t)
	mockResultRepo := mockUsecase.NewResultRepository(t)
	mockChildRepo := mockUsecase.NewChildRepository(t)

	uc := usecase.NewQuestionnaireUsecase(mockPackageRepo, mockChildRepo, mockResultRepo, nil)

	testCases := []struct {
		name                 string
		input                usecase.SubmitQuestionnaireInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.SubmitQuestionnaireOutput
		expectedFunctionCall func()
	}{
		{
			name: "invalid input: missing package ID",
			input: usecase.SubmitQuestionnaireInput{
				ChildID: childID,
				Answers: validAnswersZeroed,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid input: missing answers",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid input: answers missing a whole group",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   answersMissingAGroup,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid input: answers contain extra questions",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   answersHasMoreQuestions,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository returning an unexpected error on find by id",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "target package not found",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "unable to submit to inactive package",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(inactivePackage, nil).Once()
			},
		},
		{
			name: "submitted answer value exceed possible value",
			input: usecase.SubmitQuestionnaireInput{
				ChildID:   childID,
				PackageID: packageID,
				Answers:   answersValueExceedPossibleValues,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(selectedLockedPackage, nil).Once()
			},
		},
		{
			name: "success submitting with no defined child id from unknown user",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:     ctx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroed.ID,
				Result:     zeroedResultDetail,
				Indication: selectedLockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(selectedLockedPackage, nil).Once()
				mockResultRepo.EXPECT().Create(ctx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
				}).Return(resultZeroed, nil).Once()
			},
		},
		{
			name: "failed to create result",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(selectedLockedPackage, nil).Once()
				mockResultRepo.EXPECT().Create(ctx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "success submitting with no defined child id from a parent",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				Result:     zeroedResultDetail,
				Indication: selectedLockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parent.ID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()
			},
		},
		{
			name: "submitting to an unlocked package should trigger the locking mechanism",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				Result:     zeroedResultDetail,
				Indication: selectedUnlockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parentID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedUnlockedPackage, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()

				ctxWithCancel := context.WithoutCancel(parentCtx)
				mockPackageRepo.EXPECT().Update(ctxWithCancel, packageID, usecase.RepoUpdatePackageInput{
					LockStatus: &truth,
				}).Return(&model.Package{}, nil).Once()
			},
		},
		{
			name: "failure when trying to lock package should not affecting the result",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				Result:     zeroedResultDetail,
				Indication: selectedUnlockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parentID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedUnlockedPackage, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()

				ctxWithCancel := context.WithoutCancel(parentCtx)
				mockPackageRepo.EXPECT().Update(ctxWithCancel, packageID, usecase.RepoUpdatePackageInput{
					LockStatus: &truth,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "submitting to a child must be done by logged in user",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(selectedLockedPackage, nil).Once()
			},
		},
		{
			name: "repository failed to find child",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:         parentCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "no child found on repository",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:         parentCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "other than parents / therapist should not be able to submit to other child",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:         randomParentCtx,
			wantErr:     true,
			expectedErr: usecase.ErrForbidden,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(randomParentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(randomParentCtx, childID).Return(child, nil).Once()
			},
		},
		{
			name: "parent success submiting to their child",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				PackageID:  packageID,
				Answers:    validAnswersZeroed,
				Result:     zeroedResultDetail,
				Indication: selectedLockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parentID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
					ChildID:   childID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()
			},
		},
		{
			name: "therapist should be able to submit to any child",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:     therapistCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByTherapist.ID,
				PackageID:  packageID,
				Answers:    validAnswersZeroed,
				Result:     zeroedResultDetail,
				Indication: selectedLockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  therapist.ID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(therapistCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(therapistCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Create(therapistCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: therapist.ID,
					ChildID:   childID,
				}).Return(resultZeroedSubmittedByTherapist, nil).Once()
			},
		},
		{
			name: "repository returning failure when creating result",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:         parentCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedLockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parent.ID,
					ChildID:   childID,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "submit to a child using unlocked package should also trigger locking mechanism",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				PackageID:  packageID,
				Answers:    validAnswersZeroed,
				Result:     zeroedResultDetail,
				Indication: selectedUnlockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parentID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedUnlockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
					ChildID:   childID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()
				ctxWithoutCancel := context.WithoutCancel(parentCtx)
				mockPackageRepo.EXPECT().Update(ctxWithoutCancel, packageID, usecase.RepoUpdatePackageInput{
					LockStatus: &truth,
				}).Return(&model.Package{}, nil).Once()
			},
		},
		{
			name: "submit to a child using unlocked package should not affected by failure when trying to lock the package",
			input: usecase.SubmitQuestionnaireInput{
				PackageID: packageID,
				Answers:   validAnswersZeroed,
				ChildID:   childID,
			},
			ctx:     parentCtx,
			wantErr: false,
			expectedOutput: &usecase.SubmitQuestionnaireOutput{
				ResultID:   resultZeroedSubmittedByParent.ID,
				PackageID:  packageID,
				Answers:    validAnswersZeroed,
				Result:     zeroedResultDetail,
				Indication: selectedUnlockedPackage.IndicationCategories.GetIndicationCategoryByScore(0),
				CreatedBy:  parentID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(parentCtx, packageID).Return(selectedUnlockedPackage, nil).Once()
				mockChildRepo.EXPECT().FindByID(parentCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Create(parentCtx, usecase.RepoCreateResultInput{
					PackageID: packageID,
					Answer:    validAnswersZeroed,
					Result:    zeroedResultDetail,
					CreatedBy: parentID,
					ChildID:   childID,
				}).Return(resultZeroedSubmittedByParent, nil).Once()
				ctxWithoutCancel := context.WithoutCancel(parentCtx)
				mockPackageRepo.EXPECT().Update(ctxWithoutCancel, packageID, usecase.RepoUpdatePackageInput{
					LockStatus: &truth,
				}).Return(nil, assert.AnError).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleSubmitQuestionnaire(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)

				// here are the important things to be checked
				assert.Equal(t, res.Result, tc.expectedOutput.Result)
				assert.Equal(t, res.ResultID, tc.expectedOutput.ResultID)
				assert.Equal(t, res.CreatedBy, tc.expectedOutput.CreatedBy)
				assert.Equal(t, res.Indication, tc.expectedOutput.Indication)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestQuestionnaireUsecase_HandleDownloadQuestionnaireResult(t *testing.T) {
	ctx := context.Background()
	resultID := uuid.New()

	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesParent,
	}
	userCtx := model.SetUserToCtx(ctx, user)

	randomUser := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesParent,
	}
	randomUserCtx := model.SetUserToCtx(ctx, randomUser)

	therapistUser := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesTherapist,
	}
	therapistCtx := model.SetUserToCtx(ctx, therapistUser)

	mockPackageRepo := mockUsecase.NewPackageRepo(t)
	mockResultRepo := mockUsecase.NewResultRepository(t)

	fontBytes, err := os.ReadFile("../../assets/font.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	uc := usecase.NewQuestionnaireUsecase(mockPackageRepo, nil, mockResultRepo, font)

	pack := &model.Package{
		ID:                      uuid.New(),
		IndicationCategories:    validIndicationCategories,
		ImageResultAttributeKey: validImageResultAttributeKey,
	}

	resultWithOwner := &model.Result{
		ID:        resultID,
		CreatedBy: user.ID,
		PackageID: pack.ID,
	}

	resultWithoutOwner := &model.Result{
		ID:        resultID,
		PackageID: pack.ID,
	}

	testCases := []struct {
		name                 string
		input                usecase.DownloadQuestionnaireResultInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:        "missing result id",
			input:       usecase.DownloadQuestionnaireResultInput{},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository failed to find result",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "no result found",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "result without owner can be downloaded by anyone",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:     ctx,
			wantErr: false,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(resultWithoutOwner, nil).Once()
				mockPackageRepo.EXPECT().FindByID(ctx, resultWithoutOwner.PackageID).Return(pack, nil).Once()
			},
		},
		{
			name: "simulate when package repo returning error when finding package",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(resultWithoutOwner, nil).Once()
				mockPackageRepo.EXPECT().FindByID(ctx, resultWithoutOwner.PackageID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "simulate when package not found on repository",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(resultWithoutOwner, nil).Once()
				mockPackageRepo.EXPECT().FindByID(ctx, resultWithoutOwner.PackageID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "result with owner can only be downloaded by the owner or the therapist",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(ctx, resultID).Return(resultWithOwner, nil).Once()
			},
		},
		{
			name: "result with owner should not be downloaded by other non therapist user",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:         randomUserCtx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(randomUserCtx, resultID).Return(resultWithOwner, nil).Once()
			},
		},
		{
			name: "result owner downloading their own result",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:     userCtx,
			wantErr: false,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(userCtx, resultID).Return(resultWithOwner, nil).Once()
				mockPackageRepo.EXPECT().FindByID(userCtx, resultWithOwner.PackageID).Return(pack, nil).Once()
			},
		},
		{
			name: "therapist should be able to download any result",
			input: usecase.DownloadQuestionnaireResultInput{
				ResultID: resultID,
			},
			ctx:     therapistCtx,
			wantErr: false,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindByID(therapistCtx, resultID).Return(resultWithOwner, nil).Once()
				mockPackageRepo.EXPECT().FindByID(therapistCtx, resultWithOwner.PackageID).Return(pack, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			_, err := uc.HandleDownloadQuestionnaireResult(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestQuestionnaireUsecae_HandleSearchQuestionnaireResult(t *testing.T) {
	therapist := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesTherapist,
	}

	ctx := model.SetUserToCtx(context.Background(), therapist)

	mockResultRepo := mockUsecase.NewResultRepository(t)

	uc := usecase.NewQuestionnaireUsecase(nil, nil, mockResultRepo, nil)

	validInput := usecase.SearchQuestionnaireResultInput{
		Limit:     10,
		Offset:    0,
		ID:        uuid.New(),
		PackageID: uuid.New(),
		ChildID:   uuid.New(),
		CreatedBy: uuid.New(),
	}

	expectedResultLen := 10

	testCases := []struct {
		name                 string
		input                usecase.SearchQuestionnaireResultInput
		wantErr              bool
		ctx                  context.Context
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name:        "should only allow authorized request",
			input:       usecase.SearchQuestionnaireResultInput{},
			wantErr:     true,
			ctx:         context.Background(),
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "should only allow therapist - parent",
			input:       usecase.SearchQuestionnaireResultInput{},
			wantErr:     true,
			ctx:         model.SetUserToCtx(context.Background(), model.AuthUser{Role: model.RolesParent}),
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "should only allow therapist - administrator",
			input:       usecase.SearchQuestionnaireResultInput{},
			wantErr:     true,
			ctx:         model.SetUserToCtx(context.Background(), model.AuthUser{Role: model.RolesAdministrator}),
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "limit is missing from input",
			input:       usecase.SearchQuestionnaireResultInput{},
			wantErr:     true,
			ctx:         ctx,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "offset should not be negative",
			input: usecase.SearchQuestionnaireResultInput{
				Limit:  10,
				Offset: -1,
			},
			wantErr:     true,
			ctx:         ctx,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "repository returning unexpected error on search",
			input:       validInput,
			wantErr:     true,
			ctx:         ctx,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().Search(ctx, usecase.RepoSearchResultInput{
					Limit:     validInput.Limit,
					Offset:    validInput.Offset,
					ID:        validInput.ID,
					PackageID: validInput.PackageID,
					ChildID:   validInput.ChildID,
					CreatedBy: validInput.CreatedBy,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:        "repository returning not found error on search",
			input:       validInput,
			wantErr:     true,
			ctx:         ctx,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().Search(ctx, usecase.RepoSearchResultInput{
					Limit:     validInput.Limit,
					Offset:    validInput.Offset,
					ID:        validInput.ID,
					PackageID: validInput.PackageID,
					ChildID:   validInput.ChildID,
					CreatedBy: validInput.CreatedBy,
				}).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:              "ok",
			input:             validInput,
			wantErr:           false,
			ctx:               ctx,
			expectedOutputLen: expectedResultLen,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().Search(ctx, usecase.RepoSearchResultInput{
					Limit:     validInput.Limit,
					Offset:    validInput.Offset,
					ID:        validInput.ID,
					PackageID: validInput.PackageID,
					ChildID:   validInput.ChildID,
					CreatedBy: validInput.CreatedBy,
				}).Return(make([]model.Result, expectedResultLen), nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleSearchQuestionnaireResult(ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Len(t, res, tc.expectedOutputLen)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestQuestionnaireUsecase_HandleGetUserHistory(t *testing.T) {
	userID := uuid.New()
	user := model.AuthUser{
		ID:   userID,
		Role: model.RolesParent,
	}

	ctx := context.Background()
	userCtx := model.SetUserToCtx(ctx, user)

	mockResultRepo := mockUsecase.NewResultRepository(t)

	uc := usecase.NewQuestionnaireUsecase(nil, nil, mockResultRepo, nil)

	expectedOutputLen := 78

	validInput := usecase.GetUserHistoryInput{
		Limit:  10,
		Offset: 0,
	}

	testCases := []struct {
		name                 string
		input                usecase.GetUserHistoryInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name:        "limit is missing from input",
			input:       usecase.GetUserHistoryInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "offset should not be negative",
			input: usecase.GetUserHistoryInput{
				Limit:  10,
				Offset: -1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "missing user in context",
			input:       validInput,
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "repository find all user history returning unexpected error",
			input:       validInput,
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindAllUserHistory(userCtx, usecase.RepoFindAllUserHistoryInput{
					Limit:  validInput.Limit,
					Offset: validInput.Offset,
					UserID: userID,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:        "repository returning not found",
			input:       validInput,
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindAllUserHistory(userCtx, usecase.RepoFindAllUserHistoryInput{
					Limit:  validInput.Limit,
					Offset: validInput.Offset,
					UserID: userID,
				}).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:              "ok",
			input:             validInput,
			ctx:               userCtx,
			wantErr:           false,
			expectedOutputLen: expectedOutputLen,
			expectedFunctionCall: func() {
				mockResultRepo.EXPECT().FindAllUserHistory(userCtx, usecase.RepoFindAllUserHistoryInput{
					Limit:  validInput.Limit,
					Offset: validInput.Offset,
					UserID: userID,
				}).Return(make([]model.Result, expectedOutputLen), nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleGetUserHistory(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Len(t, res, tc.expectedOutputLen)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestQuestionnaireUsecase_HandleInitializeATECQuestionnaire(t *testing.T) {
	ctx := context.Background()

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewQuestionnaireUsecase(mockPackageRepo, nil, nil, nil)

	targetPackageID := uuid.New()
	input := usecase.InitializeATECQuestionnaireInput{
		PackageID: targetPackageID,
	}

	defaultQuestionnaire := &model.Package{
		ID:            uuid.New(),
		Questionnaire: validQuestionnaire,
		Name:          "default",
	}

	targetPackage := &model.Package{
		ID:            targetPackageID,
		Questionnaire: validQuestionnaire,
		Name:          "target",
	}

	testCases := []struct {
		name                 string
		input                usecase.InitializeATECQuestionnaireInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.InitializeATECQuestionnaireOutput
		expectedFunctionCall func()
	}{
		{
			name:    "undefined package id on input will use default questionnaire",
			input:   usecase.InitializeATECQuestionnaireInput{},
			wantErr: false,
			expectedOutput: &usecase.InitializeATECQuestionnaireOutput{
				ID:            defaultQuestionnaire.ID,
				Questionnaire: defaultQuestionnaire.Questionnaire,
				Name:          defaultQuestionnaire.Name,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindOldestActiveAndLockedPackage(ctx).Return(defaultQuestionnaire, nil).Once()
			},
		},
		{
			name:        "undefined package id on input will use default questionnaire: no default questionnaire found",
			input:       usecase.InitializeATECQuestionnaireInput{},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindOldestActiveAndLockedPackage(ctx).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:        "undefined package id on input will use default questionnaire: unexpected error",
			input:       usecase.InitializeATECQuestionnaireInput{},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindOldestActiveAndLockedPackage(ctx).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:        "the target package wasn't found",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, targetPackageID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:        "failed to get the target package",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, targetPackageID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "ok",
			input:   input,
			wantErr: false,
			expectedOutput: &usecase.InitializeATECQuestionnaireOutput{
				ID:            targetPackageID,
				Questionnaire: targetPackage.Questionnaire,
				Name:          targetPackage.Name,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, targetPackageID).Return(targetPackage, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleInitializeATECQuestionnaire(ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, res)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}
