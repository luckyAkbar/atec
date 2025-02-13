package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// @Summary		Create new ATEC questionaire package
// @Description	Create new ATEC questionaire package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			create_package_input	body		CreatePackageInput									true	"ATEC questionnarie package details"
// @Success		200						{object}	StandardSuccessResponse{data=CreatePackageOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/atec/packages [post]
//
// @Example
//
//	{
//	  "package_name": "string",
//	  "image_result_attribute_key": {
//	    "indication": "Indikasi",
//	    "result_id": "ID Hasil",
//	    "submitted_at": "Dikerjakan Pada",
//	    "title": "Skor ATEC",
//	    "total": "Total Skor"
//	  },
//	  "indication_categories": [
//	    {
//	      "detail": "Menunjukkan kemungkinan bahwa anak memiliki pola perilaku dan keterampuilan komunikasi yang agak normal dan memiliki peluang tinggi untuk menjalani kehidupan normal dan independen yang menunjukkan gejala ASD minimal.",
//	      "maximum_score": 30,
//	      "minimum_score": 0,
//	      "name": "gejala ringan"
//	    },
//	    {
//	      "detail": "Menunjukkan kemungkinan bahwa anak kemungkinan besar akan dapat menjalani kehidupan semi independen tanpa perlu ditempatkan di fasilitas perawatan formal",
//	      "maximum_score": 50,
//	      "minimum_score": 31,
//	      "name": "gejala sedang"
//	    },
//	    {
//	      "detail": "Menunjukkan kemungkinan bahwa anak jatuh ke persentil ke-90 (sangat autis). Anak mungkin akan membutuhkan perawatan berkelanjutan (mungkin di sebuah institusi), dan mungkin tidak dapat mencapai tingkat kebebasan apapun dari orang lain",
//	      "maximum_score": 179,
//	      "minimum_score": 51,
//	      "name": "gejala berat"
//	    }
//	  ],
//	  "questionnaire": {
//	    "0": {
//	      "custom_name": "Kemampuan Bicara/Berbahasa",
//	      "options": [
//	        {
//	          "description": "Tidak Benar",
//	          "id": 2,
//	          "score": 2
//	        },
//	        {
//	          "description": "Agak Benar",
//	          "id": 1,
//	          "score": 1
//	        },
//	        {
//	          "description": "Sangat Benar",
//	          "id": 0,
//	          "score": 0
//	        }
//	      ],
//	      "questions": [
//	        "Mengetahui namanya sendiri",
//	        "Berespon pada \"Tidak\" atau \"Stop\"",
//	        "Dapat mengikuti perintah",
//	        "Dapat menggunakan 1 kata (Tidak!, Makan, Air, dll)",
//	        "Dapat menggunakan 2 kata sekaligus bersamaan (Tidak mau!, Pergi pulang, dll)",
//	        "Dapat menggunakan 3 kata sekaligus bersamaan (Mau minum susu, dll)",
//	        "Mengetahui 10 kata atau lebih",
//	        "Dapat membuat kalimat yang berisi 4 kata atau lebih",
//	        "Mampu menjelaskan apa yang dia inginkan",
//	        "Mampu menanyakan pertanyaan yang bermakna",
//	        "Isi pembicaraan cenderung relevan/bermakna",
//	        "Sering menggunakan kalimat-kalimat yang berurutan",
//	        "Bisa mengikuti pembicaraan dengan cukup baik",
//	        "Memiliki kemampuan bicara/berbahasa yang sesuai dengan seusianya"
//	      ]
//	    },
//	    "1": {
//	      "custom_name": "Kemampuan Bersosialisasi",
//	      "options": [
//	        {
//	          "description": "Tidak Cocok",
//	          "id": 0,
//	          "score": 0
//	        },
//	        {
//	          "description": "Agak Cocok",
//	          "id": 1,
//	          "score": 1
//	        },
//	        {
//	          "description": "Sangat Cocok",
//	          "id": 2,
//	          "score": 2
//	        }
//	      ],
//	      "questions": [
//	        "Terlihat seperti berada dalam ‘tempurung’ – Anda tidak bisa menjangkau dia",
//	        "Mengabaikan orang lain",
//	        "Ketika dipanggil, hanya sedikit atau malah tidak memperhatikan",
//	        "Tidak kooperatif dan menolak",
//	        "Tidak ada kontak mata",
//	        "Lebih suka menyendiri",
//	        "Tidak menunjukkan rasa kasih sayang",
//	        "Tidak mampu menyapa orang tua",
//	        "Menghindari kontak dengan orang lain",
//	        "Tidak mampu menirukan orang lain",
//	        "Tidak suka dipegang atau dipeluk",
//	        "Tidak mau berbagi atau menunjukkan",
//	        "Tidak bisa melambaikan tangan \"Da..Dahh\"",
//	        "Sering tidak setuju / menolak (not compliant)",
//	        "Tantrum, marah-marah",
//	        "Tidak mempunyai teman",
//	        "Jarang tersenyum",
//	        "Tidak peka terhadap perasaan orang lain",
//	        "Acuh tak acuh ketika disukai orang lain",
//	        "Acuh tak acuh ketika ditinggal pergi oleh orang tuanya"
//	      ]
//	    },
//	    "2": {
//	      "custom_name": "Kesadaran Sensori/Kognitif",
//	      "options": [
//	        {
//	          "description": "Sangat Cocok",
//	          "id": 0,
//	          "score": 0
//	        },
//	        {
//	          "description": "Agak Cocok",
//	          "id": 1,
//	          "score": 1
//	        },
//	        {
//	          "description": "Tidak Cocok",
//	          "id": 2,
//	          "score": 2
//	        }
//	      ],
//	      "questions": [
//	        "Merespon saat dipanggil namanya",
//	        "Merespon saat dipuji",
//	        "Melihat pada orang dan binatang",
//	        "Melihat pada gambar (dan TV)",
//	        "Menggambar, mewarnai dan melakukan kesenian",
//	        "Bermain dengan mainannya secara sesuai",
//	        "Menggunakan ekspresi wajah yang sesuai",
//	        "Memahami cerita yang ditayangkan di TV",
//	        "Memahami penjelasan",
//	        "Sadar akan lingkungannya",
//	        "Sadar akan bahaya",
//	        "Mampu berimajinasi",
//	        "Memulai aktivitas",
//	        "Mampu berpakaian sendiri",
//	        "Memiliki rasa penasaran dan ketertarikan",
//	        "Suka tantangan, senang mengeksplorasi",
//	        "Tampak selaras, tidak tampak ‘kosong’",
//	        "Mampu mengikuti pandangan ke arah semua orang memandang."
//	      ]
//	    },
//	    "3": {
//	      "custom_name": "Kesehatan Umum, Fisik dan Perilaku",
//	      "options": [
//	        {
//	          "description": "Sangat Bermasalah",
//	          "id": 3,
//	          "score": 3
//	        },
//	        {
//	          "description": "Cukup Bermasalah",
//	          "id": 2,
//	          "score": 2
//	        },
//	        {
//	          "description": "Sedikit Bermasalah",
//	          "id": 1,
//	          "score": 1
//	        },
//	        {
//	          "description": "Tidak bermasalah",
//	          "id": 0,
//	          "score": 0
//	        }
//	      ],
//	      "questions": [
//	        "Mengompol saat tidur",
//	        "Mengompol di celana/popok",
//	        "Buang air besar di celana/popok",
//	        "Diare",
//	        "Konstipasi / Sembelit",
//	        "Gangguan Tidur",
//	        "Makan terlalu banyak / terlalu sedikit",
//	        "Pilihan makanan yang diinginkan sangat terbatas (extremely limited diet, picky eater)",
//	        "Hiperaktif",
//	        "Letargi, lemah, lesu",
//	        "Memukul atau melukai diri sendiri",
//	        "Memukul atau melukai orang lain",
//	        "Destruktif",
//	        "Sensitif terhadap suara",
//	        "Cemas / penuh ketakutan",
//	        "Tidak senang/ mudah rewel/ menangis",
//	        "Kejang",
//	        "Bicara secara obsesif",
//	        "Kaku terhadap rutinitas",
//	        "Berteriak / menjerit-jerit",
//	        "Menuntut hal atau cara yang sama berulang-ulang",
//	        "Sering gelisah / agitasi",
//	        "Tidak peka terhadap nyeri",
//	        "Terfokus atau sulit dialihkan dari objek atau topik tertentu",
//	        "Gerakan repetitive (stimming, menggoyang-goyangkan bagian badan)"
//	      ]
//	    }
//	  }
//	}
//
//nolint:lll
func (s *service) HandleCreatePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &CreatePackageInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.packageUsecase.Create(c.Request().Context(), usecase.CreatePackageInput{
			PackageName:             input.PackageName,
			Questionnaire:           input.Quesionnaire,
			IndicationCategories:    input.IndicationCategories,
			ImageResultAttributeKey: input.ImageResultAttributeKey,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: CreatePackageOutput{
				ID: output.ID,
			},
		})
	}
}

