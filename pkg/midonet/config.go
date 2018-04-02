package midonet

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/yamt/midonet-kubernetes/pkg/config"
)

type Config struct {
	API				string
	ClusterRouter	uuid.UUID
}

func NewConfigFromEnvConfig(config *config.Config) *Config {
	router, err := uuid.Parse(config.ClusterRouter)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse cluster router")
	}
	return &Config{
		API: config.MidoNetAPI,
		ClusterRouter: router,
	}
}