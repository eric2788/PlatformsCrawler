package rest

import (
	"fmt"
	"github.com/eric2788/PlatformsCrawler/logging"
	"github.com/gin-gonic/gin"
)

type RouterHook func(group *gin.RouterGroup)

var (
	logger       = logging.GetMainLogger()
	routerGroups = make(map[string]RouterHook, 0)
)

func StartServe(port int) {

	r := gin.Default()

	for path, hook := range routerGroups {
		hook(r.Group(path))
	}

	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal(err)
	}

}

func AddHook(path string, hook RouterHook) {
	routerGroups[path] = hook
}
