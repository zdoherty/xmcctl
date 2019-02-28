package config

type Config struct {
	// Selected is the name of a Device which should be used as the default to send commands to
	Selected string `yaml:"selected,omitempty"`

	// DiscoveryRequestPort is the port devices expect discovery requests to be sent to
	DiscoveryRequestPort int

	// DiscoveryResponsePort is the port to listen for discovery responses on
	DiscoveryResponsePort int

	// Devices is a list of all active devices
	Devices []Device

	// Archive is a list of devices which were previously discovered but no longer responsive.
	Archive []Device
}

type Device struct {
	Name        string
	Model       string
	IP          string
	ControlPort int
	NotifyPort  int
	InfoPort    int
}
