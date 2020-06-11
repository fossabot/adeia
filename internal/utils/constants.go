package utils

const (
	// EnvPrefix is used as the prefix for all env variables related to adeia.
	EnvPrefix = "ADEIA"

	// ==========
	// Keys of env variables to override config from config.yaml
	// ==========
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

	// EnvCacheNetworkKey is the env key for redis network.
	EnvCacheNetworkKey = EnvPrefix + "_CACHE_NETWORK"

	// EnvCacheAddrKey is the env key for redis address.
	EnvCacheAddrKey = EnvPrefix + "_CACHE_ADDR"

	// EnvCacheConnsizeKey is the env key for redis connection pool size.
	EnvCacheConnsizeKey = EnvPrefix + "_CACHE_CONNSIZE"
)
