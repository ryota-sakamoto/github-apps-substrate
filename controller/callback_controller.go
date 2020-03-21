package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v30/github"

	"github.com/ryota-sakamoto/github-apps-substrate/service"
)

type CallbackController struct {
	installationService service.InstallationService
}

func NewCallbackController(is service.InstallationService) CallbackController {
	return CallbackController{
		installationService: is,
	}
}

func (ca CallbackController) Endpoint(c *gin.RouterGroup) {
	c.POST("/callback", ca.callback)
}

func (ca CallbackController) callback(c *gin.Context) {
	switch c.GetHeader("X-Github-Event") {
	case "installation", "installation_repositories":
		ca.installation(c)
	default:
		c.AbortWithStatus(400)
		return
	}
}

func (ca CallbackController) installation(c *gin.Context) {
	var event github.InstallationEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(400)
		return
	}

	if err := ca.installationService.Action(&event); err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(500)
		return
	}
}
