// Package config contains all function to access system's configuration
package config

import (
	"fmt"
	"time"

	"github.com/sendinblue/APIv3-go-library/lib"
	"github.com/spf13/viper"

	// used to indirectly call the init function of this stdlib/config package
	_ "github.com/sweet-go/stdlib/config"
)

// LogLevel log level
func LogLevel() string {
	return viper.GetString("log_level")
}

// PostgresDSN postgres dsn
func PostgresDSN() string {
	host := viper.GetString("postgres.host")
	db := viper.GetString("postgres.db")
	user := viper.GetString("postgres.user")
	pw := viper.GetString("postgres.pw")
	port := viper.GetString("postgres.port")
	sslMode := viper.GetString("postgres.ssl_mode")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, pw, db, port, sslMode)
}

// ServerAccountVerificationBaseURL contains the url for user when clicking the
// validate account button on account verification email after signing up. Could
// be used to point to the front end page along with the required validation token
func ServerAccountVerificationBaseURL() string {
	return viper.GetString("server.account_verification_base_url")
}

// ServerPort returns application server port
func ServerPort() string {
	return ":" + viper.GetString("server.port")
}

// PrivateKeyFilePath private key path
func PrivateKeyFilePath() string {
	return viper.GetString("private_key_path")
}

// IVKey iv key
func IVKey() string {
	return viper.GetString("iv_key")
}

// SignupTokenExpiry expiry time in time.Duration
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

// LoginTokenExpiry expiry time for login token in time.Duration
func LoginTokenExpiry() time.Duration {
	return viper.GetDuration("login_token_expiry")
}

// ChangePasswordTokenExpiry change password token expiry in time.Duration
func ChangePasswordTokenExpiry() time.Duration {
	return viper.GetDuration("change_password_token_expiry")
}

// ServerResetPasswordBaseURL reset password base URL
func ServerResetPasswordBaseURL() string {
	return viper.GetString("server.reset_password_base_url")
}

// RedisAddr get redis address
func RedisAddr() string {
	return viper.GetString("caching.redis.host")
}

// RedisPassword get redis password
func RedisPassword() string {
	return viper.GetString("caching.redis.password")
}

// RedisDB get redis db
func RedisDB() int {
	return viper.GetInt("caching.redis.db")
}

// RedisMinIdleConns get redis min idle connections
func RedisMinIdleConns() int {
	return viper.GetInt("caching.redis.min_idle_conns")
}

// RedisConnMaxLifetimeSec get redis connection max lifetime in seconds
func RedisConnMaxLifetimeSec() int {
	return viper.GetInt("caching.redis.conn_max_lifetime_sec")
}

// CacheExpiryDuration get cache expiry duration
func CacheExpiryDuration() struct {
	AllActivePackage       time.Duration
	DefaultNilCacheTimeout time.Duration
} {
	return struct {
		AllActivePackage       time.Duration
		DefaultNilCacheTimeout time.Duration
	}{
		AllActivePackage:       viper.GetDuration("caching.expiration.all_active_packages"),
		DefaultNilCacheTimeout: viper.GetDuration("caching.expiration.default_nil_value"),
	}
}

// RedisLockAddr redis lock address
func RedisLockAddr() string {
	return viper.GetString("caching.redis_lock.host")
}

// RedisLockPassword redis lock password
func RedisLockPassword() string {
	return viper.GetString("caching.redis_lock.password")
}

// RedisLockDB redis lock db
func RedisLockDB() int {
	return viper.GetInt("caching.redis_lock.db")
}

// RedisLockMinIdleConns redis lock min idle connections
func RedisLockMinIdleConns() int {
	return viper.GetInt("caching.redis_lock.min_idle_conns")
}

// RedisLockConnMaxLifetimeSec redis lock connection max lifetime in seconds
func RedisLockConnMaxLifetimeSec() int {
	return viper.GetInt("caching.redis_lock.conn_max_lifetime_sec")
}
