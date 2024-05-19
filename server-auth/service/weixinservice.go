package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mini/models"
	"net/http"
)

var (
	appId     = "wx089b65ae0c4a3957"
	appSecret = "a67ab882cebc27c4e27e3bf60e175388"
	grantType = "authorization_code"
)

type OpenIdResp struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
}

func GetOpenId(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 获取路径参数 js_code
	jsCode := ctx.Query("code")
	if jsCode == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "code不能为空"))
		return
	}
	//	2. 拼接 url
	url := "https://api.weixin.qq.com/sns/jscode2session?"
	url += "appid=" + appId + "&" + "secret=" + appSecret + "&" + "grant_type=" + grantType + "&"
	url += "js_code=" + jsCode

	//	3. 发起 HTTPS 请求
	r, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, "请求失败"))
		return
	}
	defer r.Body.Close()

	//	4. 读取响应内容
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, "读取响应失败"))
		return
	}
	//	5. 解析
	var openIdResp OpenIdResp
	if err2 := json.Unmarshal(body, &openIdResp); err2 != nil || openIdResp.Openid == "" {
		ctx.JSON(http.StatusOK, resp.Fail(500, "获取id失败"))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(openIdResp.Openid))
}
