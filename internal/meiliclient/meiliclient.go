package meiliclient

import (
	"germa66/internal/config"

	meilisearch "github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

type MeiliClient interface {
	HealthCheck() bool
}

type Service struct {
	index  meilisearch.IndexManager
	client meilisearch.ServiceManager
	conf   *config.Config
}

// New creates a new MeiliClient using the provided configuration,
// connects to the MeiliSearch instance and creates the index if it doesn't exist.
func New(conf *config.Config) *Service {
	log.Infof("Creating connection to Meilisearch on %s", conf.MeilisearchHost)
	client := meilisearch.New(conf.MeilisearchHost, meilisearch.WithAPIKey(conf.MeilisearchAPIKey))
	index := client.Index(conf.MeiliIndex)

	return &Service{
		client: client,
		index:  index,
		conf:   conf,
	}
}

// HealthCheck checks the health of the MeiliSearch instance.
func (mc *Service) HealthCheck() bool {
	return mc.client.IsHealthy()
}
