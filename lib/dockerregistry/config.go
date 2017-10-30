package dockerregistry

import (
	"code.uber.internal/infra/kraken/lib/dockerregistry/transfer"
	"code.uber.internal/infra/kraken/lib/store"
	docker "github.com/docker/distribution/configuration"
	"github.com/uber-go/tally"
)

// Config contains docker registry config, disable torrent flag, and tag deletion config
type Config struct {
	Docker         docker.Configuration `yaml:"docker"`
	DisableTorrent bool                 `yaml:"disable_torrent"`
	TagDir         string               `yaml:"tag_dir"`
	TagDeletion    TagDeletionConfig    `yaml:"tag_deletion"`
}

// TagDeletionConfig contains configuration to delete tags
type TagDeletionConfig struct {
	Enable bool `yaml:"enable"`
	// Interval for running tag deletion in seconds
	Interval int `yaml:"interval"`
	// Number of tags we keep for each repo
	RetentionCount int `yaml:"retention_count"`
	// Least number of seconds we keep tags for
	RetentionTime int `yaml:"retention_time"`
}

// CreateDockerConfig returns docker specified configuration
func (c *Config) CreateDockerConfig(name string, imageTransferer transfer.ImageTransferer, fileStore store.FileStore, stats tally.Scope) *docker.Configuration {
	c.Docker.Storage = docker.Storage{
		name: docker.Parameters{
			"config":     c,
			"transferer": imageTransferer,
			"store":      fileStore,
			"metrics":    stats,
		},
		// Redirect is enabled by default in docker registry.
		// We implement redirect on proxy level so we do not need this in storage driver for now.
		"redirect": docker.Parameters{
			"disable": true,
		},
	}
	return &c.Docker
}