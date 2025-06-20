package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/RuriYS/RePorter/types"
	"gopkg.in/yaml.v3"
)

const default_config = `
server:
  host: 0.0.0.0
  port: 42000
  ttl: 3h
  allowed_ips: 
   - 0.0.0.0
  allowed_ports:
   - 8080
   - 8443

client:
  host: 127.0.0.1
  port: 42000
  broadcast_interval: 5m
  timeout: 3s
  ports:
   - 22
  whitelist_mode: false # forward ports in the list if true, otherwise don't (default: false)
`

var (
	config	*types.Config
)

func LoadConfig(path string) (err error) {
	if len(path) == 0 {
		path = "/etc/RePort/config.yml"
	}

	slog.Debug("[Config] loading config", "path", path)
	
	_, err = os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("[Config] config does not exist", "path", path, "error", err)
			createConfig(path)
		}
	}
	
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	c := &types.Config{}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return err
	}

	config = c
	slog.Debug("[Config] found config", "config", c)
	return nil
}

func GetConfig() (c *types.Config) {
	if config == nil {
		slog.Error("[Config] config not initialized", "config", c)
		os.Exit(1)
	}

	return config
}

func createConfig(path string) {
	slog.Debug("[Config] creating directory", "path", path)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("[Config] failed to create directory", "error", err.Error())
		os.Exit(1)
	}
	slog.Debug("[Config] writing default config", "path", path, "data", default_config)
	err = os.WriteFile(path, []byte(default_config), 0755)
	if err != nil {
		slog.Error("[Config] failed to create config", "error", err.Error())
		os.Exit(1)
	}
}
