package config

import (
    "log"
    "os"
    "strconv"
)

type DBConfig struct {
    Host           string
    Port           int
    User           string
    Password       string
    Name           string
    MigrationsPath string
}

type ServerConfig struct {
    Port string
}

type RabbitMQConfig struct {
    URL string
}

type Config struct {
    DB       DBConfig
    Server   ServerConfig
    RabbitMQ RabbitMQConfig
}

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() *Config {
    port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        log.Fatalf("Invalid DB_PORT: %v", err)
    }

    return &Config{
        DB: DBConfig{
            Host:           getEnv("DB_HOST", "localhost"),
            Port:           port,
            User:           getEnv("DB_USER", "postgres"),
            Password:       getEnv("DB_PASS", "postgres"),
            Name:           getEnv("DB_NAME", "aqua_sec_cloud_inventory"),
            MigrationsPath: getEnv("DB_MIGRATIONS_PATH", "migrations"),
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        RabbitMQ: RabbitMQConfig{
            URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
        },
    }
}

func getEnv(key, defaultVal string) string {
    val := os.Getenv(key)
    if val == "" {
        return defaultVal
    }
    return val
}
