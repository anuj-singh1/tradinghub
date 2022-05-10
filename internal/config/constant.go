package config

const (
	PostgresSource          = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	DbDriver                = "postgres"
	GIN_ENV_GLOBAL_INSTANCE = "globalinstance"
	EnvVariablePath         = "./"
	EnvFileName             = "app"
	EnvFileExtension        = "env"
	FyersBaseUrl            = "https://api.fyers.in/api/v2/"
	FyersDataUrl            = "https://api.fyers.in/data-rest/v2/"
	ValidateCodePath        = "validate-authcode"
	GenerateAuthCodePath    = "generate-authcode"
	QuotesPath              = "quotes/"
)
