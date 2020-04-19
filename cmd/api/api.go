package main

import (
	"log"
	"os"

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
		os.Exit(1)
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

	r.Run(":8080")
}
