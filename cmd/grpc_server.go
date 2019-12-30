package cmd

import (
	"fmt"
	"net/http"

	"github.com/orensimple/otus_events_api/config"
	"github.com/orensimple/otus_events_api/internal/domain/services"
	"github.com/orensimple/otus_events_api/internal/grpc/api"
	"github.com/orensimple/otus_events_api/internal/logger"
	"github.com/orensimple/otus_events_api/internal/maindb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addrGRPC, addr string

// TODO: dependency injection, orchestrator
func construct() (*api.CalendarServer, error) {
	err := config.Init(addr)
	if err != nil {
		logger.ContextLogger.Errorf("Eror init config, viper.ReadInConfig", err.Error())
	}
	logger.InitLogger()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/events", viper.GetString("postgres.user"), viper.GetString("postgres.passwd"), viper.GetString("postgres.ip"), viper.GetString("postgres.port"))
	eventStorage, err := maindb.NewPgEventStorage(dsn)
	if err != nil {
		return nil, err
	}
	eventService := &services.EventService{
		EventStorage: eventStorage,
	}
	server := &api.CalendarServer{
		EventService: eventService,
	}
	return server, nil
}

var RootCmd = &cobra.Command{
	Use:   "grpc_server",
	Short: "Run grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := construct()
		if err != nil {
			logger.ContextLogger.Errorf(" Cannot init logger, config, storage", err)
		}
		logger.ContextLogger.Infof(" [*] GRPC server run. To exit press CTRL+C")
		go func() {
			err = server.Serve(addrGRPC)
			if err != nil {
				logger.ContextLogger.Errorf(" Cannot start GRPC server", err)
			}
		}()

		http.Handle("/metrics", promhttp.Handler())
		logger.ContextLogger.Infof("Starting web server at %s\n", "events-api:9110")
		err = http.ListenAndServe("events-api:9120", nil)
		if err != nil {
			logger.ContextLogger.Errorf("http.ListenAndServer for metrics: %v\n", err.Error())
		}
	},
}

func init() {
	RootCmd.Flags().StringVar(&addrGRPC, "addr", "events-api:8088", "host:port to listen")
	RootCmd.Flags().StringVar(&addr, "config", "./config", "")
}
