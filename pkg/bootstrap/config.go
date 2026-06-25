package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config   Config
	configMu sync.RWMutex
)

type Config struct {
	Environment string
	ServiceName string
	Port        int

	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	GRPCTLSCertPath string
	GRPCTLSKeyPath  string

	RabbitMQURL        string
	RabbitMQExchange   string
	RabbitMQQueue      string
	RabbitMQRoutingKey string
	RabbitMQDLX        string
	RabbitMQDLQ        string
	RabbitMQPrefetch   int

	FileManagementGRPCAddr               string
	FileManagementGRPCTLSEnabled         bool
	FileManagementGRPCCACertPath         string
	FileManagementGRPCInsecureSkipVerify bool

	PaginationLimitDefault int32
	PaginationLimitMax     int32
}

func LoadConfig(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		_ = godotenv.Load(dotenvPath)
	}

	cfg := Config{
		Environment:                          getEnv("ENVIRONMENT", "development"),
		ServiceName:                          getEnv("SERVICE_NAME", "neuraclinic-records"),
		Port:                                 getEnvInt("PORT", 8000),
		DBHost:                               getEnv("DB_HOST", ""),
		DBPort:                               getEnvInt("DB_PORT", 5432),
		DBUser:                               getEnv("DB_USER", ""),
		DBPass:                               getEnv("DB_PASS", ""),
		DBName:                               getEnv("DB_NAME", ""),
		DBSSLMode:                            getEnv("DB_SSLMODE", "disable"),
		GRPCTLSCertPath:                      getEnv("GRPC_TLS_CERT_PATH", ""),
		GRPCTLSKeyPath:                       getEnv("GRPC_TLS_KEY_PATH", ""),
		RabbitMQURL:                          getEnv("RABBITMQ_URL", ""),
		RabbitMQExchange:                     getEnv("RABBITMQ_EXCHANGE", "neuraclinic.events"),
		RabbitMQQueue:                        getEnv("RABBITMQ_QUEUE", "records.file_status_changed.v1"),
		RabbitMQRoutingKey:                   getEnv("RABBITMQ_ROUTING_KEY", "file.record.status_changed.v1"),
		RabbitMQDLX:                          getEnv("RABBITMQ_DLX", "neuraclinic.records.dlx"),
		RabbitMQDLQ:                          getEnv("RABBITMQ_DLQ", "records.file_status_changed.v1.dlq"),
		RabbitMQPrefetch:                     getEnvInt("RABBITMQ_PREFETCH", 10),
		FileManagementGRPCAddr:               getEnv("FILE_MANAGEMENT_GRPC_ADDR", ""),
		FileManagementGRPCTLSEnabled:         getEnvBool("FILE_MANAGEMENT_GRPC_TLS_ENABLED", true),
		FileManagementGRPCCACertPath:         getEnv("FILE_MANAGEMENT_GRPC_CA_CERT_PATH", ""),
		FileManagementGRPCInsecureSkipVerify: getEnvBool("FILE_MANAGEMENT_GRPC_INSECURE_SKIP_VERIFY", false),
		PaginationLimitDefault:               int32(getEnvInt("PAGINATION_LIMIT_DEFAULT", 10)),
		PaginationLimitMax:                   int32(getEnvInt("PAGINATION_LIMIT_MAX", 100)),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	setConfig(cfg)
	return cfg, nil
}

func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

func setConfig(cfg Config) {
	configMu.Lock()
	config = cfg
	configMu.Unlock()
}

func (c Config) Validate() error {
	required := map[string]string{
		"DB_HOST":                   c.DBHost,
		"DB_USER":                   c.DBUser,
		"DB_PASS":                   c.DBPass,
		"DB_NAME":                   c.DBName,
		"GRPC_TLS_CERT_PATH":        c.GRPCTLSCertPath,
		"GRPC_TLS_KEY_PATH":         c.GRPCTLSKeyPath,
		"RABBITMQ_URL":              c.RabbitMQURL,
		"FILE_MANAGEMENT_GRPC_ADDR": c.FileManagementGRPCAddr,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required config key: %s", key)
		}
	}

	if c.Port <= 0 {
		return fmt.Errorf("PORT must be greater than zero")
	}
	if c.DBPort <= 0 {
		return fmt.Errorf("DB_PORT must be greater than zero")
	}
	if c.PaginationLimitDefault <= 0 {
		return fmt.Errorf("PAGINATION_LIMIT_DEFAULT must be greater than zero")
	}
	if c.PaginationLimitMax < c.PaginationLimitDefault {
		return fmt.Errorf("PAGINATION_LIMIT_MAX must be greater than or equal to PAGINATION_LIMIT_DEFAULT")
	}
	if c.RabbitMQPrefetch <= 0 {
		return fmt.Errorf("RABBITMQ_PREFETCH must be greater than zero")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
