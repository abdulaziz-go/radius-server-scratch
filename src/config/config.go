package config

import (
	"os"
	"radius-server/src/common/logger"
	"strconv"
	"strings"

	typeUtil "radius-server/src/utils/type"

	"github.com/joho/godotenv"
)

type RedisConnectionConfig struct {
	MaxNumber       int
	OpenMinNumber   int
	MaxLifetimeSec  int
	MaxIdleTimeSec  int
	DeleteBatchSize int
	ResponseLimit   int
}

type DbConnectionConfig struct {
	MaxNumber      int
	OpenMaxNumber  int
	MaxLifetimeSec int
	MaxIdleTimeSec int
}

type DbConfig struct {
	Dsn              string
	Connection       DbConnectionConfig
	BatchInsertSize  int
	AutoRunMigration bool
	Logger           bool
}

type ClickhouseConnectionConfig struct {
	MaxNumber      int
	OpenMaxNumber  int
	MaxLifetimeSec int
}
type ClickhouseConfig struct {
	Host             string
	Port             int
	DbName           string
	Username         string
	Password         string
	BatchInsertSize  int
	AutoRunMigration bool
	Connection       ClickhouseConnectionConfig
	MaxMemoryLimit   int64 // In GB
}

type ServerConfig struct {
	Id   string
	Port int
	Host string
}

type SecurityConfig struct {
	JwtKey            string
	JwtLifetimeSec    int
	BotJwtLifetimeSec int
	XApiKey           string
}

type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Enabled     bool
}

type CorsConfig struct {
	Origins     string
	Methods     string
	Headers     string
	Credentials bool
}

type RedisConfig struct {
	Host       string
	Connection RedisConnectionConfig
}

type CacheConfig struct {
	FlowLifetimeSec int
}

type WebsocketConfig struct {
	HeartbeatPingTimeSec int
	HeartbeatPongTimeSec int
	Logger               bool
}

type CronConfig struct {
	WebsocketStatsSchedule string
	FlowMetricsSchedule    string
}

type CookieConfig struct {
	Secure bool
}

type PrometheusConfig struct {
	Address string
	Points  int64
}

type FakeMetricsConfig struct {
	DoFakeMetrics bool
	MinApps       int
	MaxApps       int
	IpPool        string
	MinIps        int
	MaxIps        int
	MinTotalBytes int
	MaxTotalBytes int
	MinFlows      int
	MaxFlows      int
}

type Config struct {
	AppName    string
	AppHost    string
	AppLang    string
	AppVersion string
	IsDebug    bool
	IsLocal    bool
	Database   DbConfig
	Server     ServerConfig
	Security   SecurityConfig
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Fatal().Msgf("Loading .env file error. Error - %s", err.Error())
	}
	appName := getEnvAsString("APP_NAME", typeUtil.String("radius-server"))
	appHost := getEnvAsString("APP_HOST", typeUtil.String(""))
	appVersion := getEnvAsString("APP_VERSION", typeUtil.String("v1.0.0"))
	appLang := getEnvAsString("APP_LANG", typeUtil.String("en"))
	isDebug := getEnvAsBool("IS_DEBUG", typeUtil.Bool(true))
	isLocal := getEnvAsBool("IS_LOCAL", typeUtil.Bool(true))

	serverPort := getEnvAsInt("SERVER_PORT", typeUtil.Int(3000), nil, nil)
	serverHost := getEnvAsString("SERVER_HOST", typeUtil.String(""))

	dbDns := getEnvAsString("DB_DNS", typeUtil.String("postgres://postgres:postgres@localhost:5532/postgres?sslmode=disable"))
	dbConnectionMaxNumber := getEnvAsInt("DB_CONNECTION_MAX_NUMBER", typeUtil.Int(10), nil, nil)
	dbConnectionOpenMaxNumber := getEnvAsInt("DB_CONNECTION_OPEN_MAX_NUMBER", typeUtil.Int(100), nil, nil)
	dbConnectionMaxLifetimeSec := getEnvAsInt("DB_CONNECTION_MAX_LIFETIME_SEC", typeUtil.Int(3600), nil, nil)
	dbConnectionMaxIdleTimeSec := getEnvAsInt("DB_CONNECTION_MAX_IDLE_TIME_SEC", typeUtil.Int(300), nil, nil)
	dbBatchInsertSize := getEnvAsInt("DB_BATCH_INSERT_SIZE", typeUtil.Int(1000), typeUtil.Int(1), nil)
	dbAutoRunMigration := getEnvAsBool("DB_AUTO_RUN_MIGRATION", typeUtil.Bool(true))
	dbLogger := getEnvAsBool("DB_LOGGER", typeUtil.Bool(true))

	securityJwtKey := getEnvAsString("SECURITY_JWT_KEY", typeUtil.String(""))
	securityJwtLifetimeSec := getEnvAsInt("SECURITY_JWT_LIFETIME_SEC", typeUtil.Int(86400), nil, nil)  // 24 * 60 * 60
	botJwtLifetimeSec := getEnvAsInt("SECURITY_BOT_JWT_LIFETIME_SEC", typeUtil.Int(2592000), nil, nil) // 30 * 24 * 60 * 60
	securityXApiKey := getEnvAsString("SECURITY_X_API_KEY", typeUtil.String(""))
	if securityXApiKey == "" {
		logger.Logger.Fatal().Msg("Security SECURITY_X_API_KEY  is required.")
	}

	AppConfig = &Config{
		AppName:    appName,
		AppHost:    appHost,
		AppVersion: appVersion,
		AppLang:    appLang,
		IsDebug:    isDebug,
		IsLocal:    isLocal,

		Database: DbConfig{
			Dsn: dbDns,
			Connection: DbConnectionConfig{
				MaxNumber:      dbConnectionMaxNumber,
				OpenMaxNumber:  dbConnectionOpenMaxNumber,
				MaxLifetimeSec: dbConnectionMaxLifetimeSec,
				MaxIdleTimeSec: dbConnectionMaxIdleTimeSec,
			},
			BatchInsertSize:  dbBatchInsertSize,
			AutoRunMigration: dbAutoRunMigration,
			Logger:           dbLogger,
		},

		Server: ServerConfig{
			Port: serverPort,
			Host: serverHost,
		},
		Security: SecurityConfig{
			JwtKey:            securityJwtKey,
			JwtLifetimeSec:    securityJwtLifetimeSec,
			BotJwtLifetimeSec: botJwtLifetimeSec,
			XApiKey:           securityXApiKey,
		},
	}
}

