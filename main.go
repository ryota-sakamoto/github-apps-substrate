package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/ryota-sakamoto/github-apps-substrate/controller"
	"github.com/ryota-sakamoto/github-apps-substrate/middleware"
	"github.com/ryota-sakamoto/github-apps-substrate/service"
)

func main() {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	r := gin.Default()
	r.Use(middleware.ValidatePayload([]byte(secret)))

	is := service.NewInstallationService()
	callback := controller.NewCallbackController(is)

	api := r.Group("/api")
	{
		callback.Endpoint(api)
	}

	r.Run("localhost:8080")
}
