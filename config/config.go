package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string         `yaml:"address"`
	PG      PostgresConfig `yaml:"postgres"`
	Tokens  TokensConfig   `yaml:"tokens"`
}

type TokensConfig struct {
	Secret     string `yaml:"secret"`
	RefreshExp int    `yaml:"refreshExp"`
	AccessExp  int    `yaml:"accessExp"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type Flags struct {
	ConfigPath string
}

func ParseFlags() Flags {
	processorCfgPath := flag.String("config", "", "Path to service cfg")
	flag.Parse()
	return Flags{
		ConfigPath: *processorCfgPath,
	}
}

func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}
