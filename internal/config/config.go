package config

import (
	"io/ioutil"

	"github.com/tarent/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Games []Game `yaml:"games"`
	Cores []Core `yaml:"cores"`
}

type Game struct {
	Name     string    `yaml:"name"`
	Core     string    `yaml:"core"`
	Image    string    `yaml:"image"`
	Index    int       `yaml:"index"`
	Mappings []Mapping `yaml:"mappings"`
	Disks    int	   `yaml:"disks"`
}

type Mapping struct {
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Autofire bool   `yaml:"autofire"`
}

type Core struct {
	Name      string   `yaml:"name"`
	Enter     []string `yaml:"enter"`
	Exit      []string `yaml:"exit"`
	Load      []string `yaml:"load"`
	Run       []string `yaml:"run"`
	BootSleep uint64   `yaml:"bootSleep"`
	Speed1    uint64   `yaml:"speed1"`
	Speed2    uint64   `yaml:"speed2"`
}

func NewConfig() *Config {
	c := Config{}
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Error reading yaml file")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing yaml file")
	}
	return &c
}

func (c *Config) GetCoreByName(name string) *Core {
	for _, v := range c.Cores {
		if v.Name == name {
			return &v
		}
	}
	logrus.WithField("name", name).Fatal("Unknown core")
	return nil
}
