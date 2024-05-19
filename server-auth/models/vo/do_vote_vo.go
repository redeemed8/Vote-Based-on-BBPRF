package vo

import "errors"

type DoVoteVo struct {
	Vid           int   `json:"vid"`
	Status        int   `json:"status"`
	IsMultiChoice int   `json:"imc"`
	SelectOptions []int `json:"select_options"`
	//  .............
}

func (doVoteVo *DoVoteVo) Right() (bool, error) {
	if doVoteVo.Vid < 1 {
		return false, errors.New("投票不存在或已被删除")
	}
	if doVoteVo.Status == 0 {
		return false, errors.New("投票尚未发布")
	}
	if doVoteVo.Status != 1 {
		doVoteVo.Status = 1
	}
	if doVoteVo.IsMultiChoice != 0 && doVoteVo.IsMultiChoice != 1 {
		doVoteVo.IsMultiChoice = 1
	}
	if len(doVoteVo.SelectOptions) < 2 {
		return false, errors.New("至少应有2个选项")
	}
	counter := 0
	for _, select_ := range doVoteVo.SelectOptions {
		if select_ != 0 {
			counter++
		}
	}
	if counter < 1 {
		return false, errors.New("至少应该选择一个选项")
	}
	return true, nil
}
