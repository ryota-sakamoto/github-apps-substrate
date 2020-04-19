package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ryota-sakamoto/github-apps-substrate/internal/config"
	"github.com/ryota-sakamoto/github-apps-substrate/internal/handler"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/middleware"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/repository"
	"github.com/ryota-sakamoto/github-apps-substrate/pkg/service"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	r := gin.Default()

	rr := repository.NewRepositoryRepository(conf.GitHub.PrivateKey, conf.GitHub.AppID)
	ss := service.NewSubscribeService(rr)
	callback := handler.NewCallbackHandler(ss)

	api := r.Group("/api")
	api.Use(middleware.ValidatePayload([]byte(conf.GitHub.Secret)))
	{
		callback.Endpoint(api)
	}

	port := 8080
	if conf.Server.Port != 0 {
		port = conf.Server.Port
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server")

	ctx, cannel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cannel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%+v\n", err)
	}

	log.Println("shutdown success")
}
