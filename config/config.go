package config

import (
	"encoding/json"

	"github.com/pastelnetwork/supernode/internal/fileserver"
	"github.com/pastelnetwork/supernode/internal/restserver"
	"github.com/pastelnetwork/supernode/nats"
	"github.com/pastelnetwork/supernode/pastel"
)

// Config contains configuration of all components of the SuperNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel *pastel.Config     `mapstructure:"pastel" json:"pastel,omitempty"`
	Nats   *nats.Config       `mapstructure:"nats" json:"nats,omitempty"`
	REST   *restserver.Config `mapstructure:"rest" json:"rest,omitempty"`
	P2P    *fileserver.Config `mapstructure:"p2p" json:"p2p,omitempty"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		Main:   *NewMain(),
		Pastel: pastel.NewConfig(),
		Nats:   nats.NewConfig(),
		REST:   restserver.NewConfig(),
		P2P:    fileserver.NewConfig(),
	}
}
