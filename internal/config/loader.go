package config

import (
	"fmt"
	"os"
)

func LoadConfig(configPath string) (*Config, error) {
	// Existing code to load the config file into the `config` variable

	// Retrieve MySQL username and password from environment variables
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")

	if user == "" || password == "" {
		return nil, fmt.Errorf("MySQL username or password not set in environment variables")
	}

	// Append the username and password to the DSN
	config.MySQLDSN = fmt.Sprintf("%s:%s@%s", user, password, config.MySQLDSN)

	return &config, nil
}
