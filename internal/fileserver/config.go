package fileserver

const (
	defaultHostname = "localhost"
	defaultPort     = "4444"
)

type P2PSeeds struct {
	Hostname string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
}

type Config struct {
	Hostname string     `mapstructure:"hostname"`
	Port     string     `mapstructure:"port"`
	Stun     bool       `mapstructure:"stun"`
	Seeds    []P2PSeeds `mapstructure:"seeds"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
