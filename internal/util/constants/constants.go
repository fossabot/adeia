package constants

// TimeUnit is used to distinguish between different time formats in requests.
type TimeUnit string

const (
	// APIVersion represents the current major version of the API. It is used as URL prefix.
	APIVersion = "v1"

	// EmployeeIDLength represents the length of the generated employee IDs.
	EmployeeIDLength = 6

	// ==========
	// Session-related constants
	// ==========

	// ContextUserKey is the key used to store the authenticated user in the context.
	ContextUserKey = "user"
	// RefreshTokenCookieName is the name of the cookie that stores the refresh token.
	RefreshTokenCookieName = "token"
	// AccessTokenExpiry is the max-age of the access token (jwt) (in seconds).
	AccessTokenExpiry = 30 * 60 // 30 minutes
	// RefreshTokenExpiry is the max-age of the refresh token in the cookie (in seconds).
	RefreshTokenExpiry = 7 * 24 * 60 * 60 // 7 days

	// ==========
	// Keys of env variables to override config from config.yaml
	// ==========

	// EnvPrefix is used as the prefix for all env variables related to adeia.
	EnvPrefix = "ADEIA"

	// Server keys

	// EnvServerJWTSecretKey is the env key for server's jwt secret.
	EnvServerJWTSecretKey = EnvPrefix + "_SERVER_JWT_SECRET"

	// Mailer keys

	// EnvMailerUsernameKey is the env key for mailer username.
	EnvMailerUsernameKey = EnvPrefix + "_MAILER_USERNAME"
	// EnvMailerPasswordKey is the env key for mailer password.
	EnvMailerPasswordKey = EnvPrefix + "_MAILER_PASSWORD"

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

	// ==========
	// General constants
	// ==========

	// Epoch represents the TimeUnit epoch.
	Epoch TimeUnit = "Epoch"
	// Month represents the TimeUnit Month.
	Month TimeUnit = "Month"
	// DayOfMonth represents the TimeUnit DayOfMonth.
	DayOfMonth TimeUnit = "DayOfMonth"
	// Year represents the TimeUnit Year.
	Year TimeUnit = "Year"
)
