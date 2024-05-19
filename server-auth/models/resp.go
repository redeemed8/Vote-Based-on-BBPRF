package models

type StatusCode int
type Resp struct {
	Code StatusCode `json:"code"`
	Data any        `json:"data"`
	Msg  string     `json:"msg"`
}

func NewResp() *Resp {
	return &Resp{}
}

const OK = 200

func (r *Resp) Success(data any) *Resp {
	r.Code = OK
	r.Msg = "success"
	r.Data = data
	return r
}

func (r *Resp) Fail(code int, msg string) *Resp {
	r.Code = StatusCode(code)
	r.Msg = msg
	return r
}
