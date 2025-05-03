package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockUsecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
)

var validQuestionnaire = model.Questionnaire{
	0: {
		CustomName: "Kemampuan Bicara/Berbahasa",
		Options: []model.AnswerOption{
			{ID: 2, Description: "Tidak Benar", Score: 2},
			{ID: 1, Description: "Agak Benar", Score: 1},
			{ID: 0, Description: "Sangat Benar", Score: 0},
		},
		Questions: []string{
			"Mengetahui namanya sendiri",
			`Berespon pada "Tidak" atau "Stop"`,
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
	},
	1: {
		CustomName: "Kemampuan Bersosialisasi",
		Options: []model.AnswerOption{
			{ID: 0, Description: "Tidak Cocok", Score: 0},
			{ID: 1, Description: "Agak Cocok", Score: 1},
			{ID: 2, Description: "Sangat Cocok", Score: 2},
		},
		Questions: []string{
			"Terlihat seperti berada dalam \"tempurung\" - Anda tidak bisa menjangkau dia",
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
			`Tidak bisa melambaikan tangan "Da..Dahh"`,
			"Sering tidak setuju / menolak (not compliant)",
			"Tantrum, marah-marah",
			"Tidak mempunyai teman",
			"Jarang tersenyum",
			"Tidak peka terhadap perasaan orang lain",
			"Acuh tak acuh ketika disukai orang lain",
			"Acuh tak acuh ketika ditinggal pergi oleh orang tuanya",
		},
	},
	2: {
		CustomName: "Kesadaran Sensori/Kognitif",
		Options: []model.AnswerOption{
			{ID: 0, Description: "Sangat Cocok", Score: 0},
			{ID: 1, Description: "Agak Cocok", Score: 1},
			{ID: 2, Description: "Tidak Cocok", Score: 2},
		},
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
	},
	3: {
		CustomName: "Kesehatan Umum, Fisik dan Perilaku",
		Options: []model.AnswerOption{
			{ID: 3, Description: "Sangat Bermasalah", Score: 3},
			{ID: 2, Description: "Cukup Bermasalah", Score: 2},
			{ID: 1, Description: "Sedikit Bermasalah", Score: 1},
			{ID: 0, Description: "Tidak bermasalah", Score: 0},
		},
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
	},
}

var validIndicationCategories = model.IndicationCategories{
	{
		MinimumScore: 0,
		MaximumScore: 30,
		Name:         "gejala ringan",
		Detail:       "Menunjukkan kemungkinan bahwa anak memiliki pola perilaku dan keterampuilan komunikasi yang agak normal dan memiliki peluang tinggi untuk menjalani kehidupan normal dan independen yang menunjukkan gejala ASD minimal.",
	},
	{
		MinimumScore: 31,
		MaximumScore: 50,
		Name:         "gejala sedang",
		Detail:       "Menunjukkan kemungkinan bahwa anak kemungkinan besar akan dapat menjalani kehidupan semi independen tanpa perlu ditempatkan di fasilitas perawatan formal",
	},
	{
		MinimumScore: 51,
		MaximumScore: 179,
		Name:         "gejala berat",
		Detail:       "Menunjukkan kemungkinan bahwa anak jatuh ke persentil ke-90 (sangat autis). Anak mungkin akan membutuhkan perawatan berkelanjutan (mungkin di sebuah institusi), dan mungkin tidak dapat mencapai tingkat kebebasan apapun dari orang lain",
	},
}
var validImageResultAttributeKey = model.ImageResultAttributeKey{
	Title:       "Title",
	Total:       "Total",
	Indication:  "Indication",
	ResultID:    "ResultID",
	SubmittedAt: "SubmittedAt",
}

func TestCreatePackageInputValidate(t *testing.T) {
	t.Run("valid input should pass validation", func(t *testing.T) {
		input := usecase.CreatePackageInput{
			PackageName:             "Valid Package",
			Questionnaire:           validQuestionnaire,
			IndicationCategories:    validIndicationCategories,
			ImageResultAttributeKey: validImageResultAttributeKey,
		}

		err := input.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing package name should return error", func(t *testing.T) {
		input := usecase.CreatePackageInput{
			PackageName:             "",
			Questionnaire:           validQuestionnaire,
			IndicationCategories:    validIndicationCategories,
			ImageResultAttributeKey: validImageResultAttributeKey,
		}

		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid questionnaire should return error", func(t *testing.T) {
		input := usecase.CreatePackageInput{
			PackageName: "Valid Package",
			Questionnaire: model.Questionnaire{
				0: {
					CustomName: "",
					Questions:  []string{"Question 1"},
					Options: []model.AnswerOption{
						{ID: 0, Description: "Option 1", Score: 0},
					},
				},
			},
			IndicationCategories:    validIndicationCategories,
			ImageResultAttributeKey: validImageResultAttributeKey,
		}

		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid indication categories should return error", func(t *testing.T) {
		input := usecase.CreatePackageInput{
			PackageName:   "Valid Package",
			Questionnaire: validQuestionnaire,
			IndicationCategories: model.IndicationCategories{
				{MinimumScore: 0, MaximumScore: 10, Name: "Low", Detail: "Low indication"},
				{MinimumScore: 11, MaximumScore: 20, Name: "", Detail: "Medium indication"}, // Invalid category
				{MinimumScore: 21, MaximumScore: 30, Name: "High", Detail: "High indication"},
			},
			ImageResultAttributeKey: validImageResultAttributeKey,
		}

		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid image result attribute key should return error", func(t *testing.T) {
		input := usecase.CreatePackageInput{
			PackageName:          "Valid Package",
			Questionnaire:        validQuestionnaire,
			IndicationCategories: validIndicationCategories,
			ImageResultAttributeKey: model.ImageResultAttributeKey{
				Title:       "",
				Total:       "Total",
				Indication:  "Indication",
				ResultID:    "ResultID",
				SubmittedAt: "SubmittedAt",
			},
		}

		err := input.Validate()
		assert.Error(t, err)
	})
}

func TestChangeActiveStatusInputValidate(t *testing.T) {
	t.Run("Valid ChangeActiveStatusInput", func(t *testing.T) {
		input := usecase.ChangeActiveStatusInput{
			PackageID:    uuid.New(),
			ActiveStatus: true,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid ChangeActiveStatusInput - Missing PackageID", func(t *testing.T) {
		input := usecase.ChangeActiveStatusInput{
			ActiveStatus: true,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})
}

func TestUpdatePackageInputValidate(t *testing.T) {
	validQuestionnaire := model.Questionnaire{
		0: {
			CustomName: "Kemampuan Bicara/Berbahasa",
			Options: []model.AnswerOption{
				{ID: 2, Description: "Tidak Benar", Score: 2},
				{ID: 1, Description: "Agak Benar", Score: 1},
				{ID: 0, Description: "Sangat Benar", Score: 0},
			},
			Questions: []string{
				"Mengetahui namanya sendiri",
				`Berespon pada "Tidak" atau "Stop"`,
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
		},
		1: {
			CustomName: "Kemampuan Bersosialisasi",
			Options: []model.AnswerOption{
				{ID: 0, Description: "Tidak Cocok", Score: 0},
				{ID: 1, Description: "Agak Cocok", Score: 1},
				{ID: 2, Description: "Sangat Cocok", Score: 2},
			},
			Questions: []string{
				"Terlihat seperti berada dalam \"tempurung\" - Anda tidak bisa menjangkau dia",
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
				`Tidak bisa melambaikan tangan "Da..Dahh"`,
				"Sering tidak setuju / menolak (not compliant)",
				"Tantrum, marah-marah",
				"Tidak mempunyai teman",
				"Jarang tersenyum",
				"Tidak peka terhadap perasaan orang lain",
				"Acuh tak acuh ketika disukai orang lain",
				"Acuh tak acuh ketika ditinggal pergi oleh orang tuanya",
			},
		},
		2: {
			CustomName: "Kesadaran Sensori/Kognitif",
			Options: []model.AnswerOption{
				{ID: 0, Description: "Sangat Cocok", Score: 0},
				{ID: 1, Description: "Agak Cocok", Score: 1},
				{ID: 2, Description: "Tidak Cocok", Score: 2},
			},
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
		},
		3: {
			CustomName: "Kesehatan Umum, Fisik dan Perilaku",
			Options: []model.AnswerOption{
				{ID: 3, Description: "Sangat Bermasalah", Score: 3},
				{ID: 2, Description: "Cukup Bermasalah", Score: 2},
				{ID: 1, Description: "Sedikit Bermasalah", Score: 1},
				{ID: 0, Description: "Tidak bermasalah", Score: 0},
			},
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
		},
	}

	t.Run("valid input should pass validation", func(t *testing.T) {
		input := usecase.UpdatePackageInput{
			PackageID:     uuid.New(),
			PackageName:   "Valid Package",
			Questionnaire: validQuestionnaire,
		}

		err := input.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing package ID should return error", func(t *testing.T) {
		input := usecase.UpdatePackageInput{
			PackageID:     uuid.Nil,
			PackageName:   "Valid Package",
			Questionnaire: validQuestionnaire,
		}

		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("missing package name should return error", func(t *testing.T) {
		input := usecase.UpdatePackageInput{
			PackageID:     uuid.New(),
			PackageName:   "",
			Questionnaire: validQuestionnaire,
		}

		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid questionnaire should return error", func(t *testing.T) {
		input := usecase.UpdatePackageInput{
			PackageID:   uuid.New(),
			PackageName: "Valid Package",
			Questionnaire: model.Questionnaire{
				0: {
					CustomName: "",
					Questions:  []string{"Question 1"},
					Options: []model.AnswerOption{
						{ID: 0, Description: "Option 1", Score: 0},
					},
				},
			},
		}

		err := input.Validate()
		assert.Error(t, err)
	})
}

func TestPackageUsecase_Create(t *testing.T) {
	ctx := context.Background()

	packageID := uuid.New()
	userID := uuid.New()
	user := model.AuthUser{
		ID:   userID,
		Role: model.RolesAdministrator,
	}

	userCtx := model.SetUserToCtx(ctx, user)

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewPackageUsecase(mockPackageRepo)

	validInput := usecase.CreatePackageInput{
		PackageName:             "valid package name",
		Questionnaire:           validQuestionnaire,
		IndicationCategories:    validIndicationCategories,
		ImageResultAttributeKey: validImageResultAttributeKey,
	}

	testCases := []struct {
		name                 string
		input                usecase.CreatePackageInput
		wantErr              bool
		ctx                  context.Context
		expectedErr          error
		expectedOutput       *usecase.CreatePackageOutput
		expectedFunctionCall func()
	}{
		{
			name:        "empty user on ctx",
			input:       usecase.CreatePackageInput{},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "missing required field triggers validation error",
			input:       usecase.CreatePackageInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "repository return error on create",
			input:       validInput,
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().Create(userCtx, usecase.RepoCreatePackageInput{
					UserID:                  userID,
					PackageName:             validInput.PackageName,
					Questionnaire:           validInput.Questionnaire,
					IndicationCategories:    validInput.IndicationCategories,
					ImageResultAttributeKey: validInput.ImageResultAttributeKey,
				}).Return(nil, usecase.ErrInternal).Once()
			},
		},

		{
			name:    "ok",
			input:   validInput,
			ctx:     userCtx,
			wantErr: false,
			expectedOutput: &usecase.CreatePackageOutput{
				ID: packageID,
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().Create(userCtx, usecase.RepoCreatePackageInput{
					UserID:                  userID,
					PackageName:             validInput.PackageName,
					Questionnaire:           validInput.Questionnaire,
					IndicationCategories:    validInput.IndicationCategories,
					ImageResultAttributeKey: validInput.ImageResultAttributeKey,
				}).Return(&model.Package{ID: packageID}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.Create(tc.ctx, tc.input)

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

func TestPackageUsecase_ChangeActiveStatus(t *testing.T) {
	ctx := context.Background()

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewPackageUsecase(mockPackageRepo)

	packageID := uuid.New()
	statusEnabled := true
	statusDisabled := false

	activePackage := &model.Package{
		ID:       packageID,
		IsActive: statusEnabled,
	}

	inactivePackage := &model.Package{
		ID:       packageID,
		IsActive: statusDisabled,
	}

	lockedActivePackage := &model.Package{
		ID:       packageID,
		IsActive: statusEnabled,
		IsLocked: true,
	}

	testCases := []struct {
		name                 string
		input                usecase.ChangeActiveStatusInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.ChangeActiveStatusOutput
		expectedFunctionCall func()
	}{
		{
			name:        "missing package id",
			input:       usecase.ChangeActiveStatusInput{ActiveStatus: statusDisabled},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository returning unexpected error",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusEnabled,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "package was not found on repository",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusEnabled,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "package was already active",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusEnabled,
			},
			wantErr: false,
			expectedOutput: &usecase.ChangeActiveStatusOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(activePackage, nil).Once()
			},
		},
		{
			name: "package still inactive",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusDisabled,
			},
			wantErr: false,
			expectedOutput: &usecase.ChangeActiveStatusOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(inactivePackage, nil).Once()
			},
		},
		{
			name: "unable to change active status because package is locked",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusDisabled,
			},
			wantErr:     true,
			expectedErr: usecase.ErrForbidden,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(lockedActivePackage, nil).Once()
			},
		},
		{
			name: "repository failed to update package",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusEnabled,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(inactivePackage, nil).Once()
				mockPackageRepo.EXPECT().Update(ctx, packageID, usecase.RepoUpdatePackageInput{
					ActiveStatus: &statusEnabled,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.ChangeActiveStatusInput{
				PackageID:    packageID,
				ActiveStatus: statusEnabled,
			},
			wantErr: false,
			expectedOutput: &usecase.ChangeActiveStatusOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(inactivePackage, nil).Once()
				mockPackageRepo.EXPECT().Update(ctx, packageID, usecase.RepoUpdatePackageInput{
					ActiveStatus: &statusEnabled,
				}).Return(&model.Package{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.ChangeActiveStatus(ctx, tc.input)

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

func TestPackageUsecase_Update(t *testing.T) {
	ctx := context.Background()

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewPackageUsecase(mockPackageRepo)

	packageID := uuid.New()
	unlockedPackage := &model.Package{
		ID:       packageID,
		IsLocked: false,
	}
	lockedPackage := &model.Package{
		ID:       packageID,
		IsLocked: true,
	}
	input := usecase.UpdatePackageInput{
		PackageID:     packageID,
		PackageName:   "Valid Package",
		Questionnaire: validQuestionnaire,
	}

	testCases := []struct {
		name                 string
		input                usecase.UpdatePackageInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.UpdatePackageOutput
		expectedFunctionCall func()
	}{
		{
			name:        "missing required field triggers validation error",
			input:       usecase.UpdatePackageInput{},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name:        "repository failed to find package by id",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:        "package not found",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:        "unable to update locked package",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrForbidden,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(lockedPackage, nil).Once()
			},
		},
		{
			name:        "repository failed to perform update package",
			input:       input,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(unlockedPackage, nil).Once()
				mockPackageRepo.EXPECT().Update(ctx, packageID, usecase.RepoUpdatePackageInput{
					PackageName:   input.PackageName,
					Questionnaire: &input.Questionnaire,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "ok",
			input:   input,
			wantErr: false,
			expectedOutput: &usecase.UpdatePackageOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindByID(ctx, packageID).Return(unlockedPackage, nil).Once()
				mockPackageRepo.EXPECT().Update(ctx, packageID, usecase.RepoUpdatePackageInput{
					PackageName:   input.PackageName,
					Questionnaire: &input.Questionnaire,
				}).Return(&model.Package{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.Update(ctx, tc.input)

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

func TestPackageUsecase_Delete(t *testing.T) {
	ctx := context.Background()

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewPackageUsecase(mockPackageRepo)

	packageID := uuid.New()

	testCases := []struct {
		name                 string
		input                uuid.UUID
		wantErr              bool
		expectedErr          error
		expectedOutput       error
		expectedFunctionCall func()
	}{
		{
			name:        "repository failed to delete package",
			input:       packageID,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().Delete(ctx, packageID).Return(assert.AnError).Once()
			},
		},
		{
			name:           "ok",
			input:          packageID,
			wantErr:        false,
			expectedOutput: nil,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().Delete(ctx, packageID).Return(nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			err := uc.Delete(ctx, tc.input)

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

func TestPackageUsecase_FindActiveQuestionnaires(t *testing.T) {
	ctx := context.Background()

	mockPackageRepo := mockUsecase.NewPackageRepo(t)

	uc := usecase.NewPackageUsecase(mockPackageRepo)

	expectedOutputLen := 10

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name:        "repository failed to find active questionnaires",
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindAllActivePackages(ctx).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:        "no active package was found",
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindAllActivePackages(ctx).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:              "ok",
			wantErr:           false,
			expectedOutputLen: expectedOutputLen,
			expectedFunctionCall: func() {
				mockPackageRepo.EXPECT().FindAllActivePackages(ctx).Return(make([]model.Package, expectedOutputLen), nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.FindActiveQuestionnaires(ctx)

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
