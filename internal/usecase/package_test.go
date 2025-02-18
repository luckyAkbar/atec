package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen,lll
func TestCreatePackageInputValidate(t *testing.T) {
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
	validIndicationCategories := model.IndicationCategories{
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
	validImageResultAttributeKey := model.ImageResultAttributeKey{
		Title:       "Title",
		Total:       "Total",
		Indication:  "Indication",
		ResultID:    "ResultID",
		SubmittedAt: "SubmittedAt",
	}

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
