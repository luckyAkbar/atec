package config

import (
	"fmt"
	"time"

	"github.com/sendinblue/APIv3-go-library/lib"
	"github.com/spf13/viper"

	// used to indirectly call the init function of this stdlib/config package
	_ "github.com/sweet-go/stdlib/config"
)

func LogLevel() string {
	return viper.GetString("log_level")
}

func PostgresDSN() string {
	host := viper.GetString("postgres.host")
	db := viper.GetString("postgres.db")
	user := viper.GetString("postgres.user")
	pw := viper.GetString("postgres.pw")
	port := viper.GetString("postgres.port")
	sslMode := viper.GetString("postgres.ssl_mode")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, pw, db, port, sslMode)
}

// ServerPort returns application server port
func ServerPort() string {
	return fmt.Sprintf(":%s", viper.GetString("server.port"))
}

func PrivateKeyFilePath() string {
	return viper.GetString("private_key_path")
}

func IVKey() string {
	return viper.GetString("iv_key")
}

func SignupTokenExpiry() time.Duration {
	return viper.GetDuration("signup_token_expiry")
}

// SendinblueAPIKey get API key for send in blue
func SendinblueAPIKey() string {
	return viper.GetString("sendinblue.api_key")
}

// SendInBlueSender generate sendinblue sender using configured sender name and sender email
func SendInBlueSender() *lib.SendSmtpEmailSender {
	return &lib.SendSmtpEmailSender{
		Name:  viper.GetString("sendinblue.sender_name"),
		Email: viper.GetString("sendinblue.sender_email"),
	}
}

// SendInBlueIsActivated is activated sendinblue
func SendInBlueIsActivated() bool {
	return viper.GetBool("sendinblue.is_activated")
}
