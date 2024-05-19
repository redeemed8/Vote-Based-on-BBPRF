package dto

type VoteDto struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Participants int    `json:"participants"` //	参与人数
	Status       int    `json:"status"`       //	是否发布
}

type GetVoteDto struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	IsMultiChoice int      `json:"is_multi_choice"`
	Status        int      `json:"status"`
	Options       []string `json:"options"`

	Participants int   `json:"participants"`
	AnsCount     []Ans `json:"ans_count"`
}

type Ans struct {
	OptionName  string `json:"option_name"`
	OptionCount int    `json:"option_count"`
}
