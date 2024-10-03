package configuration

type Configuration struct {
	// Filesystem - filesystem configurations
	Filesystem Filesystem `yaml:"filesystem"`
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