// @Summary		Update existing ATEC questionnarie package
// @Description	Update existing ATEC questionnarie package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			package_id				path		string												true	"package id (UUID v4)"
// @Param			update_package_input	body		UpdatePackageInput									true	"ATEC questionnarie package details"
// @Success		200						{object}	StandardSuccessResponse{data=UpdatePackageOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/atec/packages/{package_id} [put]
func (s *service) HandleUpdatePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &UpdatePackageInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.packageUsecase.Update(c.Request().Context(), usecase.UpdatePackageInput{
			PackageID:     input.PackageID,
			PackageName:   input.PackageName,
			Questionnaire: input.Quesionnaire,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: UpdatePackageOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Update package activation status
// @Description	Update existing ATEC questionnarie package activation status
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization				header		string													true	"JWT Token"
//
// @Param			package_id					path		string													true	"package ID to be activated/deactivated (UUID v4)"
// @Param			activation_package_input	body		ActivationPackageInput									true	"activation status"
// @Success		200							{object}	StandardSuccessResponse{data=ActivationPackageOutput}	"Successful response"
// @Failure		400							{object}	StandardErrorResponse									"Bad request"
// @Failure		500							{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/atec/packages/{package_id} [patch]
func (s *service) HandleActivationPackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &ActivationPackageInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.packageUsecase.ChangeActiveStatus(c.Request().Context(), usecase.ChangeActiveStatusInput{
			PackageID:    input.PackageID,
			ActiveStatus: input.Status,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: ActivationPackageOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Delete ATEC questionnaire package
// @Description	Delete ATEC questionnaire package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization	header	string	true	"JWT Token"
//
// @Param			package_id		path	string	true	"package ID to be deleted (UUID v4)"
// @Success		200				"No Content"
// @Failure		400				{object}	StandardErrorResponse	"Bad request"
// @Failure		500				{object}	StandardErrorResponse	"Internal Error"
// @Router			/v1/atec/packages/{package_id} [delete]
func (s *service) HandleDeletePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &DeletePackageInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		err := s.packageUsecase.Delete(c.Request().Context(), input.PackageID)
		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// @Summary		Get all active packages
// @Description	Get all active packages
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=[]SearchActivePackageOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse										"Bad request"
// @Failure		500	{object}	StandardErrorResponse										"Internal Error"
// @Router			/v1/atec/packages/active [get]
func (s *service) HandleSearchActivePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		packages, err := s.packageUsecase.FindActiveQuestionnaires(c.Request().Context())
		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		output := []SearchActivePackageOutput{}
		for _, pack := range packages {
			output = append(output, SearchActivePackageOutput{
				ID:                   pack.ID,
				Questionnaire:        pack.Questionnaire,
				IndicationCategories: pack.IndicationCategories,
				Name:                 pack.Name,
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       output,
		})
	}
}
