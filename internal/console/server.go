package console

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/db"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sweet-go/stdlib/encryption"
	"github.com/sweet-go/stdlib/helper"
	"golang.org/x/crypto/bcrypt"
)

var serverCMD = &cobra.Command{
	Use:  "server",
	Long: "run the API server",
	Run:  serverFn,
}

//nolint:gochecknoinits
func init() {
	serverCMD.Flags().Bool("init-admin-account", false, "if set, will check whether an admin account already exists or not, if not will create one")
	serverCMD.Flags().Bool("init-package", false, "if set, will check whether a package already exists or not, if not will create one")

	rootCMD.AddCommand(serverCMD)
}

//nolint:funlen
func serverFn(cmd *cobra.Command, _ []string) {
	key, err := encryption.ReadKeyFromFile(config.PrivateKeyFilePath())
	if err != nil {
		panic(err)
	}

	fontBytes, err := os.ReadFile("./assets/font.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	sharedCryptor := common.NewSharedCryptor(&common.CreateCryptorOpts{
		HashCost:      bcrypt.DefaultCost,
		EncryptionKey: key.Bytes,
		IV:            config.IVKey(),
		BlockSize:     common.DefaultBlockSize,
	})

	brevoClient := common.NewBrevoClient(config.SendinblueAPIKey())

	mailer := common.NewMailer(config.SendInBlueSender().Name, config.SendInBlueSender().Email, brevoClient)

	db.InitializePostgresConn()

	userRepo := repository.NewUserRepository(db.PostgresDB)
	packageRepo := repository.NewPackageRepo(db.PostgresDB)
	childRepo := repository.NewChildRepository(db.PostgresDB)
	resultRepo := repository.NewResultRepository(db.PostgresDB)

	authUsecase := usecase.NewAuthUsecase(sharedCryptor, userRepo, mailer)
	packageUsecase := usecase.NewPackageUsecase(packageRepo)
	childUsecase := usecase.NewChildUsecase(childRepo, resultRepo)
	questionnaireUsecase := usecase.NewQuestionnaireUsecase(packageRepo, childRepo, resultRepo, font)

	initAdmin, err := cmd.Flags().GetBool("init-admin-account")
	if err != nil {
		panic(err)
	}

	if initAdmin {
		initAdminAccount(userRepo, sharedCryptor)
	}

	initPackageFlag, err := cmd.Flags().GetBool("init-package")
	if err != nil {
		panic(err)
	}

	if initPackageFlag {
		initPackage(packageRepo, userRepo)
	}

	httpServer := echo.New()

	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

	rootGroup := httpServer.Group("")
	rootGroup.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{
			Message: "pong",
		})
	})

	v1Group := httpServer.Group("v1")

	rest.NewService(v1Group, authUsecase, packageUsecase, childUsecase, questionnaireUsecase)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)

	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for {
			select {
			case <-sigCh:
				logrus.Info("shutting down the server")
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				logrus.Error(e)
				gracefulShutdown(httpServer)
				quitCh <- true
			}
		}
	}()

	// Start server
	go func() {
		if err := httpServer.Start(config.ServerPort()); err != nil && err != http.ErrServerClosed {
			httpServer.Logger.Fatal("shutting down the server: ", err.Error())
		}
	}()

	<-quitCh
	logrus.Info("exiting")
}

func gracefulShutdown(srv *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:mnd
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal(err)
	}
}

