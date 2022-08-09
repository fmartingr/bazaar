package server

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

// readDotEnv reads the configuration from variables in a .env file (only for contributing)
func readDotEnv() map[string]string {
	file, err := os.Open(".env")
	if err != nil {
		return nil
	}
	defer file.Close()

	result := make(map[string]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		keyval := strings.SplitN(line, "=", 2)
		result[keyval[0]] = keyval[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

type ServerConfig struct {
	Hostname string `env:"HOSTNAME,required"`
	Http     struct {
		Enabled bool `env:"HTTP_ENABLED,default=True"`
		Port    int  `env:"HTTP_PORT,default=8080"`
	}
}

func ParseServerConfiguration(ctx context.Context) *ServerConfig {
	var cfg ServerConfig

	lookuper := envconfig.MultiLookuper(
		envconfig.MapLookuper(map[string]string{"HOSTNAME": os.Getenv("HOSTNAME")}),
		envconfig.MapLookuper(readDotEnv()),
		envconfig.PrefixLookuper("BAZAAR_", envconfig.OsLookuper()),
		envconfig.OsLookuper(),
	)
	if err := envconfig.ProcessWith(ctx, &cfg, lookuper); err != nil {
		log.Fatalf("Error parsing configuration: %s", err)
	}

	return &cfg
}
