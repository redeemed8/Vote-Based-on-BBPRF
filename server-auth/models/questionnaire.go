package models

import (
	"gorm.io/gorm"
	"mini/config"
	"mini/models/dto"
	"mini/util"
)

type Questionnaire struct {
	Id          string `gorm:"primarykey"`
	Title       string //	问卷标题
	PublisherId uint32
	QuestionIds string //	所有问题id的集合  形如  "1,2,3"
	AnsNumber   int    //	答卷数量
	Status      int    //	发布状态
}

type QuestionnaireDao_ struct{ db *gorm.DB }
type QuestionnaireUtil_ struct{}

var QuestionnaireDao QuestionnaireDao_
var QuestionnaireUtil QuestionnaireUtil_

func InitQuestionnaire() {
	QuestionnaireDao.db = config.DB
	QuestionnaireDao.CreateTable()
}

func (questionnaire *Questionnaire) TableName() string {
	return "1634_questionnaire"
}

type Questionnaires []Questionnaire

func (Questionnaires *Questionnaires) ToDtos(isSimple bool) []dto.QuestionnaireDto {
	var dtos = make([]dto.QuestionnaireDto, 0)
	for i := len(*Questionnaires) - 1; i >= 0; i-- {
		var dto_ dto.QuestionnaireDto
		dto_.Id = (*Questionnaires)[i].Id
		if isSimple {
			dto_.Title = util.Simply((*Questionnaires)[i].Title)
		} else {
			dto_.Title = (*Questionnaires)[i].Title
		}
		dto_.AnsNumber = (*Questionnaires)[i].AnsNumber
		dto_.Status = (*Questionnaires)[i].Status
		dtos = append(dtos, dto_)
	}
	return dtos
}

func (dao *QuestionnaireDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&Questionnaire{})
}
func (dao *QuestionnaireDao_) CreateQuestionnaire(questionnaire Questionnaire) error {
	return dao.db.Model(&Questionnaire{}).Create(questionnaire).Error
}

func (dao *QuestionnaireDao_) SearchQuestionnaire(id uint32, keyword string) (Questionnaires, error) {
	var questionnaires = make(Questionnaires, 0)
	result := dao.db.Model(&Questionnaire{}).Where("publisher_id = ? and title like ?", id, "%"+keyword+"%").Find(&questionnaires)
	return questionnaires, result.Error
}

func (dao *QuestionnaireDao_) GetQuestionnaire(qid string) (Questionnaire, error) {
	var ret Questionnaire
	result := dao.db.Model(&Questionnaire{}).Where("id = ?", qid).First(&ret)
	return ret, result.Error
}

func (dao *QuestionnaireDao_) DelQuestionnaire(id string, publisherId uint32) error {
	return dao.db.Model(&Questionnaire{}).Where("id = ? and publisher_id = ?", id, publisherId).Delete(&Questionnaire{}).Error
}
func (dao *QuestionnaireDao_) UptQuestionnaire(id string, publisherId uint32, updates map[string]interface{}) error {
	return dao.db.Model(&Questionnaire{}).Where("id = ? and publisher_id = ?", id, publisherId).Updates(updates).Error
}

// -----------------------------------------------------

func (util *QuestionnaireUtil_) GetIdsSplit() string {
	return ","
}

func (util *QuestionnaireUtil_) UpdateStatus(status int) int {
	if status == 0 {
		return 1
	} else {
		return 0
	}
}
