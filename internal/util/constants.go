package util

const (
	// ==========
	// Keys of env variables to override config from config.yaml
	// ==========

	// EnvPrefix is used as the prefix for all env variables related to adeia.
	EnvPrefix = "ADEIA"

	// Mailer keys

	// EnvMailerUsername is the env key for mailer username.
	EnvMailerUsername = EnvPrefix + "_MAILER_USERNAME"
	// EnvMailerPassword is the env key for mailer password.
	EnvMailerPassword = EnvPrefix + "_MAILER_PASSWORD"

	// Database (Postgres) keys

	// EnvConfPathKey is the env key for confPath.
	EnvConfPathKey = EnvPrefix + "_CONF_PATH"
	// EnvDBNameKey is the env key for database name.
	EnvDBNameKey = EnvPrefix + "_DB_NAME"
	// EnvDBUserKey is the env key for database user.
	EnvDBUserKey = EnvPrefix + "_DB_USER"
	// EnvDBPasswordKey is the env key for database password.
	EnvDBPasswordKey = EnvPrefix + "_DB_PASSWORD"
	// EnvDBHostKey is the env key for database host.
	EnvDBHostKey = EnvPrefix + "_DB_HOST"
	// EnvDBPortKey is the env key for database port.
	EnvDBPortKey = EnvPrefix + "_DB_PORT"

	// Cache keys

	// EnvCacheHostKey is the env key for redis host.
	EnvCacheHostKey = EnvPrefix + "_CACHE_HOST"
	// EnvCachePortKey is the env key for redis port.
	EnvCachePortKey = EnvPrefix + "_CACHE_PORT"
)
