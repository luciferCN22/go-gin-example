package main

import (
	"context"
	"fmt"
	"github.com/luciferCN22/go-gin-example/models"
	"github.com/luciferCN22/go-gin-example/pkg/gredis"
	"github.com/luciferCN22/go-gin-example/pkg/logging"
	"github.com/luciferCN22/go-gin-example/pkg/setting"
	"github.com/luciferCN22/go-gin-example/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen:%s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
