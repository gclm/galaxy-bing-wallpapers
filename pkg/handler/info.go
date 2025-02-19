package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoResponse struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
}

func GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, InfoResponse{
		Name:        "galaxy-bing-wallpapers",
		Version:     "1.0.0",
		Author:      "gclm",
		Description: "Bing wallpaper API service",
		Repository:  "https://github.com/gclm/galaxy-bing-wallpapers",
	})
}
