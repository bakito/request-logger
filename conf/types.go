package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf type to define custom request mapping
type Conf struct {
	Echo    []string   `yaml:"echo,omitempty"`
	LogBody []LogBody  `yaml:"logBody,omitempty"`
	Replay  []Response `yaml:"replay,omitempty"`
}

// LogBody config type
type LogBody struct {
	Path       string `yaml:"path"`
	LineLength bool   `yaml:"lineLength"`
}

// Response config type
type Response struct {
	Path        string `yaml:"path"`
	Body        string `yaml:"body,omitempty"`
	BodyFile    string `yaml:"bodyFile,omitempty"`
	ContentType string `yaml:"contentType"`
}

// GetConf get the config from the given file
func GetConf(configFile string) (*Conf, error) {

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := &Conf{}

	err = yaml.Unmarshal(yamlFile, c)

	return c, err
}
