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
	retentionPeriod = flag.Int("retentionPeriod", 1, "Retention period in months for the stored metrics. "+
		"Older data is automatically deleted")

	// storageDataPath is the path to the directory for storing data.
	storageDataPath = flag.String("storageDataPath", "victoria-metrics-data", "+
		"Path to storage data directory")

	// maxInsertRequestSize is the maximum size of a single insert request in bytes.
	maxInsertRequestSize = flag.Int("maxInsertRequestSize", 32*1024*1024, "+
		"The maximum size in bytes of a single insert request")

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

	// Ensure storage directory exists.
	if err := os.MkdirAll(*storageDataPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create storage data directory %q: %s\n", *storageDataPath, err)
		os.Exit(1)
	}

	// Set up HTTP request router.
	router := newRouter()

	// Start HTTP server.
	srv := &fasthttp.Server{
		Handler:            router,
		Name:               "VictoriaMetrics",
		ReadTimeout:        60 * time.Second,
		WriteTimeout:       60 * time.Second,
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
			handleMetrics(ctx)
		case "/health":
			handleHealth(ctx)
		default:
			ctx.Error(fmt.Sprintf("unsupported path: %s", path), http.StatusNotFound)
		}
	}
}

// handleHealth responds to health check requests.
func handleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprint(ctx, "OK")
}

// handleMetrics exposes internal metrics in Prometheus format.
func handleMetrics(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: write actual internal metrics
	fmt.Fprintf(ctx, "# VictoriaMetrics internal metrics\n")
}

// handleQuery handles instant PromQL queries.
func handleQuery(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: implement PromQL instant query
	fmt.Fprint(ctx, `{"status":"success","data":{"resultType":"vector","result":[]}}`)
}

// handleQueryRange handles range PromQL queries.
func handleQueryRange(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: implement PromQL range query
	fmt.Fprint(ctx, `{"status":"success","data":{"resultType":"matrix","result":[]}}`)
}

// handleSeries handles series metadata queries.
func handleSeries(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: implement series lookup
	fmt.Fprint(ctx, `{"status":"success","data":[]}`)
}

// handleLabels handles label names queries.
func handleLabels(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: implement label names lookup
	fmt.Fprint(ctx, `{"status":"success","data":[]}`)
}

// handleWrite handles remote write requests (Prometheus remote write protocol).
func handleWrite(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("only POST method is supported for /api/v1/write", http.StatusMethodNotAllowed)
		return
	}
	// TODO: implement Prometheus remote write ingestion
	if *logNewSeries {
		fmt.Println("Received write request")
	}
	ctx.SetStatusCode(http.StatusNoContent)
}
