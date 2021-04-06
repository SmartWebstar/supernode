package fileserver

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/supernode/internal/common"

	"golang.org/x/sync/errgroup"
)

type P2PServer struct {
	config *Config
	dht    *DHT
}

func New(cfg *Config) *P2PServer {
	return &P2PServer{config: cfg}
}

func (s *P2PServer) Start(ctx context.Context, app *common.Application) func() error {

	var bootstrapNodes []*NetworkNode

	for _, seed := range s.config.Seeds {
		if seed.Hostname != "" || seed.Port != "" {
			bootstrapNode := NewNetworkNode(seed.Hostname, seed.Port)
			bootstrapNodes = append(bootstrapNodes, bootstrapNode)
		}
	}

	var err error
	s.dht, err = NewDHT(&MemoryStore{}, &Options{
		BootstrapNodes: bootstrapNodes,
		IP:             s.config.Hostname,
		Port:           s.config.Port,
		UseStun:        s.config.Stun,
	})

	return app.CreateServer(ctx, "p2p_node",
		//initServer
		func(ctx context.Context) error {
			log.Infoln("p2p_node - Opening socket...")
			if err = s.dht.CreateSocket(); err != nil {
				return fmt.Errorf("p2p_node - error openning Socket for p2p server: %s", err)
			}
			log.Infoln("p2p_node - Socket opened")
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			eg, _ := errgroup.WithContext(ctx)

			eg.Go(func() error {
				log.Infoln("p2p_node is listening on " + s.dht.GetNetworkAddr())
				if err = s.dht.Listen(); err != nil && err.Error() != "closed" {
					return fmt.Errorf("p2p_node - error running p2p server: %s", err)
				}
				return nil
			})
			eg.Go(func() error {
				if len(bootstrapNodes) > 0 {
					log.Infoln("p2p_node - bootstrapping")
					if err = s.dht.Bootstrap(); err != nil {
						return fmt.Errorf("p2p_node - error bootstrapping p2p server: %s", err)
					}
					log.Infoln("p2p_node - bootstrapping done")
				}
				return nil
			})

			log.Infoln("p2p_node started")

			if err := eg.Wait(); err != nil {
				return fmt.Errorf("p2p_node - error in p2p server: %s", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			if err := s.dht.Disconnect(); err != nil {
				return fmt.Errorf("p2p_node - error stopping p2p server: %s", err)
			}
			return nil
		})
}
