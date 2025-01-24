package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type StorageMode int

const (
	FileMode StorageMode = iota
	DatabaseMode
)

type Config struct {
	Env           string        `yaml:"env" env-default:"local"` // environment: local/dev/prod
	Storage       StorageConfig `yaml:"storage" env-required:"true"`
	GRPC          GRPCConfig    `yaml:"grpc"`
	MigrationPath string        `yaml:"migration_path"`
	TokenTTL      time.Duration `yaml:"token_ttl" env_default:"1h"`
}

type StorageConfig struct {
	StorageMode StorageMode
	FileStorage FileStorageConfig `yaml:"file_storage"`
	Database    DatabaseConfig    `yaml:"database"`
}

type FileStorageConfig struct {
	FilePath string `yaml:"file_path" env-default:"./storage/in_file_storage.json"`
}

type DatabaseConfig struct {
	URI string `yaml:"uri"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout"`
}

// MustLoad loads config file *.yaml from path in flag --config.
/* Default config path is ./config/local.yaml */
func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty, use --config=./path/to/config.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not exist by path " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("error reading config file " + err.Error())
	}

	if cfg.Storage.Database.URI == "" && cfg.Storage.FileStorage.FilePath == "" {
		panic("error getting Database URI or File Storage file path: check config.")
	}

	cfg.Storage.StorageMode = FileMode
	if cfg.Storage.Database.URI != "" {
		cfg.Storage.StorageMode = DatabaseMode
	}

	return &cfg
}

// fetchConfigPath fetches config path from flag --config.
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "./config/local.yaml", "path to config file .yaml")
	flag.Parse()

	return path
}
