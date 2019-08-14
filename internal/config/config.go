package config

type Config struct {
	Games []Game `yaml:"games"`
}

type Game struct {
	Name     string    `yaml:"name"`
	Core     string    `yaml:"core"`
	Image    string    `yaml:"image"`
	Index    int       `yaml:"index"`
	Mappings []Mapping `yaml:"mappings"`
}

type Mapping struct {
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Autofire bool   `yaml:"autofire"`
}
