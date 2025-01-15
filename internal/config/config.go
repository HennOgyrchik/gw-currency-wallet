package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Postgres  PostgresConfig `env:",prefix=PSQL_" json:",omitempty"`
	Web       WebConfig      `env:",prefix=WEB_" json:",omitempty"`
	Exchanger GRPCConfig     `env:",prefix=EXCHANGER_GRPC_SERVER_" json:",omitempty"`
	Auth      GRPCConfig     `env:",prefix=AUTH_GRPC_SERVER_" json:",omitempty"`
}

type PostgresConfig struct {
	Host        string `env:"HOST,default=localhost" json:",omitempty"`
	Port        int    `env:"PORT,default=5432" json:",omitempty"`
	DBName      string `env:"DB_NAME,default=postgres" json:",omitempty"`
	User        string `env:"USER,default=postgres" json:",omitempty"`
	Password    string `env:"PASSWORD,default=postgres" json:",omitempty"`
	SSLMode     string `env:"SSL_MODE,default=disable" json:",omitempty"`
	ConnTimeout int    `env:"CONN_TIMEOUT,default=5" json:",omitempty"`
}

type WebConfig struct {
	Host string `env:"HOST,default=localhost" json:",omitempty"`
	Port int    `env:"PORT,default=80" json:",omitempty"`
}

type GRPCConfig struct {
	Host string `env:"HOST,default=localhost" json:",omitempty"`
	Port int    `env:"PORT,default=9090" json:",omitempty"`
}

func (g GRPCConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func (p PostgresConfig) ConnectionURL() (string, error) {
	host := p.Host
	port := p.Port
	if port < 1 && port > 65536 {
		return "", fmt.Errorf("PSQL_PORT invalid")
	}
	host = host + ":" + strconv.Itoa(p.Port)

	urlBuilder := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   p.DBName,
	}

	if p.User == "" || p.Password == "" {
		return "", fmt.Errorf("PSQL_USER or PSQL_PASSWORD invalid")
	}
	urlBuilder.User = url.UserPassword(p.User, p.Password)

	q := urlBuilder.Query()
	connTimeout := p.ConnTimeout
	if connTimeout < 1 {
		return "", fmt.Errorf("PSQL_CONN_TIMEOUT invalid")
	}
	q.Add("connect_timeout", strconv.Itoa(p.ConnTimeout))

	if p.SSLMode != "disable" && p.SSLMode != "enable" {
		return "", fmt.Errorf("PSQL_SSL_MODE invalid")
	}
	q.Add("sslmode", p.SSLMode)

	urlBuilder.RawQuery = q.Encode()

	return urlBuilder.String(), nil
}

func (w WebConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%d", w.Host, w.Port)
}

func LoadConfig(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func New() Config {
	return Config{Postgres: PostgresConfig{
		Host:        getEnvAsString("PSQL_HOST", "localhost"),
		Port:        getEnvAsInt("PSQL_PORT", 5432),
		DBName:      getEnvAsString("PSQL_DB_NAME", "postgres"),
		User:        getEnvAsString("PSQL_USER", "postgres"),
		Password:    getEnvAsString("PSQL_PASSWORD", "postgres"),
		SSLMode:     getEnvAsString("PSQL_SSL_MODE", "disable"),
		ConnTimeout: getEnvAsInt("PSQL_CONN_TIMEOUT", 60),
	},
		Web: WebConfig{
			Host: getEnvAsString("WEB_HOST", "localhost"),
			Port: getEnvAsInt("WEB_PORT", 80),
		},
		Exchanger: GRPCConfig{
			Host: getEnvAsString("EXCHANGER_HOST", "localhost"),
			Port: getEnvAsInt("EXCHANGER_PORT", 9090),
		},
		Auth: GRPCConfig{
			Host: getEnvAsString("AUTHORIZER_HOST", "localhost"),
			Port: getEnvAsInt("AUTHORIZER_PORT", 9090),
		}}

}

func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
