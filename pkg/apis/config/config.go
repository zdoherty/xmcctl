package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Config contains configuration parameters for the entire program, including previously
// discovered devices and the preferred device to send commands to
type Config struct {
	// Selected is the name of a Device which should be used as the default to send commands to
	Selected string `yaml:"selected,omitempty"`

	// Devices is a list of all active devices
	Devices []RawDevice `yaml:"devices"`

	// Archive is a list of devices which were previously discovered but no longer used.
	Archive []RawDevice `yaml:"archive,omitempty"`
}

// NewConfigFromDefaults make a Config with default values
func NewConfigFromDefaults() *Config {
	c := &Config{
		Selected: "",
		Devices:  []RawDevice{},
		Archive:  []RawDevice{},
	}
	return c
}

// NewConfigFromBytes makes a Config from the passed byte slice
func NewConfigFromBytes(someBytes []byte) (*Config, error) {
	c := NewConfigFromDefaults()
	err := yaml.Unmarshal(someBytes, c)
	return c, err
}

// NewConfigFromFile makes a Config from the passed file
func NewConfigFromFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return NewConfigFromBytes(b)
}

// RawDevice contains information for a specific transponder in unparsed form.
type RawDevice struct {
	Name           string `yaml:"name"`
	Model          string `yaml:"model"`
	IP             string `yaml:"ip"`
	ControlVersion string `yaml:"control-version,omitempty"`
	ControlPort    int    `yaml:"control-port,omitempty"`
	NotifyPort     int    `yaml:"notify-port,omitempty"`
	InfoPort       int    `yaml:"info-port,omitempty"`
	SetupPort      int    `yaml:"setup-port,omitempty"`
}
