package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Echo    []string   `yaml:"echo,omitempty"`
	LogBody []string   `yaml:"logBody,omitempty"`
	Replay  []Response `yaml:"replay,omitempty"`
}

type Response struct {
	Path        string `yaml:"path"`
	Body        string `yaml:"body,omitempty"`
	BodyFile    string `yaml:"bodyFile,omitempty"`
	ContentType string `yaml:"contentType"`
}

func GetConf(configFile string) (*Conf, error) {

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := &Conf{}

	err = yaml.Unmarshal(yamlFile, c)

	return c, err
}
