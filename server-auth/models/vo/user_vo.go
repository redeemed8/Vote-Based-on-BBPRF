package vo

type PhoneLoginVo struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type RegisterAccountVo struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Repassword string `json:"repassword"`
	Phone      string `json:"phone"`
	Code       string `json:"code"`
}

type AccountLoginVo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
