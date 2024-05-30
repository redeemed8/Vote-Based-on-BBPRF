package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mini/config"
	"mini/models"
	"mini/models/dto"
	"mini/models/vo"
	"mini/util"
	"mini/util/fileutil"
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
	//	2. 获取表单参数
	var newVo2 vo.NewVoteVo2
	newVo2.VoteTitle = ctx.PostForm("vote_title")
	newVo2.OptionNames = ctx.PostForm("options")
	if ctx.PostForm("is_multi_choice") != "0" {
		newVo2.IsMultiChoice = 1
	} else {
		newVo2.IsMultiChoice = 0
	}
	//	3. 转换为 string数组
	var options []string
	if err := json.Unmarshal([]byte(newVo2.OptionNames), &options); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	var newVo vo.NewVoteVo
	newVo.VoteTitle = newVo2.VoteTitle
	newVo.IsMultiChoice = newVo2.IsMultiChoice
	newVo.Options = make([]string, 0)
	for _, option := range options {
		newVo.Options = append(newVo.Options, option)
	}
	//	4. 参数校验
	if right, err := newVo.Right(); !right {
		ctx.JSON(http.StatusOK, resp.Fail(400, err.Error()))
		return
	}
	//	5. 转换成正式的vote
	vote := newVo.Parse(userClaim.Id)
	//	6. 获取图片
	file, fileHeader, fileErr := ctx.Request.FormFile("file")
	if fileErr != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "请放一张图片"))
		return
	}
	//	7. 解析图片
	ossFile, parseErr := models.OssFileUtil.FileToOssFile(file, fileHeader)
	if parseErr != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, parseErr.Error()))
		return
	}
	//	8. 图片解析成功, 将图片添加到数据库
	createErr := models.OssFileDao.CreateOssFile(ossFile)
	if util.MysqlErr(createErr) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	photoUrlChan <- *ossFile
	vote.PhotoUrl = ossFile.FileUrl
	//	9. 投票信息添加到数据库
	err := models.VoteDao.CreateVote(&vote)
	if util.MysqlErr(err) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	fmt.Println("url == ", ossFile.FileUrl)
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
		Url:           queryVote.PhotoUrl,
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

var photoUrlChan = make(chan models.OssFile)

func UploadPhoto() {
	fmt.Println("任务 --- 上传云图片  已开启")
	for {
		select {
		case ossFile, ok := <-photoUrlChan:
			if ok {
				//	上传云图片
				err1 := UploadFileToOss(ossFile)
				if err1 != nil {
					fmt.Println("上传云图片出错 , cause by : ", err1.Error())
					return
				}
				//	删除临时文件
				err2 := fileutil.DeleteFile(ossFile.LocalPath)
				if err2 != nil {
					fmt.Println("删除本地临时文件出错 , cause by : ", err2.Error())
					return
				}
			}
		}
	}
}

func UploadFileToOss(ossFile models.OssFile) error {
	err := config.IVBucket.PutObjectFromFile(ossFile.FileObjectName, ossFile.LocalPath)
	if err != nil {
		return errors.New("oss上传文件错误:" + err.Error())
	}
	return nil
}
