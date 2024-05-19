package dto

type QuestionnaireDto struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	AnsNumber int    `json:"ans_number"`
	Status    int    `json:"status"`
}
