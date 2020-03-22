package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/ryota-sakamoto/github-apps-substrate/config"
	"github.com/ryota-sakamoto/github-apps-substrate/controller"
	"github.com/ryota-sakamoto/github-apps-substrate/middleware"
	"github.com/ryota-sakamoto/github-apps-substrate/repository"
	"github.com/ryota-sakamoto/github-apps-substrate/service"
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
	callback := controller.NewCallbackController(ss)

	api := r.Group("/api")
	api.Use(middleware.ValidatePayload([]byte(conf.GitHub.Secret)))
	{
		callback.Endpoint(api)
	}

	r.Run(":8080")
}
