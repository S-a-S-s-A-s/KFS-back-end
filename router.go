package main

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	apiRouterIndex := r.Group("/index")
	apiRouterFile := r.Group("/file")
	// index apis
	apiRouterIndex.GET("/info", Info)
	apiRouterIndex.POST("/login", Login)
	apiRouterIndex.POST("/logout", Logout)
	// file apis
	apiRouterFile.GET("/:page/:limit", PageAndLimit)
	apiRouterFile.GET("/projects", GetProject)
	apiRouterFile.POST("/download/:Name/:Project", DownloadFile)
	apiRouterFile.POST("/upload", UploadFile)
	apiRouterFile.GET("/info/:Name/:Project", GetFileSize)
	apiRouterFile.POST("/download", DownloadBigFile)
	apiRouterFile.POST("/download/auth", DownloadAuth)
	apiRouterFile.GET("/uploadAuth", UploadAuth)
}
