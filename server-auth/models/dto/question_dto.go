package dto

type QuestionDto struct {
	Id       string   `json:"id"`
	FatherId string   `json:"father_id"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Options  []string `json:"options"`
}
