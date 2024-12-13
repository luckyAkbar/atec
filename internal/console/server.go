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
	rootCMD.AddCommand(serverCMD)
}

func serverFn(_ *cobra.Command, _ []string) {
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
	childUsecase := usecase.NewChildUsecase(childRepo)
	questionnaireUsecase := usecase.NewQuestionnaireUsecase(packageRepo, childRepo, resultRepo, font)

	httpServer := echo.New()

	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

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
