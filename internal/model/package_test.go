package model_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/luckyAkbar/atec/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestChecklistGroupValidation(t *testing.T) {
	t.Run("empty custom name should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "",
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty custom name, expected error, but got nil")
		}
	})

	t.Run("empty questions should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{},
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty questions, expected error, but got nil")
		}
	})

	t.Run("empty options should return error", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{"just to not trigger the error here"},
			Options:    []model.AnswerOption{},
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("on empty options, expected error, but got nil")
		}
	})

	t.Run("way over options (not defined)", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{"just to not trigger the error here"},
			Options: []model.AnswerOption{
				{
					ID:          0,
					Description: "just to not trigger the error here",
					Score:       0,
				},
				{
					ID:          1,
					Description: "just to not trigger the error here",
					Score:       1,
				}, {
					ID:          2,
					Description: "just to not trigger the error here",
					Score:       3, // should trigger error because 2 is missing
				},
			},
		}

		err := cg.Validate()
		if err == nil {
			t.Errorf("missing score 2 on 3 options, expected error, but got nil")
		}
	})

	t.Run("the least it takes to pass the validation", func(t *testing.T) {
		cg := model.ChecklistGroup{
			CustomName: "just to not trigger the error here",
			Questions:  []string{"just to not trigger the error here"},
			Options: []model.AnswerOption{
				{
					ID:          0,
					Description: "just to not trigger the error here",
					Score:       0,
				},
				{
					ID:          1,
					Description: "just to not trigger the error here",
					Score:       1,
				}, {
					ID:          2,
					Description: "just to not trigger the error here",
					Score:       2,
				},
			},
		}

		err := cg.Validate()
		if err != nil {
			t.Errorf("expecting the least it takes to pass the validation, but got %v", err)
		}
	})
}

func TestQuestionnaireValidation(t *testing.T) {
	t.Run("missing atleast 1 quesionnaire (against default template) grup should trigger an error", func(t *testing.T) {
		questionnaire := model.Questionnaire{
			1: model.ChecklistGroup{
				CustomName: "not yet important",
			},
			2: model.ChecklistGroup{
				CustomName: "not yet important",
			},
		}

		err := questionnaire.Validate()
		if err == nil {
			t.Errorf("missing atleast 1 quesionnaire (against default template) grup should trigger an error, but got nil")
		}
	})

	t.Run("all the questionnaire cheklist group must pass their internal validation", func(t *testing.T) {
		questionnaire := model.Questionnaire{
			0: {
				CustomName: "", // this should not empty
				Questions: []string{
					"Mengetahui namanya sendiri",
				},
				Options: []model.AnswerOption{
					{ID: 0, Description: "Sangat Benar", Score: 0},
				},
			},
			1: {
				CustomName: "Kemampuan Bersosialisasi",
				Questions: []string{
					"Terlihat seperti berada dalam \"tempurung\" – Anda tidak bisa menjangkau dia",
				},
				Options: []model.AnswerOption{
					{ID: 2, Description: "Sangat Cocok", Score: 2},
				},
			},
			2: {
				CustomName: "Kesadaran Sensori/Kognitif",
				Questions: []string{
					"Merespon saat dipanggil namanya",
				},
				Options: []model.AnswerOption{
					{ID: 2, Description: "Tidak Cocok", Score: 2},
				},
			},
			3: {
				CustomName: "Kesehatan Umum, Fisik dan Perilaku",
				Questions: []string{
					"Gerakan repetitive (stimming, menggoyang-goyangkan bagian badan)",
				},
				Options: []model.AnswerOption{
					{ID: 0, Description: "Tidak bermasalah", Score: 0},
				},
			},
		}

		questionnaire[1] = model.ChecklistGroup{} // force it to trigger validation error

		err := questionnaire.Validate()
		if err == nil {
			t.Errorf("expecting error on internally questionnaire validation")
		}
	})

	t.Run("questionnaire group Speech/Language/Communication missing questions", func(t *testing.T) {
		var questionnaire = model.Questionnaire{
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
					// missing question here
				},
				Options: []model.AnswerOption{
					{ID: 3, Description: "Sangat Bermasalah", Score: 3},
					{ID: 2, Description: "Cukup Bermasalah", Score: 2},
					{ID: 1, Description: "Sedikit Bermasalah", Score: 1},
					{ID: 0, Description: "Tidak bermasalah", Score: 0},
				},
			},
		}

		err := questionnaire.Validate()
		if err == nil {
			t.Errorf("each group questions must match the defined template")
		}
	})

	t.Run("missing options on group Speech/Language/Communication must trigger error", func(t *testing.T) {
		var questionnaire = model.Questionnaire{
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
					// missing options here
					{ID: 2, Description: "Cukup Bermasalah", Score: 2},
					{ID: 1, Description: "Sedikit Bermasalah", Score: 1},
					{ID: 0, Description: "Tidak bermasalah", Score: 0},
				},
			},
		}

		err := questionnaire.Validate()
		if err == nil {
			t.Errorf("each group options must match the defined template")
		}
	})

	t.Run("ok", func(t *testing.T) {
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

		err := questionnaire.Validate()
		if err != nil {
			t.Errorf("expecting no error, but got %v", err)
		}
	})
}

