package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Outputs map[string]Output `yaml:"outputs"`
	Inputs  []string          `yaml:"inputs"`
	Games   []Game            `yaml:"games"`
	Cores   []Core            `yaml:"cores"`
	Mist    Mist              `yaml:"mist"`
	Ui      Ui                `yaml:"ui"`
}

type Output struct {
	Alias    string `yaml:"alias"`
	ScanCode int    `yaml:"scanCode"`
}

type Ui struct {
	Regions map[string]Region `yaml:"regions"`
}

type Textbox struct {
	X      int    `yaml:"x"`
	Y      int    `yaml:"y"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
	Align  string `yaml:"align"`
}

type Region struct {
	X       int     `yaml:"x"`
	Y       int     `yaml:"y"`
	Width   int     `yaml:"width"`
	Height  int     `yaml:"height"`
	Textbox Textbox `yaml:"textbox"`
}

type Game struct {
	Name         string    `yaml:"name"`
	Core         string    `yaml:"core"`
	Image        string    `yaml:"image"`
	Index        int       `yaml:"index"`
	Mappings     []Mapping `yaml:"mappings"`
	Disks        int       `yaml:"disks"`
	InitSequence []string  `yaml:"initSequence"`
}

type Mist struct {
	WaitAfterReset uint64 `yaml:"waitAfterReset"`
	ResetDuration  uint64 `yaml:"resetDuration"`
}

type Mapping struct {
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Autofire bool   `yaml:"autofire"`
	Text     string `yaml:"text"`
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

	c.ValidateGames()

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

func (c *Config) InputExists(name string) bool {
	for _, v := range c.Inputs {
		if v == name {
			return true
		}
	}
	return false
}

func (c *Config) OutputExists(name string) bool {
	for k := range c.Outputs {
		if k == name {
			return true
		}
	}
	return false
}

func (c *Config) ValidateGames() {

	for _, g := range c.Games {

		// Validate Mappings
		for _, m := range g.Mappings {
			if !c.InputExists(m.Input) {
				log.Fatal().Str("game", g.Name).Str("input", m.Input).Msg("Unknown input")
			}
			if !c.OutputExists(m.Output) {
				log.Fatal().Str("game", g.Name).Str("output", m.Output).Msg("Unknown output")
			}
		}

		// Validate Core
		if c.GetCoreByName(g.Core) == nil {
			log.Fatal().Str("core", g.Core).Msg("Unknown core")
		}

	}
}
