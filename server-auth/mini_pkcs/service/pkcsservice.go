package service

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io/ioutil"
	"mini_pkcs/model"
	"mini_pkcs/util"
	"net/http"
)

const Password = "huihuchsahkjffjiopwofjopgg4545wf44d5s4d56sa4d4w844cxc12sa1dc53a1d51d1231c21d351d5aq41d854w5dwd45dqwqwd4dq4f2"

const RequestHeader = "https://mini.81jcpd.cn"

//const RequestHeader = "http://localhost:3656"

// GetSM2PublicKey
// 获取sm2的公钥		[post]  	pwd
func GetSM2PublicKey(ctx *gin.Context) {
	//	获取密码的md5值
	pwd := ctx.PostForm("pwd")
	if pwd == "" {
		ctx.JSON(http.StatusOK, gin.H{"x": "", "y": ""})
		return
	}
	//	验证密码
	if util.MD5(pwd) != util.MD5(Password) {
		ctx.JSON(http.StatusOK, gin.H{"x": "", "y": ""})
		return
	}
	//	密码正确，返回公钥
	ctx.JSON(http.StatusOK, gin.H{"x": model.PublicKeyX, "y": model.PublicKeyY})
}

const ServerError = "服务器异常，请稍后再试"

type TokenDataResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func verifyToken(token string) (bool, int, error, string) {
	client := &http.Client{}
	//	构造请求
	req, err := http.NewRequest("GET", RequestHeader+"/auth/verify/token", nil)
	if err != nil {
		return false, 500, errors.New(ServerError), ""
	}
	//	添加请求头
	req.Header.Set("X-auth", token)
	//	发起请求
	resp, err2 := client.Do(req)
	if err2 != nil {
		return false, 500, errors.New(ServerError), ""
	}
	defer resp.Body.Close()
	//	拿出响应体
	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return false, 500, errors.New(ServerError), ""
	}
	//	用结构体接收
	var respData TokenDataResp

	err4 := json.Unmarshal(respBody, &respData)
	if err4 != nil {
		return false, 500, errors.New(ServerError), ""
	}
	//	判断是否成功
	if respData.Code == 200 {
		return true, 200, nil, respData.Data
	} else {
		return false, respData.Code, errors.New(respData.Msg), ""
	}
}

// GetSignedPrivKey
// 获取客户端私钥和签名		[post]
func GetSignedPrivKey(ctx *gin.Context) {
	resp := model.NewResp()
	//	1. 获取 token
	identify := ctx.PostForm("identify") // jwt-token
	if identify == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "no"))
		return
	}
	//	2. 验证是否是 jwt-token
	right, code, err, uuid_ := verifyToken(identify)
	if !right {
		ctx.JSON(http.StatusOK, resp.Fail(code, err.Error()))
		return
	}
	//	3. 将 jwt-token 哈希后，检查数据库中是否有其的密钥
	tokenMd5 := util.MD5(uuid_)
	queryPrivKey, err1 := model.PrivKeyInfoDao.GetPrivKeyInfoByToken(tokenMd5)
	if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerError))
		return
	}
	if queryPrivKey != "" {
		signTokenMd5, err9 := sm2Sign(util.MD5(queryPrivKey))
		if err9 != nil {
			ctx.JSON(http.StatusOK, resp.Fail(500, ServerError))
			return
		}
		ctx.JSON(http.StatusOK, resp.Success(gin.H{
			"priv_key": queryPrivKey,
			"sign":     signTokenMd5,
		}))
		return
	}
	//	4. 库里没有这个token，为其创建一个新的
	newPrivKey := uuid.New().String()
	info := model.PrivKeyInfo{TokenMd5: tokenMd5, ClientPrivKey: newPrivKey}
	err2 := model.PrivKeyInfoDao.CreatePrivKeyInfo(&info)
	if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerError))
		return
	}
	//	5. 创建成功返回 新的 私钥
	signTokenMd5, err9 := sm2Sign(util.MD5(newPrivKey))
	if err9 != nil {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerError))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(gin.H{
		"priv_key": newPrivKey,
		"sign":     signTokenMd5,
	}))
}

func sm2Sign(data string) (string, error) {
	message := []byte(data)
	// 使用私钥进行签名
	signature, err := model.PrivateKey.Sign(rand.Reader, message, nil)
	if err != nil {
		return "", errors.New("1")
	}

	fmt.Println("------------------------------")
	fmt.Println("签名值为 base64 == ", util.EncodeBase64(signature))
	fmt.Println("------------------------------")

	return util.EncodeBase64(signature), nil
}
