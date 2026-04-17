package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Mode         string        `mapstructure:"mode"` // debug, release, test
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// Trust configuration
	TrustedProxies []string `mapstructure:"trusted_proxies"`

	// Middleware flags
	EnableLogger   bool `mapstructure:"enable_logger"`
	EnableRecovery bool `mapstructure:"enable_recovery"`

	// Custom settings
	MaxMultipartMemory int64 `mapstructure:"max_multipart_memory"` // bytes
	BodySizeLimit      int64 `mapstructure:"body_size_limit"`      // bytes
	EnableMetrics      bool  `mapstructure:"enable_metrics"`
	EnableCORS         bool  `mapstructure:"enable_cors"`

	// Rate limiting
	RateLimit      int `mapstructure:"rate_limit"` // requests per second
	RateLimitBurst int `mapstructure:"rate_limit_burst"`
}

type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (c *DatabaseConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
	default:
		return ""
	}
}

type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiry     time.Duration `mapstructure:"jwt_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
	BcryptCost    int           `mapstructure:"bcrypt_cost"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

// func DefaultServerConfig() *ServerConfig {
// 	return &ServerConfig{
// 		Mode:               gin.DebugMode,
// 		Port:               "8080",
// 		Host:               "0.0.0.0",
// 		ReadTimeout:        30 * time.Second,
// 		WriteTimeout:       30 * time.Second,
// 		TrustedProxies:     []string{"127.0.0.1", "10.0.0.0/8"},
// 		EnableLogger:       true,
// 		EnableRecovery:     true,
// 		MaxMultipartMemory: 32 << 20, // 32 MB
// 		BodySizeLimit:      10 << 20, // 10 MB
// 		EnableMetrics:      false,
// 		EnableCORS:         false,
// 		RateLimit:          100,
// 		RateLimitBurst:     200,
// 	}
// }

func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	return nil
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
