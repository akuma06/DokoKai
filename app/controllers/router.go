package controllers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var router *gin.Engine
var once sync.Once

// Get return a router signleton
func Get() *gin.Engine {
	once.Do(func() {
		if viper.Get("env") == "PRODUCTION" {
			gin.SetMode(gin.ReleaseMode)
		}
		router = gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	})
	return router
}