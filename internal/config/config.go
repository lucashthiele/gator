package config

import (
	"encoding/json"
	"os"
	"path"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, configFileName), nil
}

func write(data []byte) error {
	configFilePath, err := getFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Read() (*Config, error) {
	configFilePath, err := getFilePath()
	if err != nil {
		return &Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return &Config{}, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = write(data)
	if err != nil {
		return err
	}

	return nil
}