// initAdminAccount is a function to initialize admin account to be used in the application.
// This works by first checking whether an account with Admin level are already created or not.
// If found nothing, will initialize one using predetermined email and password read from environment values.
// if the envionment values are not set, will fail.
func initAdminAccount(userRepo *repository.UserRepository, sc *common.SharedCryptor) {
	adminEmail := os.Getenv("INIT_ADMIN_EMAIL")
	adminPassword := os.Getenv("INIT_ADMIN_PASSWORD")
	adminUsername := os.Getenv("INIT_ADMIN_USERNAME")

	found, err := userRepo.IsAdminAccountExists(context.Background())
	if err != nil {
		logrus.WithError(err).Error("failed to check whether admin account exists on database or not")
		os.Exit(1)
	}

	if found {
		logrus.Info("admin account already exists, skipping initialization")

		return
	}

	if adminEmail == "" || adminPassword == "" {
		logrus.Error("admin email or password is not set when running initialization of admin account")
		os.Exit(1)
	}

	if adminUsername == "" {
		adminUsername = "admin"
	}

	emailEncrypted, err := sc.Encrypt(adminEmail)
	if err != nil {
		logrus.WithError(err).Error("failed to encrypt admin email when initializing the first admin account")
		os.Exit(1)
	}

	hashedPassword, err := sc.Hash([]byte(adminPassword))
	if err != nil {
		logrus.WithError(err).Error("failed to hash admin password when initializing the first admin account")
		os.Exit(1)
	}

	adminUser, err := userRepo.Create(context.Background(), repository.CreateUserInput{
		Email:    emailEncrypted,
		Password: hashedPassword,
		Username: adminUsername,
		IsActive: true,
		Roles:    model.RolesAdmin,
	})

	if err != nil {
		logrus.WithError(err).Error("failed to create admin account when initializing the first admin account")
		os.Exit(1)
	}

	logrus.Info("admin account created: ", adminUser.ID)
}

// initPackage is a function to initialize package to be used in the application. Only works if no active package found.
// Currently, this function will only create a package in Bahasa Indonesia, taken from ATEC Jatmika's website.
// If however, more languages option want to be implemented here, it should update this function to accept more parameters.
//
//nolint:funlen,lll,mnd
func initPackage(packageRepo repository.PackageRepoIface, userRepo repository.UserRepositoryIface) {
	ctx := context.Background()
	truth := true
	logger := logrus.WithField("function", "initPackage")

	activePackages, err := packageRepo.Search(ctx, repository.SearchPackageInput{
		IsActive: &truth,
	})

	switch err {
	default:
		logger.WithError(err).Error("failed to perform package search, unable to initialize package")
		os.Exit(1)
	case nil:
		logger.WithFields(logrus.Fields{
			"activePackageCount": len(activePackages),
			"activePackages":     helper.Dump(activePackages),
		}).Info("already found atleast one active package, skipping initialization")

		return
	case repository.ErrNotFound:
		break
	}

	logger.Info("no active package found, initializing package")
	logger.Info("start to find admin account to be used as package creator")

	admin, err := userRepo.Search(ctx, repository.SearchUserInput{
		Role:   model.RolesAdmin,
		Limit:  1,
		Offset: 0,
	})

	switch err {
	default:
		logger.WithError(err).Error("failed to find admin account, unable to initialize package")
		os.Exit(1)
	case repository.ErrNotFound:
		logger.Error("no admin account found, unable to initialize package. you might need to create it first")
		os.Exit(1)
	case nil:
		break
	}

	adminAccount := admin[0]

	logger.WithField("adminAccountID", adminAccount.ID).Info("admin account found, creating package")

	input := repository.CreatePackageInput{
		UserID:      adminAccount.ID,
		PackageName: "Kuesioner ATEC Bahasa Indonesia dari ATEC Jatmika",
		Questionnaire: model.Questionnaire{
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
		},
		IndicationCategories: model.IndicationCategories{
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
		},
		ImageResultAttributeKey: model.ImageResultAttributeKey{
			Title:       "Skor ATEC",
			Total:       "Total Skor",
			Indication:  "Indikasi",
			ResultID:    "ID Hasil",
			SubmittedAt: "Dikerjakan Pada",
		},
	}

	tx := db.PostgresDB.Begin()

	createdPackage, err := packageRepo.Create(ctx, input, tx)
	if err != nil {
		logger.WithError(err).Error("failed to create package")
		tx.Rollback()
		os.Exit(1)
	}

	logger.WithField("packageID", createdPackage.ID).Info("package created successfully. will try to activate it")

	_, err = packageRepo.Update(ctx, createdPackage.ID, repository.UpdatePackageInput{
		ActiveStatus: &truth,
	}, tx)

	if err != nil {
		logger.WithError(err).Error("failed to activate package")
		tx.Rollback()
		os.Exit(1)
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(err).Error("failed to commit transaction")
		os.Exit(1)
	}

	logger.Info("package activated successfully")
}
