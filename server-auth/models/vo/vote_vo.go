package vo

import (
	"errors"
	"fmt"
	"mini/models"
	"mini/util"
	"strings"
)

type NewVoteVo struct {
	VoteTitle     string   `json:"vote_title"`
	Options       []string `json:"options"`
	IsMultiChoice int      `json:"is_multi_choice"` //	1 是-多选   0 否-单选
}

func (newVoteVo NewVoteVo) Right() (bool, error) {
	if newVoteVo.VoteTitle == "" {
		return false, errors.New("标题不能为空")
	}
	if newVoteVo.IsMultiChoice != 1 {
		newVoteVo.IsMultiChoice = 0
	}
	//	遍历所有的选项，先排除空选项
	options := make([]string, 0)
	for i, option := range newVoteVo.Options {
		if option != "" {
			options = append(options, option)
		}
		if util.Len(option) > 140 {
			return false, errors.New(fmt.Sprintf("第%d个选项长度超出限制", i+1))
		}
	}
	//	然后再进行校验
	if len(options) < 2 {
		return false, errors.New("投票至少应有两个选项")
	}
	newVoteVo.Options = options
	return true, nil
}

func (newVoteVo NewVoteVo) Parse(userId uint32) models.Vote {
	return models.Vote{
		PublisherId:   userId,
		Title:         newVoteVo.VoteTitle,
		Options:       strings.Join(newVoteVo.Options, models.VoteUtil.GetOptionSplit()),
		IsMultiChoice: newVoteVo.IsMultiChoice,
		Participants:  0,
		AnsCount:      models.VoteUtil.InitAnsCount(newVoteVo.Options),
		Status:        0,
	}
}
