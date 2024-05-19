package vo

import (
	"errors"
	"mini/models"
	"mini/util"
	"strings"
)

type QuestionVo struct {
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
}

type NewQuestionnaireVo struct {
	Title     string       `json:"title"`
	Questions []QuestionVo `json:"questions"`
}

func (newQuestionnaireVo *NewQuestionnaireVo) Right() (bool, error) {
	if newQuestionnaireVo.Title == "" {
		return false, errors.New("问卷标题不能为空")
	}
	if len(newQuestionnaireVo.Questions) < 1 {
		return false, errors.New("问卷至少有1个问题")
	}
	for i, question := range newQuestionnaireVo.Questions {
		if question.Title == "" {
			return false, errors.New("问题必须有1个标题")
		}
		if question.Type != models.Blank && question.Type != models.Select {
			return false, errors.New("无效的问题类型")
		}
		if question.Title == models.Select {
			if len(question.Options) < 2 {
				return false, errors.New("选择题至少有2个选项")
			}
			if len(question.Options) > 10 {
				return false, errors.New("选择题最多有10个选项")
			}
		}
		if question.Type == models.Blank {
			newQuestionnaireVo.Questions[i].Options = []string{models.VoidOption}
		}
	}
	return true, nil
}

func (newQuestionnaireVo *NewQuestionnaireVo) Parse(userId uint32) (models.Questionnaire, []models.Question) {
	var questionnaire models.Questionnaire
	var questions = make([]models.Question, 0)

	//	生成问卷id
	questionnaireId := util.MakeRandStr(16)
	//	填写问卷id 、 标题 、 发布人id
	questionnaire.Id = questionnaireId
	questionnaire.Title = newQuestionnaireVo.Title
	questionnaire.PublisherId = userId
	//	生成用于保存 问卷中 所有问题id 的集合
	questionIds := make([]string, 0)
	//	选项的分隔符
	optionSplit := "$%&"
	//	问题id的分隔符
	idSplit := ","
	//	遍历所有的问题
	for _, q := range newQuestionnaireVo.Questions {
		var question models.Question                   //	新建问题
		question.Id = util.MakeRandStr(6)              //	 问题随机id
		questionIds = append(questionIds, question.Id) //	保存问题id 到列表

		question.FatherId = questionnaireId                     //	问卷id
		question.Title = q.Title                                //	标题
		question.Type = q.Type                                  //	问题类型
		question.Options = strings.Join(q.Options, optionSplit) //	问题选项
		questions = append(questions, question)                 //	添加问题到列表
	}
	//	填写问卷的所有问题id
	questionnaire.QuestionIds = strings.Join(questionIds, idSplit)

	return questionnaire, questions
}
