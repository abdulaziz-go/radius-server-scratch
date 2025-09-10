package config

import (
	"os"
	"radius-server/src/common/logger"
	"strconv"
	"strings"

	typeUtil "radius-server/src/utils/type"

	"github.com/joho/godotenv"
)

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

type CorsConfig struct {
	Origins     string
	Methods     string
	Headers     string
	Credentials bool
}

type SecurityConfig struct {
	XApiKey string
}

type RadiusServerConfig struct {
	AccessHandlerServerPort     int
	AccountingHandlerServerPort int
	CoaHandlerServerPort        int
}

type Config struct {
	AppName      string
	AppHost      string
	AppLang      string
	AppVersion   string
	IsDebug      bool
	IsLocal      bool
	ServerPort   int
	Database     DbConfig
	Security     SecurityConfig
	Cors         CorsConfig
	RadiusServer RadiusServerConfig
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
	serverPort := getEnvAsInt("HTTP_SERVER_PORT", typeUtil.Int(8080), typeUtil.Int(0), typeUtil.Int(6666665))

	dbDns := getEnvAsString("DB_DNS", typeUtil.String("postgres://postgres:postgres@localhost:5532/postgres?sslmode=disable"))
	dbConnectionMaxNumber := getEnvAsInt("DB_CONNECTION_MAX_NUMBER", typeUtil.Int(10), nil, nil)
	dbConnectionOpenMaxNumber := getEnvAsInt("DB_CONNECTION_OPEN_MAX_NUMBER", typeUtil.Int(100), nil, nil)
	dbConnectionMaxLifetimeSec := getEnvAsInt("DB_CONNECTION_MAX_LIFETIME_SEC", typeUtil.Int(3600), nil, nil)
	dbConnectionMaxIdleTimeSec := getEnvAsInt("DB_CONNECTION_MAX_IDLE_TIME_SEC", typeUtil.Int(300), nil, nil)
	dbBatchInsertSize := getEnvAsInt("DB_BATCH_INSERT_SIZE", typeUtil.Int(1000), typeUtil.Int(1), nil)
	dbAutoRunMigration := getEnvAsBool("DB_AUTO_RUN_MIGRATION", typeUtil.Bool(true))
	dbLogger := getEnvAsBool("DB_LOGGER", typeUtil.Bool(true))

	securityXApiKey := getEnvAsString("SECURITY_X_API_KEY", typeUtil.String(""))
	if securityXApiKey == "" {
		logger.Logger.Fatal().Msg("Security SECURITY_X_API_KEY  is required.")
	}

	radiusAccessHanlderServerPort := getEnvAsInt("ACCESS_HANDLER_SERVER_PORT", typeUtil.Int(1812), typeUtil.Int(0), typeUtil.Int(6666665))
	radiusAccountingHanlderServerPort := getEnvAsInt("ACCOUNTING_HANDLER_SERVER_PORT", typeUtil.Int(1813), typeUtil.Int(0), typeUtil.Int(6666665))
	radiusCoaHandlerServerPort := getEnvAsInt("COA_HANDLER_SERVER_PORT", typeUtil.Int(1812), typeUtil.Int(0), typeUtil.Int(6666665))

	AppConfig = &Config{
		AppName:    appName,
		AppHost:    appHost,
		AppVersion: appVersion,
		AppLang:    appLang,
		IsDebug:    isDebug,
		IsLocal:    isLocal,
		ServerPort: serverPort,
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
		Security: SecurityConfig{
			XApiKey: securityXApiKey,
		},
		Cors: CorsConfig{
			Origins:     "*",
			Methods:     "GET,POST,PUT,DELETE,OPTIONS",
			Headers:     "Origin, Content-Type, Accept, Authorization",
			Credentials: false,
		},
		RadiusServer: RadiusServerConfig{
			AccessHandlerServerPort:     radiusAccessHanlderServerPort,
			AccountingHandlerServerPort: radiusAccountingHanlderServerPort,
			CoaHandlerServerPort:        radiusCoaHandlerServerPort,
		},
	}

}

func getEnvAsString(key string, defaultValue *string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msgf("Required environment variable %s is not set", key)
			return ""
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
			return false
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
			return []T{}
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
