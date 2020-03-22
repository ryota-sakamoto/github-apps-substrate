package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v30/github"

	"github.com/ryota-sakamoto/github-apps-substrate/service"
)

type CallbackController struct {
	installationService service.InstallationService
	subscribeService    service.SubscribeService
}

func NewCallbackController(is service.InstallationService, ss service.SubscribeService) CallbackController {
	return CallbackController{
		installationService: is,
		subscribeService:    ss,
	}
}

func (ca CallbackController) Endpoint(c *gin.RouterGroup) {
	c.POST("/callback", ca.callback)
}

func (ca CallbackController) callback(c *gin.Context) {
	var event interface{}
	var hook func() error
	switch c.GetHeader("X-Github-Event") {
	case "installation", "installation_repositories":
		event = github.InstallationEvent{}
		hook = func() error {
			e := event.(github.InstallationEvent)
			return ca.installationService.Action(&e)
		}
	case "push":
		event = github.PushEvent{}
		hook = func() error {
			e := event.(github.PushEvent)
			return ca.subscribeService.SubscribePush(&e)
		}
	default:
		c.AbortWithStatus(400)
		return
	}

	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(400)
		return
	}

	if err := hook(); err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(500)
		return
	}
}
