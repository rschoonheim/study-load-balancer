package configuration

type Configuration struct {
	// Filesystem - filesystem configurations
	Filesystem Filesystem `yaml:"filesystem"`

	// Listeners - listeners configurations
	Listeners []Listener `yaml:"listeners"`
}

type Filesystem struct {
	// Configurations - filesystem configurations
	Configurations FilesystemConfigurations `yaml:"configurations"`
}

type FilesystemConfigurations struct {
	// Path - path to the configuration file
	Path string `yaml:"path"`

	// Rescan - rescan interval in seconds
	Rescan string `yaml:"rescan"`
}

type Listener struct {
	// Name - listener name
	Name string `yaml:"name"`

	// Socket - listener socket
	Socket Socket `yaml:"socket"`
}

type Socket struct {
	// Address - socket address
	Address string `yaml:"address"`

	// Port - socket port
	Port string `yaml:"port"`

	// Protocol - socket protocol
	Protocol string `yaml:"protocol"`
}
