package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fs714/github-star-manager/api"
	"github.com/fs714/github-star-manager/db/jsondb"
	"github.com/fs714/github-star-manager/pkg/config"
	"github.com/fs714/github-star-manager/pkg/utils/log"
	"github.com/fs714/github-star-manager/pkg/utils/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	httpHost     string
	httpPort     string
	readTimeout  int
	writeTimeout int
	dbPath       string
	logFile      string
	logLevel     string
	logFormat    string
)

var StartCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start http server",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		initLog()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func InitStartCmd() {
	StartCmd.PersistentFlags().SortFlags = false
	StartCmd.Flags().SortFlags = false

	StartCmd.Flags().StringVarP(&httpHost, "host", "l", config.DefaultConfig.HttpServer.Host,
		"Http server listening address")
	config.Viper.BindPFlag("http_server.host", StartCmd.Flags().Lookup("host"))
	config.Viper.BindEnv("http_server.host", "HTTP_HOST")

	StartCmd.Flags().StringVarP(&httpPort, "port", "p", config.DefaultConfig.HttpServer.Port,
		"Http server listening port")
	config.Viper.BindPFlag("http_server.port", StartCmd.Flags().Lookup("port"))
	config.Viper.BindEnv("http_server.port", "HTTP_PORT")

	StartCmd.Flags().IntVarP(&readTimeout, "read-timeout", "", config.DefaultConfig.HttpServer.ReadTimeout,
		"Http server read timeout")
	config.Viper.BindPFlag("http_server.read_timeout", StartCmd.Flags().Lookup("read-timeout"))
	config.Viper.BindEnv("http_server.read_timeout", "HTTP_READ_TIMEOUT")

	StartCmd.Flags().IntVarP(&writeTimeout, "write-timeout", "", config.DefaultConfig.HttpServer.WriteTimeout,
		"Http server write timeout")
	config.Viper.BindPFlag("http_server.write_timeout", StartCmd.Flags().Lookup("write-timeout"))
	config.Viper.BindEnv("http_server.write_timeout", "HTTP_WRITE_TIMEOUT")

	StartCmd.Flags().StringVarP(&dbPath, "db-path", "", config.DefaultConfig.Database.Path, "Path for database file")
	config.Viper.BindPFlag("database.path", StartCmd.Flags().Lookup("db-path"))
	config.Viper.BindEnv("database.path", "DB_PATH")

	StartCmd.Flags().StringVarP(&logFile, "log-file", "", config.DefaultConfig.Logging.File,
		"Set logging file, stderr will be used if file is empty string")
	config.Viper.BindPFlag("logging.file", StartCmd.Flags().Lookup("log-file"))
	config.Viper.BindEnv("logging.file", "LOGGING_FILE")

	StartCmd.Flags().StringVarP(&logLevel, "log-level", "", config.DefaultConfig.Logging.Level,
		"Set logging level, could be info or debug")
	config.Viper.BindPFlag("logging.level", StartCmd.Flags().Lookup("log-level"))
	config.Viper.BindEnv("logging.level", "LOGGING_LEVEL")

	StartCmd.Flags().StringVarP(&logFormat, "log-format", "", config.DefaultConfig.Logging.Format,
		"Set logging format, could be console or json")
	config.Viper.BindPFlag("logging.format", StartCmd.Flags().Lookup("log-format"))
	config.Viper.BindEnv("logging.format", "LOGGING_FORMAT")
}

func initLog() {
	if config.Config.Logging.File == "" {
		logger := log.New(os.Stderr, log.ParseFormat(config.Config.Logging.Format),
			log.ParseLevel(config.Config.Logging.Level), true)
		log.ResetCurrentLog(logger)
	} else {
		var tops = []log.TeeWithRotateOption{
			{
				Filename:   config.Config.Logging.File,
				MaxSize:    config.Config.Logging.MaxSize,
				MaxAge:     config.Config.Logging.MaxAge,
				MaxBackups: config.Config.Logging.MaxBackups,
				Compress:   config.Config.Logging.Compress,
				Lef: func(lvl log.Level) bool {
					return lvl >= log.ParseLevel(config.Config.Logging.Level)
				},
				F: log.ParseFormat(config.Config.Logging.Format),
			},
		}

		logger := log.NewTeeWithRotate(tops, true)
		log.ResetCurrentLog(logger)
	}
}

func startServer() (err error) {
	log.Infow("start http server", "BaseVersion", version.BaseVersion, "GitVersion", version.GitVersion)

	log.Infow("initialize database", "path", config.Config.Database.Path)
	err = jsondb.InitJsondbFromConfig()
	if err != nil {
		log.Errorf("failed to initialize database:\n%+v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exitWg := &sync.WaitGroup{}

	{
		router := api.InitRouter()
		srv := &http.Server{
			Addr:           fmt.Sprintf("%s:%s", config.Config.HttpServer.Host, config.Config.HttpServer.Port),
			Handler:        router,
			ReadTimeout:    time.Duration(time.Duration(config.Config.HttpServer.ReadTimeout) * time.Second),
			WriteTimeout:   time.Duration(time.Duration(config.Config.HttpServer.WriteTimeout) * time.Second),
			MaxHeaderBytes: 1 << 20,
		}

		exitWg.Add(1)
		go func() {
			defer exitWg.Done()

			log.Infow("start http server", "Host", config.Config.HttpServer.Host, "Port", config.Config.HttpServer.Port)
			err = srv.ListenAndServe()
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					err = nil
				} else {
					log.Errorf("failed to start http server:\n%+v", err)
					return
				}
			}
		}()

		exitWg.Add(1)
		go func(ctx context.Context) {
			defer exitWg.Done()

			<-ctx.Done()
			cctx, ccancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer ccancel()
			err := srv.Shutdown(cctx)
			if err != nil {
				log.Errorf("failed to close http server:\n%+v", err)
			}
			log.Infow("http server exit")
		}(ctx)
	}

	<-signalCh
	cancel()
	exitWg.Wait()
	log.Infow("github-star-manager exit")

	return
}
