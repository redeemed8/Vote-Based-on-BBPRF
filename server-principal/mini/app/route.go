package app

import (
	"github.com/gin-gonic/gin"
	"mini/models"
	"mini/service"
)

func InitModel() {
	models.InitUserInfo()
	models.InitQuestion()
	models.InitQuestionnaire()
	models.InitVote()
	models.InitH()
	models.InitOssFile()
}

func AppsRouter() *gin.Engine {
	r := gin.Default()

	InitModel()

	r.GET("/goi", service.GetOpenId)

	user := r.Group("/user")
	{
		user.GET("/send/code", service.SendVerifyCode)
		user.POST("/login/phone", service.PhoneLogin)
		user.POST("/register/account", service.RegisterAccount)
		user.POST("/login/account", service.AccountLogin)
		user.GET("/get/name/avatar", service.GetNameAndUrl)
		user.POST("/update/name/avatar", service.UpdateNameAndUrl)
		user.GET("/get/publish/num", service.GetPublishNum)
	}

	question := r.Group("/question")
	{
		question.POST("/new", service.NewQuestionnaire)
		question.GET("/search", service.SearchQuestionnaire)
		question.GET("/get-detail", service.GetQuestionnaireDetail)
		question.GET("/del", service.DeleteQuestionnaire)
		question.GET("/update/status", service.UpdateQuestionnaireStatus)
	}

	vote := r.Group("/vote")
	{
		vote.POST("/new", service.NewVote)
		vote.GET("/search", service.SearchVote)
		vote.GET("/get-detail", service.GetVoteDetail)
		vote.GET("/del", service.DelVote)
		vote.GET("/update/status", service.UpdateVoteStatus)
	}

	auth := r.Group("/auth")
	{
		auth.GET("/verify/token", service.VerifyJWTToken)
		auth.GET("/get-pprms", service.GetPprmS)
		auth.POST("/act/sign", service.Sign)
		auth.POST("/act/verify", service.Verify)
	}

	go service.CreateH()
	go service.AddH()
	go service.DealWithV()
	go service.UploadPhoto()

	return r
}
