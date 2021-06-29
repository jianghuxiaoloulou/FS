package routers

import (
	"WowjoyProject/FileServer/pkg/setting"
	v1 "WowjoyProject/FileServer/routers/api/v1"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RUN_MODE)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		// 上传单文件
		apiv1.POST("/file", v1.UploadFile)
		// 检查号上传
		apiv1.POST("/numbers", v1.UploadNumbers)
		// 单文件下载
		apiv1.GET("/file/:id", v1.DownFile)
		// 检查号下载
		apiv1.GET("/numbers/:id", v1.DownNumbers)
		// 单文件删除
		apiv1.DELETE("/file", v1.DeleteFile)
		// 检查号删除
		apiv1.DELETE("/numbers", v1.DeleteNumbers)
	}
	return r
}
