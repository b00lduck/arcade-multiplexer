package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Games []Game `yaml:"games"`
	Cores []Core `yaml:"cores"`
	Mist  Mist   `yaml:"mist"`
	Ui    Ui     `yaml:"ui"`
}

type Ui struct {
	Regions map[string]Region `yaml:"regions"`
}

type Region struct {
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type Game struct {
	Name     string    `yaml:"name"`
	Core     string    `yaml:"core"`
	Image    string    `yaml:"image"`
	Index    int       `yaml:"index"`
	Mappings []Mapping `yaml:"mappings"`
	Disks    int       `yaml:"disks"`
}

type Mist struct {
	WaitAfterReset uint64 `yaml:"waitAfterReset"`
	ResetDuration  uint64 `yaml:"resetDuration"`
}

type Mapping struct {
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Autofire bool   `yaml:"autofire"`
}

type Core struct {
	Name         string   `yaml:"name"`
	Enter        []string `yaml:"enter"`
	Exit         []string `yaml:"exit"`
	Load         []string `yaml:"load"`
	LoadSameCore []string `yaml:"loadSameCore"`
	Run          []string `yaml:"run"`
	BootSleep    uint64   `yaml:"bootSleep"`
	Speed1       uint64   `yaml:"speed1"`
	Speed2       uint64   `yaml:"speed2"`
}

func NewConfig() *Config {
	c := Config{}
	yamlFile, err := os.ReadFile("config.yml")
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
	if name == "none" {
		return nil
	}
	for _, v := range c.Cores {
		if v.Name == name {
			return &v
		}
	}
	log.Fatal().Str("name", name).Msg("Unknown core")
	return nil
}
