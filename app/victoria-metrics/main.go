package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	// httpListenAddr is the address to listen for incoming HTTP requests.
	httpListenAddr = flag.String("httpListenAddr", ":8428", "TCP address to listen for incoming http requests")

	// retentionPeriod is the data retention period in months.
	// Increased default from 1 to 3 months for more useful local data retention.
	retentionPeriod = flag.Int("retentionPeriod", 3, "Retention period in months for the stored metrics. "+
		"Older data is automatically deleted")

	// storageDataPath is the path to the directory for storing data.
	storageDataPath = flag.String("storageDataPath", "victoria-metrics-data", "+
		\"Path to storage data directory\"")

	// maxInsertRequestSize is the maximum size of a single insert request in bytes.
	// Increased from 32MB to 64MB to handle larger batch writes from my local scrapers.
	maxInsertRequestSize = flag.Int("maxInsertRequestSize", 64*1024*1024, "+
		\"The maximum size in bytes of a single insert request\"")

	// logNewSeries enables logging of new time series.
	logNewSeries = flag.Bool("logNewSeries", false, "Whether to log new series. "+
		"This option is for debug purposes only. It can slow down ingestion performance")
)

func main() {
	// Parse command-line flags.
	flag.Parse()

	// Print startup information.
	fmt.Printf("Starting VictoriaMetrics at %s\n", *httpListenAddr)
	fmt.Printf("Storage data path: %s\n", *storageDataPath)
	fmt.Printf("Retention period: %d months\n", *retentionPeriod)
	fmt.Printf("Max insert request size: %d bytes\n", *maxInsertRequestSize)

	// Ensure storage directory exists.
	if err := os.MkdirAll(*storageDataPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create storage data directory %q: %s\n", *storageDataPath, err)
		os.Exit(1)
	}

	// Set up HTTP request router.
	router := newRouter()

	// Start HTTP server.
	// Increased read/write timeouts from 60s to 120s to avoid timeouts on
	// slow queries over large time ranges in my local Grafana dashboards.
	// Bumped IdleTimeout to 30s to keep persistent connections alive longer.
	srv := &fasthttp.Server{
		Handler:            router,
		Name:               "VictoriaMetrics",
		ReadTimeout:        120 * time.Second,
		WriteTimeout:       120 * time.Second,
		IdleTimeout:        30 * time.Second,
		MaxRequestBodySize: *maxInsertRequestSize,
	}

	fmt.Printf("VictoriaMetrics is ready to accept queries at http://%s\n", *httpListenAddr)
	if err := srv.ListenAndServe(*httpListenAddr); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot start HTTP server at %s: %s\n", *httpListenAddr, err)
		os.Exit(1)
	}
}

// newRouter creates and returns the HTTP request router.
func newRouter() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		switch path {
		case "/api/v1/query":
			handleQuery(ctx)
		case "/api/v1/query_range":
			handleQueryRange(ctx)
		case "/api/v1/series":
			handleSeries(ctx)
		case "/api/v1/labels":
			handleLabels(ctx)
		case "/api/v1/write":
			handleWrite(ctx)
		case "/metrics":
			hand