func TestQuestionnaireValue(t *testing.T) {
	t.Run("valid field value should return marshaled JSON", func(t *testing.T) {
		questionnaire := model.Questionnaire{
			0: {
				CustomName: "Test Group",
				Questions:  []string{"Question 1"},
				Options: []model.AnswerOption{
					{ID: 0, Description: "Option 1", Score: 0},
				},
			},
		}

		fieldValue := questionnaire

		expectedJSON, err := json.Marshal(fieldValue)
		if err != nil {
			t.Fatalf("failed to marshal field value: %v", err)
		}

		value, err := questionnaire.Value(context.Background(), nil, reflect.Value{}, fieldValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(value, expectedJSON) {
			t.Errorf("expected %s, got %s", expectedJSON, value)
		}
	})

	t.Run("invalid field value should return error", func(t *testing.T) {
		questionnaire := model.Questionnaire{}

		invalidFieldValue := make(chan int) // invalid type for JSON marshaling

		_, err := questionnaire.Value(context.Background(), nil, reflect.Value{}, invalidFieldValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})
}

func TestQuestionnaireScan(t *testing.T) {
	t.Run("nil dbValue should return nil error", func(t *testing.T) {
		questionnaire := &model.Questionnaire{}

		err := questionnaire.Scan(context.Background(), nil, reflect.Value{}, nil)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}
	})

	t.Run("valid JSON byte slice should unmarshal correctly", func(t *testing.T) {
		questionnaire := &model.Questionnaire{}
		dbValue := []byte(`{
			"0": {
				"custom_name": "Test Group",
				"questions": ["Question 1"],
				"options": [{"id": 0, "description": "Option 1", "score": 0}]
			}
		}`)

		err := questionnaire.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.Questionnaire{
			0: {
				CustomName: "Test Group",
				Questions:  []string{"Question 1"},
				Options: []model.AnswerOption{
					{ID: 0, Description: "Option 1", Score: 0},
				},
			},
		}

		if !reflect.DeepEqual(*questionnaire, expected) {
			t.Errorf("expected %v, but got %v", expected, *questionnaire)
		}
	})

	t.Run("valid JSON string should unmarshal correctly", func(t *testing.T) {
		questionnaire := &model.Questionnaire{}
		dbValue := `{
			"0": {
				"custom_name": "Test Group",
				"questions": ["Question 1"],
				"options": [{"id": 0, "description": "Option 1", "score": 0}]
			}
		}`

		err := questionnaire.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}

		expected := model.Questionnaire{
			0: {
				CustomName: "Test Group",
				Questions:  []string{"Question 1"},
				Options: []model.AnswerOption{
					{ID: 0, Description: "Option 1", Score: 0},
				},
			},
		}

		if !reflect.DeepEqual(*questionnaire, expected) {
			t.Errorf("expected %v, but got %v", expected, *questionnaire)
		}
	})

	t.Run("invalid type should return error", func(t *testing.T) {
		questionnaire := &model.Questionnaire{}
		dbValue := 12345 // invalid type

		err := questionnaire.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("invalid JSON should return error", func(t *testing.T) {
		questionnaire := &model.Questionnaire{}
		dbValue := []byte(`invalid json`)

		err := questionnaire.Scan(context.Background(), nil, reflect.Value{}, dbValue)
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})
}
func TestIndicationCategoriesValidate(t *testing.T) {
	t.Run("indication categories doesn't cover up to the maxScore", func(t *testing.T) {
		indicationCategories := model.IndicationCategories{
			{MinimumScore: 0, MaximumScore: 5, Name: "Category 1", Detail: "Detail 1"},
		}

		err := indicationCategories.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("value between minScore and maxScore doesn't fit into one of the IndicationCategories", func(t *testing.T) {
		indicationCategories := model.IndicationCategories{
			{MinimumScore: 0, MaximumScore: 5, Name: "Category 1", Detail: "Detail 1"}, // 6 is missing here
			{MinimumScore: 7, MaximumScore: 10, Name: "Category 2", Detail: "Detail 2"},
			{MinimumScore: 11, MaximumScore: 179, Name: "Category 3", Detail: "Detail 3"},
		}

		err := indicationCategories.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("value between minScore and maxScore has more than 1 matching categories", func(t *testing.T) {
		indicationCategories := model.IndicationCategories{
			{MinimumScore: 0, MaximumScore: 5, Name: "Category 1", Detail: "Detail 1"},
			{MinimumScore: 3, MaximumScore: 10, Name: "Category 2", Detail: "Detail 2"},
			{MinimumScore: 11, MaximumScore: 179, Name: "Category 3", Detail: "Detail 3"},
		}

		err := indicationCategories.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("at least one of the supplied indication categories is not valid", func(t *testing.T) {
		indicationCategories := model.IndicationCategories{
			{MinimumScore: 0, MaximumScore: 5, Name: "Category 1", Detail: "Detail 1"},
			{MinimumScore: 6, MaximumScore: 10, Name: "", Detail: "Detail 2"}, // invalid category
			{MinimumScore: 11, MaximumScore: 179, Name: "Category 3", Detail: "Detail 3"},
		}

		err := indicationCategories.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("valid indication categories passed the function", func(t *testing.T) {
		indicationCategories := model.IndicationCategories{
			{MinimumScore: 0, MaximumScore: 5, Name: "Category 1", Detail: "Detail 1"},
			{MinimumScore: 6, MaximumScore: 10, Name: "Category 2", Detail: "Detail 2"},
			{MinimumScore: 11, MaximumScore: 179, Name: "Category 3", Detail: "Detail 3"},
		}

		err := indicationCategories.Validate()
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}
	})
}

func TestIndicationCategoriesValuerAndScanner(t *testing.T) {
	// Define a sample IndicationCategories
	indicationCategories := model.IndicationCategories{
		{
			MinimumScore: 0,
			MaximumScore: 10,
			Name:         "Low",
			Detail:       "Low indication",
		},
		{
			MinimumScore: 11,
			MaximumScore: 20,
			Name:         "Medium",
			Detail:       "Medium indication",
		},
	}

	// Test Value method
	t.Run("Valuer OK", func(t *testing.T) {
		fieldValue := indicationCategories
		valuer := reflect.ValueOf(fieldValue)

		value, err := indicationCategories.Value(context.Background(), nil, valuer, fieldValue)
		if err != nil {
			t.Fatalf("Value() error = %v", err)
		}

		expected, _ := json.Marshal(fieldValue)
		if !reflect.DeepEqual(value, expected) {
			t.Errorf("Value() = %v, want %v", value, expected)
		}
	})

	// Test Scan method
	t.Run("Scan Nil dbValue", func(t *testing.T) {
		var ic model.IndicationCategories

		err := ic.Scan(context.Background(), nil, reflect.ValueOf(ic), nil)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Scan Valid []byte dbValue", func(t *testing.T) {
		var ic model.IndicationCategories

		dbValue := []byte(`[{"minimum_score":0,"maximum_score":10,"name":"Low","detail":"Low indication"},{"minimum_score":11,"maximum_score":20,"name":"Medium","detail":"Medium indication"}]`)

		err := ic.Scan(context.Background(), nil, reflect.ValueOf(ic), dbValue)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}

		if !reflect.DeepEqual(ic, indicationCategories) {
			t.Errorf("Scan() = %v, want %v", ic, indicationCategories)
		}
	})

	t.Run("Scan Valid string dbValue", func(t *testing.T) {
		var ic model.IndicationCategories

		dbValue := `[{"minimum_score":0,"maximum_score":10,"name":"Low","detail":"Low indication"},{"minimum_score":11,"maximum_score":20,"name":"Medium","detail":"Medium indication"}]`

		err := ic.Scan(context.Background(), nil, reflect.ValueOf(ic), dbValue)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}

		if !reflect.DeepEqual(ic, indicationCategories) {
			t.Errorf("Scan() = %v, want %v", ic, indicationCategories)
		}
	})

	t.Run("Scan Invalid dbValue type", func(t *testing.T) {
		var ic model.IndicationCategories

		dbValue := 123

		err := ic.Scan(context.Background(), nil, reflect.ValueOf(ic), dbValue)
		if err == nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Scan Invalid JSON dbValue", func(t *testing.T) {
		var ic model.IndicationCategories

		dbValue := []byte(`invalid json`)

		err := ic.Scan(context.Background(), nil, reflect.ValueOf(ic), dbValue)
		if err == nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, true)
		}
	})
}

func TestIndicationCategoriesGetIndicationCategoryByScore(t *testing.T) {
	ics := model.IndicationCategories{
		{
			MinimumScore: 0,
			MaximumScore: 10,
			Name:         "Low",
			Detail:       "Low indication",
		},
		{
			MinimumScore: 11,
			MaximumScore: 20,
			Name:         "Medium",
			Detail:       "Medium indication",
		},
	}

	t.Run("score is in the range of the first category", func(t *testing.T) {
		ic := ics.GetIndicationCategoryByScore(5)
		if ic.Name != "Low" {
			t.Errorf("expected Low, got %s", ic.Name)
		}
	})

	t.Run("score is in the range of the second category", func(t *testing.T) {
		ic := ics.GetIndicationCategoryByScore(15)
		if ic.Name != "Medium" {
			t.Errorf("expected Medium, got %s", ic.Name)
		}
	})

	t.Run("score is not in the range of any category, must not return nil", func(t *testing.T) {
		ic := ics.GetIndicationCategoryByScore(210 - 0)
		assert.NotNil(t, ic)
	})
}

func TestIndicationCategory(t *testing.T) {
	ic := model.IndicationCategory{
		MinimumScore: 0,
		MaximumScore: 10,
		Name:         "any",
		Detail:       "any",
	}

	t.Run("IsInTheRange? Yes", func(t *testing.T) {
		if !ic.IsInTheRange(5) {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("IsInTheRange? No", func(t *testing.T) {
		if ic.IsInTheRange(11) {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("minimum score is 0", func(t *testing.T) {
		newIc := ic
		newIc.MinimumScore = -1

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("max score atleast is 1", func(t *testing.T) {
		newIc := ic
		newIc.MaximumScore = 0

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("name is required", func(t *testing.T) {
		newIc := ic
		newIc.Name = ""

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("detail is required", func(t *testing.T) {
		newIc := ic
		newIc.Detail = ""

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("minimum score can't be less than predefined template minimum score", func(t *testing.T) {
		newIc := ic
		newIc.MinimumScore = model.DefaultATECTemplate.MinimumPossibleScore - 10

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("maximum score can't be more than predefined template maximum score", func(t *testing.T) {
		newIc := ic
		newIc.MaximumScore = model.DefaultATECTemplate.MaximumPossibleScore + 10

		err := newIc.Validate()
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := ic.Validate()
		if err != nil {
			t.Errorf("expected nil error, but got %v", err)
		}
	})
}

func TestImageResultAttributeKeyValidate(t *testing.T) {
	t.Run("Valid ImageResultAttributeKey", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "Valid Title",
			Total:       "Valid Total",
			Indication:  "Valid Indication",
			ResultID:    "Valid ResultID",
			SubmittedAt: "Valid SubmittedAt",
		}

		if err := iray.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid ImageResultAttributeKey - Missing Title", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "",
			Total:       "Valid Total",
			Indication:  "Valid Indication",
			ResultID:    "Valid ResultID",
			SubmittedAt: "Valid SubmittedAt",
		}

		if err := iray.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid ImageResultAttributeKey - Missing Total", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "Valid Title",
			Total:       "",
			Indication:  "Valid Indication",
			ResultID:    "Valid ResultID",
			SubmittedAt: "Valid SubmittedAt",
		}

		if err := iray.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid ImageResultAttributeKey - Missing Indication", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "Valid Title",
			Total:       "Valid Total",
			Indication:  "",
			ResultID:    "Valid ResultID",
			SubmittedAt: "Valid SubmittedAt",
		}

		if err := iray.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid ImageResultAttributeKey - Missing ResultID", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "Valid Title",
			Total:       "Valid Total",
			Indication:  "Valid Indication",
			ResultID:    "",
			SubmittedAt: "Valid SubmittedAt",
		}

		if err := iray.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid ImageResultAttributeKey - Missing SubmittedAt", func(t *testing.T) {
		iray := model.ImageResultAttributeKey{
			Title:       "Valid Title",
			Total:       "Valid Total",
			Indication:  "Valid Indication",
			ResultID:    "Valid ResultID",
			SubmittedAt: "",
		}

		if err := iray.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})
}

func TestImageResultAttributeKeyValuerAndScanner(t *testing.T) {
	// Define a sample ImageResultAttributeKey
	imageResultAttributeKey := model.ImageResultAttributeKey{
		Title:       "Valid Title",
		Total:       "Valid Total",
		Indication:  "Valid Indication",
		ResultID:    "Valid ResultID",
		SubmittedAt: "Valid SubmittedAt",
	}

	// Test Value method
	t.Run("Valuer OK", func(t *testing.T) {
		fieldValue := imageResultAttributeKey
		valuer := reflect.ValueOf(fieldValue)

		value, err := imageResultAttributeKey.Value(context.Background(), nil, valuer, fieldValue)
		if err != nil {
			t.Fatalf("Value() error = %v", err)
		}

		expected, _ := json.Marshal(fieldValue)
		if !reflect.DeepEqual(value, expected) {
			t.Errorf("Value() = %v, want %v", value, expected)
		}
	})

	// Test Scan method
	t.Run("Scan Nil dbValue", func(t *testing.T) {
		var iray model.ImageResultAttributeKey

		err := iray.Scan(context.Background(), nil, reflect.ValueOf(iray), nil)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}

		if iray != (model.ImageResultAttributeKey{}) {
			t.Error("scanning nil for ImageResultAttributeKey should yield empty value")
		}
	})

	t.Run("Scan Valid []byte dbValue", func(t *testing.T) {
		var iray model.ImageResultAttributeKey

		dbValue := []byte(`{"title":"Valid Title","total":"Valid Total","indication":"Valid Indication","result_id":"Valid ResultID","submitted_at":"Valid SubmittedAt"}`)

		err := iray.Scan(context.Background(), nil, reflect.ValueOf(iray), dbValue)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}

		if !reflect.DeepEqual(iray, imageResultAttributeKey) {
			t.Errorf("Scan() = %v, want %v", iray, imageResultAttributeKey)
		}
	})

	t.Run("Scan Valid string dbValue", func(t *testing.T) {
		var iray model.ImageResultAttributeKey

		dbValue := `{"title":"Valid Title","total":"Valid Total","indication":"Valid Indication","result_id":"Valid ResultID","submitted_at":"Valid SubmittedAt"}`

		err := iray.Scan(context.Background(), nil, reflect.ValueOf(iray), dbValue)
		if err != nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, false)
		}

		if !reflect.DeepEqual(iray, imageResultAttributeKey) {
			t.Errorf("Scan() = %v, want %v", iray, imageResultAttributeKey)
		}
	})

	t.Run("Scan Invalid dbValue type", func(t *testing.T) {
		var iray model.ImageResultAttributeKey

		dbValue := 123

		err := iray.Scan(context.Background(), nil, reflect.ValueOf(iray), dbValue)
		if err == nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Scan Invalid JSON dbValue", func(t *testing.T) {
		var iray model.ImageResultAttributeKey

		dbValue := []byte(`invalid json`)

		err := iray.Scan(context.Background(), nil, reflect.ValueOf(iray), dbValue)
		if err == nil {
			t.Errorf("Scan() error = %v, wantErr %v", err, true)
		}
	})
}
