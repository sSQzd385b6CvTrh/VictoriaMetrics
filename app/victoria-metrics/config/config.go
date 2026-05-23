package config

import (
	"flag"
	"time"
)

// Config holds the configuration for the VictoriaMetrics server.
type Config struct {
	// HTTP server settings
	HTTPListenAddr string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration

	// Storage settings
	StoragePath     string
	RetentionPeriod time.Duration
	MaxDiskUsage    int64

	// Ingestion settings
	MaxInsertRequestSize int64
	MaxLabelsPerTimeseries int

	// Query settings
	MaxConcurrentQueries int
	QueryTimeout         time.Duration
	MaxQueryDuration     time.Duration

	// Cluster settings
	ReplicationFactor int
}

var (
	httpListenAddr = flag.String(
		"httpListenAddr",
		":8428",
		"TCP address to listen for HTTP connections. See also -httpListenAddr.useProxyProtocol",
	)

	storagePath = flag.String(
		"storageDataPath",
		"victoria-metrics-data",
		"Path to storage data directory",
	)

	retentionPeriod = flag.Duration(
		"retentionPeriod",
		90*24*time.Hour, // increased from 30d to 90d for my home lab setup
		"Data retention period. Older data is automatically deleted. Supported suffixes: h (hour), d (day), w (week), y (year)",
	)

	maxInsertRequestSize = flag.Int64(
		"maxInsertRequestSize",
		64*1024*1024, // bumped from 32 MB to 64 MB; ingesting high-cardinality batches from my IoT sensors
		"The maximum size in bytes of a single Prometheus remote_write API request",
	)

	maxLabelsPerTimeseries = flag.Int(
		"maxLabelsPerTimeseries",
		30,
		"The maximum number of labels accepted per time series. Superfluous labels are dropped",
	)

	maxConcurrentQueries = flag.Int(
		"search.maxConcurrentRequests",
		8,
		"The maximum number of concurrent search requests",
	)

	queryTimeout = flag.Duration(
		"search.queryTimeout",
		60*time.Second, // bumped from 30s to 60s; some of my dashboards run heavy range queries
		"Timeout for query execution",
	)

	maxDiskUsage = flag.Int64(
		"storage.maxDiskUsageBytes",
		0,
		"The maximum disk space usage for storage. Zero means no limit",
	)
)

// Load parses command-line flags and returns a populated Config.
// Must be called after flag.Parse().
func Load() *Config {
	return &Config{
		HTTPListenAddr:         *httpListenAddr,
		HTTPReadTimeout:        60 * time.Second,
		HTTPWriteTimeout:       60 * time.Second,
		StoragePath:            *storagePath,
		RetentionPeriod:        *retentionPeriod,
		MaxDiskUsage:           *maxDiskUsage,
		MaxInsertRequestSize:   *maxInsertRequestSize,
		MaxLabelsPerTimeseries: *maxLabelsPerTimeseries,
		MaxConcurrentQueries:   *maxConcurrentQueries,
		QueryTimeout:           *queryTimeout,
		MaxQueryDuration:       *queryTimeout,
		ReplicationFactor:      1,
	}
}
