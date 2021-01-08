package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/logger"
)

// Start starts the http server
func Start() {
	serverCfg := config.GetServer()
	addr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	httpServer := http.Server{
		Addr:    addr,
		Handler: GetRouter(),
	}

	go func() {
		logger.Infoln("Http server has been started on", addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Errorln("Failed to start http server,", err)
		}
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorln("Http server couldn't shutdown gracefully", err)
	}
	logger.Infoln("Http server has been shutdown gracefully")
}
