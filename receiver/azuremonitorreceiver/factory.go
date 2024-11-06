// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuremonitorreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver"

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver/internal/metadata"
)

const (
	defaultCollectionInterval = 10 * time.Second
	defaultCloud              = azureCloud
	defaultSplitByDimensions  = true
)

var errConfigNotAzureMonitor = errors.New("Config was not a Azure Monitor receiver config")

// NewFactory creates a new receiver factory
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}

func createDefaultConfig() component.Config {
	cfg := scraperhelper.NewDefaultControllerConfig()
	cfg.CollectionInterval = defaultCollectionInterval

	return &Config{
		ControllerConfig:                  cfg,
		MetricsBuilderConfig:              metadata.DefaultMetricsBuilderConfig(),
		CacheResources:                    24 * 60 * 60,
		CacheResourcesDefinitions:         24 * 60 * 60,
		MaximumNumberOfMetricsInACall:     20,
		MaximumNumberOfRecordsPerResource: 10,
		MaximumNumberOfDimensionsInACall:  10,
		Services:                          monitorServices,
		Authentication:                    servicePrincipal,
		Cloud:                             defaultCloud,
		SplitByDimensions:                 to.Ptr(defaultSplitByDimensions),
	}
}

func createMetricsReceiver(_ context.Context, params receiver.Settings, rConf component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	cfg, ok := rConf.(*Config)
	if !ok {
		return nil, errConfigNotAzureMonitor
	}

	var scraper scraperhelper.Scraper
	var err error
	if cfg.UseBatchAPI {
		azureBatchScraper := newBatchScraper(cfg, params)
		scraper, err = scraperhelper.NewScraper(metadata.Type, azureBatchScraper.scrape, scraperhelper.WithStart(azureBatchScraper.start))
	} else {
		azureScraper := newScraper(cfg, params)
		scraper, err = scraperhelper.NewScraper(metadata.Type, azureScraper.scrape, scraperhelper.WithStart(azureScraper.start))
	}
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(&cfg.ControllerConfig, params, consumer, scraperhelper.AddScraper(scraper))
}
