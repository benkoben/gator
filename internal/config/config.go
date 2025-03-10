package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configName            = ".gatorconfig.json"
	configFilePermissions = 0644
)

type configWriter interface {
	write() error
}

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`

	// unexported fields
	configWriter configWriter // allow dependency injection to enable unit testing
}

func (c Config) String() string {
	return fmt.Sprintf("CurrentUsername: %s, DbUrl: %s", c.CurrentUsername, c.DbUrl)
}

/*
Writes the current_user_name field in the .gatorconfig file. An error is returned if the write operation fails.
*/
func (c *Config) SetUser(name string) error {
	c.CurrentUsername = name

	if c.configWriter == nil {
		// Reference to c itself so that the write method can be utilized
		c.configWriter = c
	}

	if err := c.configWriter.write(); err != nil {
		return err
	}

	return nil
}

/*
Writes a marshalled Config object to the configuration
*/
func (c Config) write() error {
	configFile, err := getConfigFilepath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal config: %s", err)
	}

	if err := os.WriteFile(configFile, data, configFilePermissions); err != nil {
		return fmt.Errorf("could not write to config file: %s", err)
	}

	return nil
}

/*
Reads .gatorconfig from the current user's home directory
*/
func Read() (*Config, error) {
	configFile, err := getConfigFilepath()
	if err != nil {
		return nil, err
	}

	// The file does not exist, therefore quickly return
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s: %s", configFile, os.ErrNotExist)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", configFile, err)
	}

	cfg := Config{}

	buf := bytes.NewBuffer(data)
	decoder := json.NewDecoder(buf)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode file contents: %s", err)
	}

	return &cfg, nil
}

/*
returns the filepath for the config file assuming the file is located in the
user's homedirectory. If for some reason the current user's homedir is unknown an error is returned.
*/
func getConfigFilepath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %s", err)
	}

	configPath := filepath.Join(homeDir, configName)

	return configPath, nil
}
