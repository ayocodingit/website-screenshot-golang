package main

import (
	"net/http"

	"github.com/ayocodingit/website-screenshot-golang/src/domain"
	"github.com/ayocodingit/website-screenshot-golang/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func main() {
	go utils.InitCron()

	r := gin.New()
	r.Static("/public", "./public")
	r.Use(location.Default())
	r.Use(cors.Default())
	r.POST("/", handler)

	r.Run()
}

func handler(c *gin.Context) {
	var screenshot domain.Screenshot
	if err := c.ShouldBind(&screenshot); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": err.Error()})
		return
	}

	filename, err := utils.GetScreenshot(screenshot.URL, screenshot.Quality)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	res := location.Get(c).String() + "/" + filename
	c.JSON(http.StatusOK, res)
}
