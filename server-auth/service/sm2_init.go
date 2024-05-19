package service

import (
	"encoding/json"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
)

const Password = "huihuchsahkjffjiopwofjopgg4545wf44d5s4d56sa4d4w844cxc12sa1dc53a1d51d1231c21d351d5aq41d854w5dwd45dqwqwd4dq4f2"

const RequestHeader = "https://calc.81jcpd.cn"

//const RequestHeader = "http://localhost:4847"

type PublicKeyResp struct {
	X string `json:"x"`
	Y string `json:"y"`
}

var PublicKey sm2.PublicKey

// 获取签名公钥
func init() {
	// 1.准备 POST 参数，以表单形式传递
	postData := url.Values{}
	postData.Set("pwd", Password)
	// 2.发送 POST 请求
	resp, err := http.PostForm(RequestHeader+"/pkcs/get/sm2pk", postData)
	if err != nil {
		fmt.Println("获取公钥的请求 发送 失败")
		os.Exit(1)
	}
	defer resp.Body.Close()
	//	3.拿出响应体
	respBody, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Println("获取公钥的请求 获取响应体 失败")
		os.Exit(1)
	}
	//	4.用变量接收
	var publicKeyResp PublicKeyResp
	err4 := json.Unmarshal(respBody, &publicKeyResp)
	if err4 != nil {
		fmt.Println("获取公钥的请求 用变量接收 失败")
		os.Exit(1)
	}
	//	5.组成公钥
	x := new(big.Int)
	y := new(big.Int)
	x.SetString(publicKeyResp.X, 10)
	y.SetString(publicKeyResp.Y, 10)
	PublicKey = sm2.PublicKey{Curve: sm2.P256Sm2(), X: x, Y: y}

	fmt.Println("sm2 公钥已加载 ...")
}
