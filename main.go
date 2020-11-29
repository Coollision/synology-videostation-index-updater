package main

import (
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"syscall"
	"time"

	"synology-videostation-reindexer/api"
	"synology-videostation-reindexer/metrics"
	"synology-videostation-reindexer/synology"
	synoConf "synology-videostation-reindexer/synology/config"

	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/toml"
	"github.com/sirupsen/logrus"
)

// Version is filled in by makefile when building (only in the main...)
var version string

//Config the large config struct that contains the whole app config
type Config struct {
	LogLevel       string `cfgDefault:"DEBUG"`
	ServerConfig   api.Config
	MetricsConfig  metrics.Config
	SynologyConfig synoConf.Config
}

func main() {
	logrus.Infof("version=%s", version)

	cfg := &Config{}
	goconfig.Path = "./"
	goconfig.File = "test.toml"
	err := goconfig.Parse(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	Log(cfg)
	customLogging(cfg.LogLevel)

	// Start Metrics server
	metrics.ServeMetrics(&cfg.MetricsConfig)

	syno := synology.NewVideoRequests(&cfg.SynologyConfig)

	srv := api.NewServer(&cfg.ServerConfig, syno)
	go srv.Start()
	defer srv.Stop()

	// Graceful shutdown
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM) // Kubernetes shutdown code
	signal.Notify(gracefulStop, syscall.SIGINT)  // CTRL + C

	// Block until we receive a shutdown signal
	<-gracefulStop

	// Stop all background work, previously defined with defers
	logrus.Info("Received a quit signal. Stopping background work now...")

	// Shut down
	defer shutdown()
}

func shutdown() {
	logrus.Info("Shutting down")
	os.Exit(0)
}

func customLogging(logLevel string) {
	// Setup logger
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Failed to parse log level. %v", err)
	}

	logrus.SetLevel(lvl)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = time.RFC3339
	customFormatter.DisableColors = true
	customFormatter.FullTimestamp = true
	customFormatter.FieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "time",
		logrus.FieldKeyLevel: "lvl",
		logrus.FieldKeyMsg:   "msg",
	}
	// get mesage as last value
	msgIsLastValue := func(s []string) {
		sort.Slice(s, func(i, j int) bool {
			if s[j] == "msg" {
				return true
			}
			return false
		})
	}
	customFormatter.SortingFunc = msgIsLastValue
	logrus.SetFormatter(customFormatter)
}

// logReflectValue recursively logs the given struct value
func logReflectValue(v reflect.Value, level int, name string) {
	t := v.Type()

	logrus.Infof("%v%v:", strings.Repeat("  ", level), name)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Log nested config
		if field.Type.Kind() == reflect.Struct {
			logReflectValue(reflect.ValueOf(value), level+1, field.Name)
			continue
		}

		// Hide sensitive data
		secret := field.Tag.Get("secret")
		if secret == "true" {
			value = "[REDACTED]"
		}
		logrus.Infof("%v%v (%v): %v", strings.Repeat("  ", level+1), field.Name, field.Type, value)
	}
}

// Log a config without exposing secrets.
// Takes a pointer to a config struct
// Secrets are marked with the tag `secret:"true"`
func Log(cfg interface{}) {
	rootElem := reflect.ValueOf(cfg).Elem()
	logReflectValue(rootElem, 0, "configuration")
	logrus.Infof("---------------------------------")
}
