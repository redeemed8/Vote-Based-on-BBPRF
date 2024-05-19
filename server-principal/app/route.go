package app

import (
	"github.com/gin-gonic/gin"
	"mini_pkcs/model"
	"mini_pkcs/service"
)

func InitModel() {

	model.InitInfo()

}

func AppsRouter() *gin.Engine {
	r := gin.Default()

	InitModel()

	pkcs := r.Group("/pkcs")
	{
		pkcs.POST("/get/sm2pk", service.GetSM2PublicKey)
		pkcs.POST("/get/signed/privk", service.GetSignedPrivKey)
	}

	return r
}
