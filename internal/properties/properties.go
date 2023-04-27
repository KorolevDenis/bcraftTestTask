package properties

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	DBSettings struct {
		DBName     string `yaml:"DBName"`
		DBPort     string `yaml:"DBPort"`
		DBHost     string `yaml:"DBHost"`
		DBUsername string `yaml:"DBUsername"`
		DBPassword string `yaml:"DBPassword"`
	} `yaml:"DBSettings"`
	ProgramSettings struct {
		BindAddress   string `yaml:"bindAddress"`
		Auth          string `yaml:"auth"`
		TokenPassword string `yaml:"tokenPassword"`
	} `yaml:"ProgramSettings"`
}

func GetConfig() (*Config, error) {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("Error while open config file: " + err.Error())
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error while read config file: " + err.Error())
	}

	var properties Config
	err = yaml.Unmarshal(data, &properties)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshal config file: " + err.Error())
	}

	return &properties, nil
}
