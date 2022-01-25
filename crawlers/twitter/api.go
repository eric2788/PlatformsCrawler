package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/eric2788/PlatformsCrawler/rest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func restApi(group *gin.RouterGroup) {
	group.GET("/userExist/:screen", checkUserExist)
}

func checkUserExist(ctx *gin.Context) {
	screen := ctx.Param("screen")
	m, err := UserLookUpCache([]string{screen})

	if err != nil {
		if twErr, ok := err.(twitter.APIError); ok {
			ctx.IndentedJSON(http.StatusBadRequest, twErr.Errors)
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if id, ok := m[screen]; ok {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"exist": true,
			"data": map[string]string{
				"screen_name": screen,
				"id":          id,
			},
		})
	} else {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"exist": false,
			"data":  map[string]string{},
		})
	}
}

func init() {
	rest.AddHook("/twitter", restApi)
}
