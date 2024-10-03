package configuration

type Configuration struct {
	// Filesystem - filesystem configurations
	Filesystem Filesystem `yaml:"filesystem"`

	// Interfaces - network interfaces configurations
	Interfaces []Interface `yaml:"interfaces"`

	// Listeners - listeners configurations
	Listeners []Listener `yaml:"listeners"`
}

type Interface struct {
	// Name - interface name
	Name string `yaml:"name"`

	// IP - interface IP address
	Ip string `yaml:"ip"`
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

	// Interface - name of the interface to listen on
	Interface string `yaml:"interface"`

	// Socket - listener socket
	Socket Socket `yaml:"socket"`
}

type Socket struct {
	// Port - socket port
	Port string `yaml:"port"`

	// Protocol - socket protocol
	Protocol string `yaml:"protocol"`
}
