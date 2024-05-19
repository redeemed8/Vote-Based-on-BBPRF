package service

import (
	"github.com/gin-gonic/gin"
	"mini/models"
	"mini/models/dto"
	"mini/models/vo"
	"mini/util"
	"mini/util/jwt"
	"net/http"
	"strconv"
)

func NewVote(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 绑定参数
	var newVo vo.NewVoteVo
	if err := ctx.ShouldBind(&newVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	3. 参数校验
	if right, err := newVo.Right(); !right {
		ctx.JSON(http.StatusOK, resp.Fail(400, err.Error()))
		return
	}
	//	4. 转换成正式的vote
	vote := newVo.Parse(userClaim.Id)
	//	5. 添加到数据库
	err := models.VoteDao.CreateVote(&vote)
	if util.MysqlErr(err) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(vote.Title))
}

// SearchVote
// url?key=xxx
func SearchVote(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	keyword := ctx.Query("key")
	//	3. 查询此人所有符合条件的投票
	votes, err := models.VoteDao.SearchVote(userClaim.Id, keyword)
	if util.MysqlErr(err) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	//	4. 返回dtos
	ctx.JSON(http.StatusOK, resp.Success(votes.ToDtos(true)))
}

// GetVoteDetail
// url?id=xxx
func GetVoteDetail(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	vid, err := strconv.Atoi(ctx.Query("id"))
	if err != nil || vid < 1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	//	3. 查询出该投票
	queryVote, queryE := models.VoteDao.GetVote(models.Vote{Id: vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryVote.Title == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	dto_ := dto.GetVoteDto{
		Id:            queryVote.Id,
		Title:         queryVote.Title,
		IsMultiChoice: queryVote.IsMultiChoice,
		Status:        queryVote.Status,
		Options:       models.VoteUtil.Options(queryVote.Options),
	}
	//	4. 是否是当前投票的发起人
	if userClaim.Id == queryVote.PublisherId {
		dto_.Participants = queryVote.Participants
		dto_.AnsCount = models.VoteUtil.ParseAnsCount(queryVote.AnsCount)
	}
	ctx.JSON(http.StatusOK, resp.Success(dto_))
}

// DelVote
// url?id=xxx
func DelVote(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	vidStr := ctx.Query("id")
	vid, tErr := strconv.Atoi(vidStr)
	if vidStr == "" || tErr != nil || vid < 1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	//	3. 查询投票信息
	queryVote, queryE := models.VoteDao.GetVote(models.Vote{Id: vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryVote.PublisherId != userClaim.Id {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权删除此投票"))
		return
	}
	//	4. 进行删除
	delErr := models.VoteDao.DelVoteById(vid)
	if util.MysqlErr(delErr) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success("投票已删除"))
}

// UpdateVoteStatus
// url?vid=xxx&imc=xxx
func UpdateVoteStatus(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 获取路径参数
	vidStr := ctx.Query("vid")
	if vidStr == "" || len(vidStr) > 12 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	vid, err := strconv.Atoi(vidStr)
	if err != nil || vid < 1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	imc := ctx.Query("imc")
	if imc != "1" && imc != "0" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效的多选值"))
		return
	}
	//	3. 获取到该投票的信息
	queryV, queryE := models.VoteDao.GetVote(models.Vote{Id: vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryV.PublisherId != userClaim.Id {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权修改发布状态"))
		return
	}
	//	4. 修改发布状态
	newStatus := models.VoteUtil.UpdateStatus(queryV.Status)
	newImc, _ := strconv.Atoi(imc)
	uptE := models.VoteDao.UpdateVote(queryV.Id, map[string]interface{}{"status": newStatus, "is_multi_choice": newImc})
	if util.MysqlErr(uptE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(gin.H{
		"status": newStatus,
		"imc":    newImc,
	}))
}

func DoVote(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	_, login := jwt.ParseToken(ctx, resp)
	if !login {
		return
	}
	//	2. 绑定参数
	var doVoteVo vo.DoVoteVo
	if err := ctx.ShouldBind(&doVoteVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	3. 参数校验完成，随时可以合并
	if right, err := doVoteVo.Right(); !right {
		ctx.JSON(http.StatusOK, resp.Fail(400, err.Error()))
		return
	}
	//	4. 检验该投票是否存在
	queryV, queryE := models.VoteDao.GetVote(models.Vote{Id: doVoteVo.Vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryV.Title == "" || queryV.Options == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	if queryV.Status == 0 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票尚未发布"))
		return
	}
	//	5. 校验是否有投票的资格
	if ok, err := IsLegalToken(); !ok {
		ctx.JSON(http.StatusOK, resp.Fail(400, err.Error()))
		return
	}
	//	6. 资格通过，进行票数的合并
	newAnsCount := models.VoteUtil.MergeOptionsToCountStr(queryV.AnsCount, doVoteVo.SelectOptions, queryV.IsMultiChoice)
	UptE := models.VoteDao.UpdateVote(doVoteVo.Vid, map[string]interface{}{"ans_count": newAnsCount})
	if util.MysqlErr(UptE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success("ok"))
}

func IsLegalToken() (bool, error) {
	return true, nil
}
