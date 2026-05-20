package config

import "os"

// Config holds all configuration values loaded from environment variables.
type Config struct {
	AppEnv                         string
	DatabaseURL                    string
	RedisAddr                      string
	JWTSecret                      string
	AuthServicePort                string
	APIGatewayPort                 string
	RunServicePort                 string
	SubmissionServicePort          string
	RepoServicePort                string
	PistonBaseURL                  string
	GitHubAppID                    string
	GitHubPrivateKeyPath           string
	GitHubDefaultBranch            string
	GitHubAPIBaseURL               string
	GitHubWebhookSecret            string
	GitHubAppName                  string
	GitHubAppInstallURL            string
	GitHubOAuthClientID            string
	GitHubOAuthClientSecret        string
	GitHubOAuthRedirectURL         string
	GitHubOAuthFrontendRedirectURL string
	FrontendURL                    string
	AuthServiceURL                 string
	RunServiceURL                  string
	SubmissionServiceURL           string
	RepoServiceURL                 string
}

// defaults defines the fallback values for configuration keys.
var defaults = map[string]string{
	"APP_ENV":                            "development",
	"DATABASE_URL":                       "postgres://innogen:innogen@localhost:5432/innogen?sslmode=disable",
	"REDIS_ADDR":                         "localhost:6379",
	"JWT_SECRET":                         "innogen-dev-secret",
	"AUTH_SERVICE_PORT":                  "8081",
	"API_GATEWAY_PORT":                   "8080",
	"RUN_SERVICE_PORT":                   "8082",
	"SUBMISSION_SERVICE_PORT":            "8083",
	"REPO_SERVICE_PORT":                  "8084",
	"PISTON_BASE_URL":                    "http://localhost:2000",
	"GITHUB_APP_ID":                      "",
	"GITHUB_PRIVATE_KEY_PATH":            "",
	"GITHUB_DEFAULT_BRANCH":              "main",
	"GITHUB_API_BASE_URL":                "https://api.github.com",
	"GITHUB_WEBHOOK_SECRET":              "",
	"GITHUB_APP_NAME":                    "rinnogen",
	"GITHUB_APP_INSTALL_URL":             "https://github.com/apps/rinnogen/installations/new",
	"GITHUB_OAUTH_CLIENT_ID":             "",
	"GITHUB_OAUTH_CLIENT_SECRET":         "",
	"GITHUB_OAUTH_REDIRECT_URL":          "http://localhost:8080/github/oauth/callback",
	"GITHUB_OAUTH_FRONTEND_REDIRECT_URL": "http://localhost:5173/github",
	"FRONTEND_URL":                       "http://localhost:5173",
	"AUTH_SERVICE_URL":                   "http://localhost:8081",
	"RUN_SERVICE_URL":                    "http://localhost:8082",
	"SUBMISSION_SERVICE_URL":             "http://localhost:8083",
	"REPO_SERVICE_URL":                   "http://localhost:8084",
}

// getEnv returns the value of the environment variable named by key, or the
// default if the variable is not set or is empty.
func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaults[key]
}

// Load reads configuration from environment variables and returns a populated
// Config pointer. Missing or empty variables fall back to sensible defaults.
func Load() *Config {
	return &Config{
		AppEnv:                         getEnv("APP_ENV"),
		DatabaseURL:                    getEnv("DATABASE_URL"),
		RedisAddr:                      getEnv("REDIS_ADDR"),
		JWTSecret:                      getEnv("JWT_SECRET"),
		AuthServicePort:                getEnv("AUTH_SERVICE_PORT"),
		APIGatewayPort:                 getEnv("API_GATEWAY_PORT"),
		RunServicePort:                 getEnv("RUN_SERVICE_PORT"),
		SubmissionServicePort:          getEnv("SUBMISSION_SERVICE_PORT"),
		RepoServicePort:                getEnv("REPO_SERVICE_PORT"),
		PistonBaseURL:                  getEnv("PISTON_BASE_URL"),
		GitHubAppID:                    getEnv("GITHUB_APP_ID"),
		GitHubPrivateKeyPath:           getEnv("GITHUB_PRIVATE_KEY_PATH"),
		GitHubDefaultBranch:            getEnv("GITHUB_DEFAULT_BRANCH"),
		GitHubAPIBaseURL:               getEnv("GITHUB_API_BASE_URL"),
		GitHubWebhookSecret:            getEnv("GITHUB_WEBHOOK_SECRET"),
		GitHubAppName:                  getEnv("GITHUB_APP_NAME"),
		GitHubAppInstallURL:            getEnv("GITHUB_APP_INSTALL_URL"),
		GitHubOAuthClientID:            getEnv("GITHUB_OAUTH_CLIENT_ID"),
		GitHubOAuthClientSecret:        getEnv("GITHUB_OAUTH_CLIENT_SECRET"),
		GitHubOAuthRedirectURL:         getEnv("GITHUB_OAUTH_REDIRECT_URL"),
		GitHubOAuthFrontendRedirectURL: getEnv("GITHUB_OAUTH_FRONTEND_REDIRECT_URL"),
		FrontendURL:                    getEnv("FRONTEND_URL"),
		AuthServiceURL:                 getEnv("AUTH_SERVICE_URL"),
		RunServiceURL:                  getEnv("RUN_SERVICE_URL"),
		SubmissionServiceURL:           getEnv("SUBMISSION_SERVICE_URL"),
		RepoServiceURL:                 getEnv("REPO_SERVICE_URL"),
	}
}
