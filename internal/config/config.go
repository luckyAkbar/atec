// Package config contains all function to access system's configuration
package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/sendinblue/APIv3-go-library/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// init this function is used to read the config file, preferably in yaml format (config.yaml)
// this was made to prevent failure to read config file when running inside testing environment where the config file is not available and not needed
//
//nolint:gochecknoinits
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Error("failed to read config file")

		if !testing.Testing() {
			panic("missing config file while run outside testing environment")
		}
	}
}

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
	Package                time.Duration
} {
	return struct {
		AllActivePackage       time.Duration
		DefaultNilCacheTimeout time.Duration
		Package                time.Duration
	}{
		AllActivePackage:       viper.GetDuration("caching.expiration.all_active_packages"),
		DefaultNilCacheTimeout: viper.GetDuration("caching.expiration.default_nil_value"),
		Package:                viper.GetDuration("caching.expiration.package"),
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

// ResendSignupVerificationLimiterDuration resend signup verification limiter duration.
// If left unset or below 1 minutes, will return the default duration of 5 minutes.
func ResendSignupVerificationLimiterDuration() time.Duration {
	const defaultDurationMinutes = 5

	const minimumDurationMinutes = 1

	defaultDuration := defaultDurationMinutes * time.Minute
	minimumDuration := minimumDurationMinutes * time.Minute

	cfg := viper.GetDuration("resend_signup_verification_limiter_duration")
	if cfg == 0 || cfg < minimumDuration {
		return defaultDuration
	}

	return cfg
}

// DBMaxIdleConn max idle conn
func DBMaxIdleConn() int {
	const defaultMaxIdleConn = 30

	cfg := viper.GetInt("postgres.pool.max_idle_conn")
	if cfg == 0 {
		return defaultMaxIdleConn
	}

	return cfg
}

// DBMaxOpenConn max open conn
func DBMaxOpenConn() int {
	const defaultMaxOpenConn = 50

	cfg := viper.GetInt("postgres.pool.max_open_conn")
	if cfg == 0 {
		return defaultMaxOpenConn
	}

	return cfg
}

// DBConnMaxLifetime max lifetime
func DBConnMaxLifetime() time.Duration {
	const defaultMaxLifetime = 30 * time.Minute

	cfg := viper.GetInt("postgres.pool.conn_max_lifetime_sec")
	if cfg == 0 {
		return defaultMaxLifetime
	}

	return time.Duration(cfg) * time.Second
}
