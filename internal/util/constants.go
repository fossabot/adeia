package util

const (
	// EmployeeIDLength represents the length of the generated employee IDs.
	EmployeeIDLength = 6

	// ContextUserKey is the key used to store authenticated user in the context.
	ContextUserKey = "user"

	// SessionCookieKey is the name of the session cookie.
	SessionCookieKey = "id"

	// SessionExpiry is the max-age of the session cookie (in seconds).
	SessionExpiry = 24 * 60 * 60

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
