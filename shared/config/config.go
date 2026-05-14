package config

import "os"

// Config holds all configuration values loaded from environment variables.
type Config struct {
	AppEnv                string
	DatabaseURL           string
	RedisAddr             string
	JWTSecret             string
	AuthServicePort       string
	APIGatewayPort        string
	RunServicePort        string
	SubmissionServicePort string
	RepoServicePort       string
	PistonBaseURL         string
	GitHubAppID           string
	GitHubPrivateKeyPath  string
	GitHubOrgName         string
}

// defaults defines the fallback values for configuration keys.
var defaults = map[string]string{
	"APP_ENV":                 "development",
	"DATABASE_URL":            "postgres://innogen:innogen@localhost:5432/innogen?sslmode=disable",
	"REDIS_ADDR":              "localhost:6379",
	"JWT_SECRET":              "innogen-dev-secret",
	"AUTH_SERVICE_PORT":       "8081",
	"API_GATEWAY_PORT":        "8080",
	"RUN_SERVICE_PORT":        "8082",
	"SUBMISSION_SERVICE_PORT": "8083",
	"REPO_SERVICE_PORT":       "8084",
	"PISTON_BASE_URL":         "http://localhost:2000",
	"GITHUB_APP_ID":           "",
	"GITHUB_PRIVATE_KEY_PATH": "",
	"GITHUB_ORG_NAME":         "",
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
		AppEnv:                getEnv("APP_ENV"),
		DatabaseURL:           getEnv("DATABASE_URL"),
		RedisAddr:             getEnv("REDIS_ADDR"),
		JWTSecret:             getEnv("JWT_SECRET"),
		AuthServicePort:       getEnv("AUTH_SERVICE_PORT"),
		APIGatewayPort:        getEnv("API_GATEWAY_PORT"),
		RunServicePort:        getEnv("RUN_SERVICE_PORT"),
		SubmissionServicePort: getEnv("SUBMISSION_SERVICE_PORT"),
		RepoServicePort:       getEnv("REPO_SERVICE_PORT"),
		PistonBaseURL:         getEnv("PISTON_BASE_URL"),
		GitHubAppID:           getEnv("GITHUB_APP_ID"),
		GitHubPrivateKeyPath:  getEnv("GITHUB_PRIVATE_KEY_PATH"),
		GitHubOrgName:         getEnv("GITHUB_ORG_NAME"),
	}
}