type UserSeed struct {
	Username string
	Password string
	Role     string
}

func parseUsersSeed(key string) []UserSeed {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Logger.Warn().Msg("Users seed is not set")
		return []UserSeed{}
	}
	userEntries := strings.Split(value, ",")
	var users []UserSeed
	for _, userEntry := range userEntries {
		userEntry = strings.TrimSpace(userEntry)
		parts := strings.Split(userEntry, ":")
		if len(parts) != 3 {
			logger.Logger.Fatal().Msgf("Invalid user format in %s. Expected format: username:password:role, got: %s", key, userEntry)
		}
		users = append(users, UserSeed{
			Username: strings.TrimSpace(parts[0]),
			Password: strings.TrimSpace(parts[1]),
			Role:     strings.TrimSpace(parts[2]),
		})
	}
	return users
}

func getEnvAsString(key string, defaultValue *string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msgf("Required environment variable %s is not set", key)
		}
		return *defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue *int, minValue *int, maxValue *int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msgf("Required environment variable %s is not set", key)
		}
		if minValue != nil && *defaultValue < *minValue {
			logger.Logger.Fatal().Msgf("Default value for %s must be at least %d, got %d", key, *minValue, *defaultValue)
		}
		if maxValue != nil && *defaultValue > *maxValue {
			logger.Logger.Fatal().Msgf("Default value for %s must be at most %d, got %d", key, *maxValue, *defaultValue)
		}
		return *defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Logger.Fatal().Msgf("Environment variable %s must be an integer, got %s", key, value)
	}
	if minValue != nil && intValue < *minValue {
		logger.Logger.Fatal().Msgf("Environment variable %s must be at least %d, got %d", key, *minValue, intValue)
	}
	if maxValue != nil && intValue > *maxValue {
		logger.Logger.Fatal().Msgf("Environment variable %s must be at most %d, got %d", key, *maxValue, intValue)
	}
	return intValue
}

func getEnvAsBool(key string, defaultValue *bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msgf("Required environment variable %s is not set", key)
		}
		return *defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		logger.Logger.Fatal().Msgf("Environment variable %s must be a valid boolean, got %s", key, value)
	}
	return boolValue
}

func getEnvAsEnums[T ~string](key string, defaultValue *[]T, allowedValues []string) []T {
	rawValue, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msgf("Required environment variable %s is not set", key)
		}
		return *defaultValue
	}
	values := strings.Split(rawValue, ",")
	var result []T
	for _, v := range values {
		v = strings.TrimSpace(v)
		valid := false
		for _, allowed := range allowedValues {
			if v == allowed {
				result = append(result, T(v))
				valid = true
				break
			}
		}
		if !valid {
			logger.Logger.Fatal().Msgf("Invalid value '%s' in environment variable %s. Allowed values: %v", v, key, allowedValues)
		}
	}
	return result
}
