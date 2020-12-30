package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bakito/request-logger/pkg/common"
	"github.com/bakito/request-logger/pkg/conf"
	"github.com/bakito/request-logger/pkg/handler"
	"github.com/bakito/request-logger/pkg/middleware"
	"github.com/bakito/request-logger/version"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	port             int
	countRequestRows bool
	disableLogger    bool
	metrics          bool
	configFile       string
	tlsKey           string
	tlsCert          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "request-logger",
	Version: version.Version,
	Short:   "A Simple webserver allowing to log incoming requests",
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *conf.Conf
		var err error
		if configFile != "" {
			config, err = conf.GetConf(configFile)
			if err != nil {
				return fmt.Errorf("error reading config %s: %v", configFile, err)
			}
		}

		functions := make(map[string]func(w http.ResponseWriter, r *http.Request))
		var paths []string
		r := router()
		if config != nil {
			for _, path := range config.Echo {
				functions[path] = handler.Echo
				paths = append(paths, path)
			}

			for _, lb := range config.LogBody {
				functions[lb.Path] = handler.ConfigLogBody(lb)
				paths = append(paths, lb.Path)
			}

			for _, resp := range config.Replay {
				functions[resp.Path] = handler.ConfigReplay(resp)
				paths = append(paths, resp.Path)
			}

			common.SortPaths(paths)

			log.Printf("Serving custom config from '%s'", configFile)
			for _, p := range paths {
				r.HandleFunc(p, functions[p])
			}

		} else {
			r.HandleFunc("/echo", handler.Echo)
			r.HandleFunc("/echo/{path:.*}", handler.Echo)

			r.HandleFunc("/body", handler.LogBody)
			r.HandleFunc("/body/{path:.*}", handler.LogBody)

			r.HandleFunc(`/code/{code:[2,4,5]\d\d}`, handler.ResponseCode)
			r.HandleFunc(`/code/{code:[2,4,5]\d\d}/{path:.*}`, handler.ResponseCode)

			r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}`, handler.RandomCode)
			r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}/{path:.*}`, handler.RandomCode)

			r.HandleFunc(`/random/sleep/{sleep:\d+}`, handler.RandomSleep)
			r.HandleFunc(`/random/sleep/{sleep:\d+}/{path:.*}`, handler.RandomSleep)

			r.HandleFunc(`/replay`, handler.Replay)
			r.HandleFunc(`/replay/{path:.*}`, handler.Replay)

			r.HandleFunc("/{path:.*}", handler.Void)
		}

		start(r)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Server port")
	rootCmd.PersistentFlags().BoolVar(&disableLogger, "disableLogger", false, "Disable the logger middleware")
	rootCmd.PersistentFlags().BoolVar(&metrics, "metrics", true, "Enable metrics")
	rootCmd.PersistentFlags().StringVar(&tlsKey, "tlsKey", "", "TLS key file")
	rootCmd.PersistentFlags().StringVar(&tlsCert, "tlsCert", "", "TLS cert file")

	rootCmd.Flags().BoolVar(&countRequestRows, "countRequestRows", true, "Enable or disable the request row count")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "The path of a config file")
}

func router() *mux.Router {
	r := mux.NewRouter()

	if metrics {
		r.Handle("/metrics", promhttp.Handler())
	}

	r.Use(middleware.CountRequests)

	if !disableLogger {
		r.Use(middleware.LogRequest)
	}
	if countRequestRows {
		r.Use(middleware.CountReqRows)
	}

	return r
}

func start(r *mux.Router) {
	if withTLS() {
		log.Printf("Running with TLS on port %d ...", port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), tlsCert, tlsKey, r))
	} else {
		log.Printf("Running on port %d ...", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}
}

func withTLS() bool {
	return tlsKey != "" && tlsCert != ""
}
