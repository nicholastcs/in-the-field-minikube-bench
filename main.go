package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	_ "go.uber.org/automaxprocs"
	"k8s.io/client-go/kubernetes"
)

var logger hclog.Logger
var clientset *kubernetes.Clientset
var startup = false
var health = false

func init() {
	podName := os.Getenv("POD_NAME")
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       "main",
		JSONFormat: false,
		Level:      hclog.Debug,
	}).Named(podName)

	cs, err := login()
	if err != nil {
		logger.Error("unable to create clientset")
	}
	clientset = cs
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	jobMode := flag.Bool("job", false, "is job")
	flag.Parse()

	// starts the http server
	srv := getGinServer()

	if *jobMode {
		doWork()
		os.Exit(0)
	} else {
		confPath := os.Getenv("CONF_PATH")
		conf, err := getConf(confPath)
		if err != nil {
			logger.Error("unable to read configuration", "file", confPath)
		}

		logger.Info("app starts with pi calculator",
			"workers", conf.Workers,
			"iterations", conf.Iterations,
		)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("unable to serve", "port", "8080")

				os.Exit(1)
			}
		}()

		time.Sleep(750)
		startup = true
		health = true

		logger.Info("serving traffic...")
	}

	// when receive shutdown signals
	s := <-sig
	logger.Warn("shutdown received!", "signal", s)

	// stop server
	if !*jobMode {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("shutdown failed", err)
		}
	}

	logger.Warn("shutdown done!", "signal", s)
}
