package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var isRunning bool = false
var metrics *Metrics

//Metrics the main metricky stuff
type Metrics struct {
	cfg    *Config
	server *http.Server
}

//ServeMetrics starts to load metrics and reteruns a metrics object
func ServeMetrics(cfg *Config) {
	if isRunning {
		panic("cannot start to serve metrics twice")
	}
	isRunning = true

	metrics = &Metrics{
		cfg: cfg,
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	metrics.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", metrics.cfg.Port),
		WriteTimeout: time.Second *30,
		ReadTimeout:  time.Second *30,
		IdleTimeout:  time.Second * 30,
		Handler:      mux,
	}

	go func() {
		if err := metrics.server.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Errorf("Failed to start server: %v", err)  // not killing app if metrics is not running,
																	// if this is needed let kube do this with a probe
																	// otherwise metrics are not that important
		}
	}()

}

// DestroyMetrics stop serving metrics and clear all global vars in this package
func DestroyMetrics() {
	isRunning = false
	metrics = nil
}
