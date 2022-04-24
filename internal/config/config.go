package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

    "github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("Error reading yaml file")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing yaml file")
	}
	return &c
}

func (c *Config) GetCoreByName(name string) *Core {
	for _, v := range c.Cores {
		if v.Name == name {
			return &v
		}
	}
	log.Fatal().Str("name", name).Msg("Unknown core")
	return nil
}
