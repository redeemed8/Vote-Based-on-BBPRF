package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mini/models"
	"mini/models/vo"
	"mini/util"
	"mini/util/jwt"
	"net/http"
)

// SendVerifyCode 发送验证码
// api: uri?phone=xxx&mode=xxx
func SendVerifyCode(ctx *gin.Context) {
	resp := models.NewResp()
	ctx.JSON(http.StatusOK, resp.Success("验证码已发送"))
}

const ServerErr = "服务器异常，请稍后再试"

func PhoneLogin(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 绑定参数
	var loginVo vo.PhoneLoginVo
	if err := ctx.ShouldBind(&loginVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	2. 参数校验
	if loginVo.Code == "" || len(loginVo.Code) != 6 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "验证码错误"))
		return
	}
	if ok := util.MatchPhone(loginVo.Phone); !ok {
		ctx.JSON(http.StatusOK, resp.Fail(400, "手机号格式错误"))
		return
	}
	//	3. 检查验证码，这里我们默认为手机号后6位
	rightCode := loginVo.Phone[len(loginVo.Phone)-6:]
	if rightCode != loginVo.Code {
		ctx.JSON(http.StatusOK, resp.Fail(400, "验证码错误"))
		return
	}
	//	4. 验证码正确，查看此手机号是否登录过
	var userUUID string
	queryUser, err1 := models.UserInfoDao.GetUserByPhone(loginVo.Phone)
	if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	if queryUser.Phone == "" {
		//	此用户是第一次登录
		userUUID = uuid.New().String()
		newUser := models.UserInfo{Phone: loginVo.Phone, UUID: userUUID}
		if err := models.UserInfoDao.CreateUser(&newUser); err != nil {
			ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
			return
		}
		queryUser = newUser
	}
	//	5. 创建 token 并返回
	token, err2 := jwt.MakeJWTToken(queryUser)
	if err2 != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(token))
}

func RegisterAccount(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 绑定参数
	var registerVo vo.RegisterAccountVo
	if err := ctx.ShouldBind(&registerVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	2. 参数校验
	if registerVo.Account == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "账号不能为空"))
		return
	}
	if registerVo.Password == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "密码不能为空"))
		return
	}
	if registerVo.Phone == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "手机号不能为空"))
		return
	}
	//	3. 检查参数准确性
	if ok := util.MatchAccount(registerVo.Account); !ok {
		ctx.JSON(http.StatusOK, resp.Fail(400, "账号格式不正确"))
		return
	}
	if registerVo.Password != registerVo.Repassword {
		ctx.JSON(http.StatusOK, resp.Fail(400, "两次密码不一致"))
		return
	}
	rightCode := "123456"
	if registerVo.Code != rightCode {
		ctx.JSON(http.StatusOK, resp.Fail(400, "验证码错误"))
		return
	}
	//	4. 参数正确，检查此账号和手机号是否已被使用
	queryUser, err1 := models.UserInfoDao.GetUserByAccount(registerVo.Account)
	if util.MysqlErr(err1) {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	if queryUser.Account != "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "该账号已被注册"))
		return
	}
	var err2 error
	queryUser, err2 = models.UserInfoDao.GetUserByPhone(registerVo.Phone)
	if util.MysqlErr(err2) {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	if queryUser.Phone != "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "该手机号已被注册"))
		return
	}
	//	5. 为其创建账号
	newUser := models.UserInfo{
		Account:  registerVo.Account,
		Password: registerVo.Password,
		Phone:    registerVo.Phone,
		UUID:     uuid.New().String(),
	}
	if err := models.UserInfoDao.CreateUser(&newUser); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success("注册成功"))
}

func AccountLogin(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 绑定参数
	var loginVo vo.AccountLoginVo
	if err := ctx.ShouldBind(&loginVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	2. 参数校验
	if loginVo.Account == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "账号不能为空"))
		return
	}
	if loginVo.Password == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "密码错误"))
		return
	}
	//	3. 查询密码
	queryUser, err1 := models.UserInfoDao.GetUserByAccount(loginVo.Account)
	if util.MysqlErr(err1) {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	if queryUser.Account == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "该账号不存在"))
		return
	}
	if queryUser.Password != loginVo.Password {
		ctx.JSON(http.StatusOK, resp.Fail(400, "密码错误"))
		return
	}
	//	4. 密码正确 创建 token 并返回
	token, err2 := jwt.MakeJWTToken(queryUser)
	if err2 != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, "服务器异常，请稍后再试"))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(token))

}
