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

	rootCMD.AddCommand(serverCMD)
}

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
		init_admin_account(userRepo, sharedCryptor)
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

// init_admin_account is a function to initialize admin account to be used in the application.
// This works by first checking whether an account with Admin level are already created or not.
// If found nothing, will initialize one using predetermined email and password read from environment values.
// if the envionment values are not set, will fail.
func init_admin_account(userRepo *repository.UserRepository, sc *common.SharedCryptor) {
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
