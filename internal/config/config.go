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
    configName = ".gatorconfig"
)

type Config struct {
    DbUrl string `json:"db_url"`
    CurrentUsername string `json:"current_user_name"`
}

/*
    Writes the current_user_name field in the .gatorconfig file
*/
func (c Config)SetUser(){
}

/*
    Reads .gatorconfig from the current user's home directory
*/
func Read() (*Config, error) {
    configFile, err := getConfigFilepath() 
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(configFile)
    if err != nil {
        return nil, fmt.Errorf("could not read %s: %s", configFile, err)
    }

    cfg := Config{}

    buf := bytes.NewBuffer(data)
    decoder := json.NewDecoder(buf)
    if err  := decoder.Decode(&cfg); err != nil {
        return nil, fmt.Errorf("could not decode file contents: %s", err)
    }

    return &cfg, nil
}


/*
    returns the filepath for the config file assuming the file exists in the 
    user's homedirectory. If it does not exist an error is returned.
*/
func getConfigFilepath() (string, error){
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("could not get home directory: %s", err)
    }

    configPath := filepath.Join(homeDir, configName)
    if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
        return "", fmt.Errorf("%s: %s", configPath, os.ErrNotExist)
    }

    return configPath, nil
}
