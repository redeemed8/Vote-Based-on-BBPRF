package service

import (
	"github.com/gin-gonic/gin"
	"mini/models"
	"mini/models/dto"
	"mini/models/vo"
	"mini/util"
	"mini/util/jwt"
	"net/http"
	"strings"
)

func NewQuestionnaire(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 绑定参数
	var newVo vo.NewQuestionnaireVo
	if err := ctx.ShouldBind(&newVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	3. 参数校验
	if ok, err := newVo.Right(); !ok {
		ctx.JSON(http.StatusOK, resp.Fail(400, err.Error()))
		return
	}
	//	4. 参数转换
	questionnaire, questions := newVo.Parse(userClaim.Id)
	//	5. 保存到数据库
	err1 := models.QuestionnaireDao.CreateQuestionnaire(questionnaire)
	if util.MysqlErr(err1) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	err2 := models.QuestionDao.CreateQuestions(questions)
	if util.MysqlErr(err2) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(newVo.Title))
}

// SearchQuestionnaire
// url?key=xxx
func SearchQuestionnaire(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	keyword := ctx.Query("key")
	//	3. 查询此人所有符合条件的问卷
	questionnaires, err := models.QuestionnaireDao.SearchQuestionnaire(userClaim.Id, keyword)
	if util.MysqlErr(err) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	//	4. 返回dtos
	ctx.JSON(http.StatusOK, resp.Success(questionnaires.ToDtos(true)))
}

type QuestionnaireDetailDto struct {
	Id        string            `json:"id"`
	Title     string            `json:"title"`
	AnsNumber int               `json:"ans_number"`
	Questions []dto.QuestionDto `json:"questions"`
	Status    int               `json:"status"`
}

// GetQuestionnaireDetail
// url?qid=xxx
func GetQuestionnaireDetail(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	qid := ctx.Query("qid")
	if qid == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "问卷id不能为空"))
		return
	}
	if len(qid) > 17 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "问卷不存在或已被删除"))
		return
	}
	//	3. 查询该问卷信息
	var detailDto QuestionnaireDetailDto
	detailDto.Questions = make([]dto.QuestionDto, 0)
	questionnaire, err1 := models.QuestionnaireDao.GetQuestionnaire(qid)
	if util.MysqlErr(err1) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if questionnaire.Id == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "未找到相关问卷"))
		return
	}
	detailDto.Id = questionnaire.Id
	detailDto.Title = questionnaire.Title
	if questionnaire.PublisherId == userClaim.Id {
		detailDto.AnsNumber = questionnaire.AnsNumber
	} else {
		detailDto.AnsNumber = 0
	}
	detailDto.Status = questionnaire.Status

	questionIds := strings.Split(questionnaire.QuestionIds, models.QuestionnaireUtil.GetIdsSplit())

	//	4. 查询所有的问题
	for _, id := range questionIds {
		question, err2 := models.QuestionDao.GetQuestion(id, questionnaire.Id)
		if util.MysqlErr(err2) {
			ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
			return
		}
		if question.Type != "" {
			var questionDto dto.QuestionDto
			questionDto.Id = question.Id
			questionDto.FatherId = question.FatherId
			questionDto.Title = question.Title
			questionDto.Type = question.Type
			questionDto.Options = []string{models.VoidOption}
			if question.Type == models.Select {
				questionDto.Options = strings.Split(question.Options, models.QuestionUtil.GetOptionSplit())
			}
			detailDto.Questions = append(detailDto.Questions, questionDto)
		}
	}

	ctx.JSON(http.StatusOK, resp.Success(detailDto))
}

// DeleteQuestionnaire
// url?id=xxx
func DeleteQuestionnaire(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	qid := ctx.Query("id")
	if qid == "" || len(qid) > 17 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "问卷不存在或已被删除"))
		return
	}
	//	3. 查询该问卷
	queryQ, queryE := models.QuestionnaireDao.GetQuestionnaire(qid)
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryQ.PublisherId != userClaim.Id {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权利删除此问卷"))
		return
	}
	//	4. 进行删除
	err := models.QuestionnaireDao.DelQuestionnaire(qid, userClaim.Id)
	if util.MysqlErr(err) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success("问卷已删除"))
}

// UpdateQuestionnaireStatus
// url?qid=xxx
func UpdateQuestionnaireStatus(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	qid := ctx.Query("qid")
	if qid == "" || len(qid) > 17 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "问卷不存在或已被删除"))
		return
	}
	//	3. 获取到该问卷的信息
	queryQ, queryE := models.QuestionnaireDao.GetQuestionnaire(qid)
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryQ.PublisherId != userClaim.Id {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权修改发布状态"))
		return
	}
	//	4. 修改发布状态
	newStatus := models.QuestionnaireUtil.UpdateStatus(queryQ.Status)
	uptE := models.QuestionnaireDao.UptQuestionnaire(queryQ.Id, queryQ.PublisherId, map[string]interface{}{"status": newStatus})
	if util.MysqlErr(uptE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(newStatus))
}
