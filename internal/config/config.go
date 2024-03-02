package config

// Config defines the structure of the configuration for the server.
type Config struct {
	ServerPort     string         `json:"server_port"` //the server's port
	DatabaseConfig DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}
