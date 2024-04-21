package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	router := httprouter.New()
	//ytApiKey : os.Getenv("YOUTUBE_API_KEY")
	//channelId := os.Getenv("YOUTUBE_CHANNEL_ID")

	ytApiKey := "AIzaSyBujADWwDBFO2hAd16A6GKwbGFh_WmiDQ0"
	channelId := "UCAsarZPd1ULXKqOHnGLdmXw"

	router.GET("/youtube/channel/stats", getChannelStats(ytApiKey, channelId))
	return router
}

func main() {
	svr := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		signit := make(chan os.Signal, 1)
		signal.Notify(signit, os.Interrupt)
		signal.Notify(signit, syscall.SIGTERM)
		<-signit

		log.Println("Service interrupted. Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server failed to shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Printf("Service started on %s", svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Service stopped")

}
