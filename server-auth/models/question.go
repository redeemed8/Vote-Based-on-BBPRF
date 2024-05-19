package models

import (
	"gorm.io/gorm"
	"mini/config"
)

const (
	Select       = "select"
	Blank        = "blank"
	OptionsSplit = "$%&"
	VoidOption   = "void"
)

type Question struct {
	Id       string `gorm:"primarykey"` //	问题id
	FatherId string //	问卷id
	Title    string //	问题标题
	Type     string //	"select" or "blank"  选择还是填空
	Options  string //	选项  :  每个选项之间用 $%& 隔开
}

type QuestionDao_ struct{ db *gorm.DB }
type QuestionUtil_ struct{}

var QuestionDao QuestionDao_
var QuestionUtil QuestionUtil_

func InitQuestion() {
	QuestionDao.db = config.DB
	QuestionDao.CreateTable()
}

func (question *Question) TableName() string {
	return "4546_question"
}

func (dao *QuestionDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&Question{})
}

func (dao *QuestionDao_) CreateQuestions(question []Question) error {
	return dao.db.Model(&Question{}).Create(&question).Error
}

func (dao *QuestionDao_) GetQuestion(id string, fid string) (Question, error) {
	var q Question
	result := dao.db.Model(&Question{}).Where("id = ? and father_id = ?", id, fid).First(&q)
	return q, result.Error
}

// ---------------------

func (util *QuestionUtil_) GetOptionSplit() string {
	return OptionsSplit
}
