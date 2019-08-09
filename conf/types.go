package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Echo     []string   `yaml:"echo"`
	EchoBody []string   `yaml:"echoBody"`
	Replay   []Response `yaml:"replay"`
}

type Response struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
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
