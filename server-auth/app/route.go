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
		vote.POST("/do-vote", service.DoVote)
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

	return r
}
