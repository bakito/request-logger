package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Conf struct {
	Echo     []string    `yaml:"echo"`
	EchoBody []string    `yaml:"echoBody"`
	Replay   [] Response `yaml:"replay"`
}

type Response struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	ContentType string `yaml:"contentType"`
}

func GetConf() *Conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		return nil
	}
	c := &Conf{}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil
	}

	return c
}
