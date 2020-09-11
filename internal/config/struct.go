package config

import "adeia/internal/util/constants"

// envOverrides holds all environment value keys for overriding the config.
var envOverrides = map[string]string{
	// server overrides
	"server.jwt_secret": constants.EnvServerJWTSecretKey,

	// mailer overrides
	"mailer.username": constants.EnvMailerUsernameKey,
	"mailer.password": constants.EnvMailerPasswordKey,

	// database overrides
	"database.dbname":   constants.EnvDBNameKey,
	"database.user":     constants.EnvDBUserKey,
	"database.password": constants.EnvDBPasswordKey,
	"database.host":     constants.EnvDBHostKey,
	"database.port":     constants.EnvDBPortKey,

	// cache overrides
	"cache.host": constants.EnvCacheHostKey,
	"cache.port": constants.EnvCachePortKey,
}

type Config struct {
	CacheConfig  `mapstructure:"cache"`
	DBConfig     `mapstructure:"database"`
	LoggerConfig `mapstructure:"logger"`
	MailerConfig `mapstructure:"mailer"`
	ServerConfig `mapstructure:"server"`
}

type CacheConfig struct {
	Network  string `mapstructure:"network"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	ConnSize int    `mapstructure:"connsize"`
}

type DBConfig struct {
	Driver      string `mapstructure:"driver"`
	DBName      string `mapstructure:"dbname"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	SSLMode     string `mapstructure:"sslmode"`
	SSLCert     string `mapstructure:"sslcert,omitempty"`
	SSLKey      string `mapstructure:"sslkey,omitempty"`
	SSLRootCert string `mapstructure:"sslrootcert,omitempty"`
}

type LoggerConfig struct {
	Level string   `mapstructure:"level"`
	Paths []string `mapstructure:"paths"`
}

type MailerConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
}

type ServerConfig struct {
	Host            string `mapstructure:"host,omitempty"`
	Port            int    `mapstructure:"port"`
	RateLimitRate   int    `mapstructure:"ratelimit_rate"`
	RateLimitWindow int    `mapstructure:"ratelimit_window"`
	JWTSecret       string `mapstructure:"jwt_secret"`
}
