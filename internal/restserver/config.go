package restserver

const (
	defaultHostname = "localhost"
	defaultPort     = 4001
)

type Config struct {
	Hostname string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
