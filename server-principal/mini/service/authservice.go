package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/big"
	"mini/models"
	"mini/models/dto"
	"mini/util"
	"mini/util/jwt"
	"net/http"
)

// VerifyJWTToken
// 验证一个token是不是jwt-token	[get]   LOGIN
func VerifyJWTToken(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 校验登录
	userClaim, login := jwt.ParseToken(ctx, resp)
	if !login {
		return //	返回的 code == 401
	}
	ctx.JSON(http.StatusOK, resp.Success(userClaim.UUID))
}

// GetPprmS
// 获取服务端公钥 [get]	NOT LOGIN
func GetPprmS(ctx *gin.Context) {
	resp := models.NewResp()
	ctx.JSON(http.StatusOK, resp.Success(dto.PprmSDto{
		N:   models.PprmS.N,
		G:   models.PprmS.G.String(),
		Y:   models.PprmS.Y.String(),
		H:   models.PprmS.H,
		CtU: [2]string{models.PprmS.CtU[0].String(), models.PprmS.CtU[1].String()},
		CtY: [2]string{models.PprmS.CtY[0].String(), models.PprmS.CtY[1].String()},
	}))
}

type SignVo struct {
	U      string `json:"u"`
	E      string `json:"e"`
	Uc     string `json:"uc"`
	UcSign string `json:"uc_sign"`
	Vid    int    `json:"vid"`
}

// Sign
// 服务器端对消息进行签名处理 	[post]	不登录，使用客户端密钥进行签名
// post_args : {客户端密钥, 客户端密钥的签名}
func Sign(ctx *gin.Context) {
	//	1. .....
	resp := models.NewResp()
	//	2. 绑定参数
	var signVo SignVo
	if err := ctx.ShouldBind(&signVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	3. 参数检验
	if signVo.U == "" || signVo.E == "" || signVo.Uc == "" || signVo.UcSign == "" || signVo.Vid < 1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "参数校验失败"))
		return
	}
	signFromBase64, errB := util.DecodeBase64(signVo.UcSign)
	if errB != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "数据有误"))
		return
	}
	//	4. 校验客户端密钥签名
	valid := PublicKey.Verify([]byte(util.MD5(signVo.Uc)), signFromBase64)
	if !valid {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权投票"))
		return
	}
	//	5. 签名验证成功, 检查 vid 是否真实
	queryV, queryE := models.VoteDao.GetVote(models.Vote{Id: signVo.Vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryV.Title == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	if queryV.Status == 0 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票尚未发布"))
		return
	}
	//	6. vid 真实，检查  H(uc,vid) 是否存在 , 存在则不给予证明
	H_sign := util.H_(signVo.Uc, signVo.Vid)
	queryH, err22 := models.HDao.GetHByName(H_sign)
	if util.MysqlErr(err22) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryH.Data != 0 && queryH.Data != 1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "您已经投过票了"))
		return
	}
	//	7. 不存在，给予证明 F , 并异步的 添加 H(uc,vid)
	e, ok1 := big.NewInt(1).SetString(signVo.E, 10)
	if !ok1 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效投票"))
		return
	}
	u, ok2 := big.NewInt(1).SetString(signVo.U, 10)
	if !ok2 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效投票"))
		return
	}
	F, errF := util.GetBindPRF(e, u, models.PrivprmS.X, models.PprmS.N, models.PprmS.G)
	if errF != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, errF.Error()))
		return
	}
	HChan <- models.H{Name: H_sign}
	ctx.JSON(http.StatusOK, resp.Success(gin.H{"F": F.String(), "h_sign": H_sign}))
}

type VerifyVo struct {
	Msg   int    `json:"msg"`
	R     int    `json:"r"`
	Uc    string `json:"uc"`
	Vid   int    `json:"vid"`
	Token string `json:"token"`
	HSign string `json:"h_sign"`
}

type V struct {
	Bins   [10]int
	HSign  string
	Imc    int
	QueryV models.Vote
}

// Verify
// 服务器端对token进行验证	 [post]	 LOGIN  带token即可
func Verify(ctx *gin.Context) {
	resp := models.NewResp()
	//	1. 接收参数
	var verifyVo VerifyVo
	if err := ctx.ShouldBind(&verifyVo); err != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, "无效参数"))
		return
	}
	//	2. 验证 HSign
	if util.H_(verifyVo.Uc, verifyVo.Vid) != verifyVo.HSign {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权投票"))
		return
	}
	//	3. 验证token  token == G ** ?
	m_u_ry := verifyVo.Msg + int(models.PrivprmS.U) + verifyVo.R*int(models.PrivprmS.Y)
	m_u_ry_mod_n_ie, err1 := util.ModInverse(big.NewInt(int64(m_u_ry)), big.NewInt(models.PprmS.N))
	if err1 != nil {
		ctx.JSON(http.StatusOK, resp.Fail(400, err1.Error()))
		return
	}
	G_m_u_ry := new(big.Int).Exp(models.PprmS.G, m_u_ry_mod_n_ie, nil)
	//	4. 比较token
	if verifyVo.Token != G_m_u_ry.String() {
		ctx.JSON(http.StatusOK, resp.Fail(400, "你无权投票"))
		return
	}
	//	5. 可以投票, 异步为其添加投票结果 并 处理 H_sign
	queryV, queryE := models.VoteDao.GetVote(models.Vote{Id: verifyVo.Vid})
	if util.MysqlErr(queryE) {
		ctx.JSON(http.StatusOK, resp.Fail(500, ServerErr))
		return
	}
	if queryV.Title == "" {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票不存在或已被删除"))
		return
	}
	if queryV.Status == 0 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票尚未发布"))
		return
	}
	//	6. ...........
	if verifyVo.Msg >= 1024 {
		ctx.JSON(http.StatusOK, resp.Fail(400, "投票无效"))
		return
	}
	ok, errA, bins := util.Analysis(verifyVo.Msg, queryV.IsMultiChoice)
	if !ok {
		ctx.JSON(http.StatusOK, resp.Fail(400, errA.Error()))
		return
	}

	VChan <- V{
		Bins:   bins,
		HSign:  verifyVo.HSign,
		Imc:    queryV.IsMultiChoice,
		QueryV: queryV,
	}
	ctx.JSON(http.StatusOK, resp.Success("投票成功,感谢你的参与"))
}

var HChan = make(chan models.H)

func CreateH() {
	fmt.Println("任务 --- 创建H 已开启")
	for {
		select {
		case h, ok := <-HChan:
			if ok {
				h.Data = 1
				_ = models.HDao.CreateH(&h)
			}
		}
	}
}

var AddHChan = make(chan string)

func AddH() {
	fmt.Println("任务 --- 增加H 已开启")
	for {
		select {
		case name, ok := <-AddHChan:
			if ok {
				_ = models.HDao.AddDataByName(name)
			}
		}
	}
}

var VChan = make(chan V)

func DealWithV() {
	fmt.Println("任务 --- 处理投票 已开启")
	for {
		select {
		case v, ok := <-VChan:
			if ok {
				//	 先更新选项数
				newAnsCount := models.VoteUtil.MergeOptionsToCountStr(v.QueryV.AnsCount, v.Bins[:], v.QueryV.IsMultiChoice)
				err := models.VoteDao.UpdateVote(v.QueryV.Id, map[string]interface{}{
					"participants": v.QueryV.Participants + 1,
					"ans_count":    newAnsCount})
				if util.MysqlErr(err) {
					continue
				}
				//	再更新 h_data
				_ = models.HDao.AddDataByName(v.HSign)
			}
		}
	}
}
