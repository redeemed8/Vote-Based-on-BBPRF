package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"math/rand"
	"regexp"
	"time"
)

func MatchPhone(phoneNumber string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phoneNumber)
}

func MatchAccount(account string) bool {
	regexpPattern := `^[a-zA-Z0-9_-]{6,18}$`
	regex := regexp.MustCompile(regexpPattern)
	return regex.MatchString(account)
}

func MysqlErr(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}

func MakeRandStr(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	time.Sleep(time.Nanosecond)
	return string(b)
}

func Len(str string) int {
	if str == "" {
		return 0
	}
	var count = 0
	for i := range str {
		if str[i] > 128 {
			count += 2
		} else {
			count++
		}
	}
	return count
}

func Simply(str string) string {
	if Len(str) <= 24 {
		return str
	}
	var count = 0
	var ret = ""
	var add int
	for _, s := range str {
		if s > 128 {
			add = 2
		} else {
			add = 1
		}
		if count+add <= 24 {
			count += add
			ret += string(s)
		} else {
			break
		}
	}
	return ret + "..."
}

func getBeta(e, u *big.Int, x int64, n int64) int64 {
	ux := new(big.Int).Exp(u, big.NewInt(x), nil)
	t1 := new(big.Int)
	eDivUx := t1.Div(e, ux)
	t2 := new(big.Int)
	eDivUx_1 := t2.Sub(eDivUx, big.NewInt(1))
	t3 := new(big.Int)
	eDivUx_1_mod_n2 := t3.Mod(eDivUx_1, big.NewInt(n*n))
	return eDivUx_1_mod_n2.Int64() / n
}

// extendedEuclidean 算法计算 (gcd, x, y) 使得 a*x + b*y = gcd(a, b)
func extendedEuclidean(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	if b.Cmp(big.NewInt(0)) == 0 {
		return a, big.NewInt(1), big.NewInt(0)
	}

	gcd, x1, y1 := extendedEuclidean(b, new(big.Int).Mod(a, b))
	x := y1
	y := new(big.Int).Sub(x1, new(big.Int).Mul(new(big.Int).Div(a, b), y1))

	return gcd, x, y
}

// ModInverse 计算 a 在模数 m 下的逆元
func ModInverse(a, m *big.Int) (*big.Int, error) {
	if a.Cmp(m) == 0 {
		return nil, fmt.Errorf("求逆失败")
	}
	if a.Cmp(m) == 1 {
		return ModInverse(new(big.Int).Mod(a, m), m)
	}
	gcd, x, _ := extendedEuclidean(a, m)

	if gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, fmt.Errorf("求逆失败")
	}
	return new(big.Int).Mod(x, m), nil
}

func GetBindPRF(e, u *big.Int, x int64, n int64, G *big.Int) (*big.Int, error) {
	beta := getBeta(e, u, x, n)
	beta_mod_n_ie, err := ModInverse(big.NewInt(beta), big.NewInt(n))
	return new(big.Int).Exp(G, beta_mod_n_ie, nil), err
}

func MD5(str string) string {
	data := []byte(str) // 要计算哈希值的数据
	hash := md5.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func H_(uc string, vid int) string {
	s := fmt.Sprintf("%s-%d", uc, vid)
	return MD5(s)
}
func getBinary(num int) [10]int {
	var bins [10]int
	for counter := 0; counter < 10; counter++ {
		bins[counter] = 1 & (num >> counter)
	}
	return bins
}

func checkBins(bins [10]int, imc int) (bool, error) {
	var counter = 0
	for i := 0; i < 10; i++ {
		if bins[i] == 1 {
			counter++
		}
	}
	if counter == 0 {
		return false, errors.New("至少选择一个选项进行投票")
	}
	if imc == 0 && counter != 1 {
		return false, errors.New("只能选择一个选项")
	}
	return true, nil
}

func Analysis(msg int, imc int) (bool, error, [10]int) {
	bins := getBinary(msg)
	ok, err := checkBins(bins, imc)
	if !ok {
		return false, err, bins
	} else {
		return true, nil, bins
	}
}

// DecodeBase64 函数用于解析 Base64 编码字符串并返回原始的字节数组
func DecodeBase64(base64Str string) ([]byte, error) {
	// 解码 Base64 编码字符串
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
