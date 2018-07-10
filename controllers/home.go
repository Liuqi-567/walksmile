package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {

	pic := BackgroundPic
	if c.GetBool("isMobile") {
		pic = strings.Replace(pic, "1920x1080", "480x800", -1)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"hotwords": Hotwords,
		"pic":      pic,
	})
}

func Wallpaper(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	if len(BackgroundPicArray) <= 0 {
		pics, err := WallPaperList()
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		BackgroundPicArray = pics
	}

	index := id % len(BackgroundPicArray)
	pic := BackgroundPicArray[index]
	if c.GetBool("isMobile") {
		pic = strings.Replace(pic, "1920x1080", "480x800", -1)
	}

	c.HTML(http.StatusOK, "wallpaper.html", gin.H{
		"pic":   pic,
		"id":    index + 1,
		"title": "bing壁纸-行笑网",
	})

}
