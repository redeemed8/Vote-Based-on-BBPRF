package models

import (
	"fmt"
	"gorm.io/gorm"
	"mini/config"
	"mini/models/dto"
	"mini/util"
	"strconv"
	"strings"
)

type VoteDao_ struct{ db *gorm.DB }
type VoteUtil_ struct{}

var VoteDao VoteDao_
var VoteUtil VoteUtil_

func InitVote() {
	VoteDao.db = config.DB
	VoteDao.CreateTable()
}

type Vote struct {
	Id            int `gorm:"primarykey"`
	PublisherId   uint32
	Title         string
	Options       string `gorm:"type:text"` //	 至少两个选项
	IsMultiChoice int
	Participants  int    //	参与人数
	AnsCount      string // "a-1;b-3;c-2"   "a-9;b-1"
	Status        int    //	是否发布
}

const (
	AnsSplit        = ";"
	VoteOptionSplit = "^&"
)

func (vote *Vote) TableName() string {
	return "7965_vote"
}

func (dao *VoteDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&Vote{})
}

func (dao *VoteDao_) CreateVote(vote *Vote) error {
	return dao.db.Model(&Vote{}).Create(vote).Error
}

func (dao *VoteDao_) GetVote(vote Vote) (Vote, error) {
	var v Vote
	result := dao.db.Model(&Vote{}).Where(vote).First(&v)
	return v, result.Error
}

func (dao *VoteDao_) DelVoteById(id int) error {
	return dao.db.Model(&Vote{}).Where("id = ?", id).Delete(&Vote{}).Error
}

type Votes []Vote

func (votes *Votes) ToDtos(isSimply bool) []dto.VoteDto {
	var dtos = make([]dto.VoteDto, 0)
	for _, vote := range *votes {
		if isSimply {
			vote.Title = util.Simply(vote.Title)
		}
		dtos = append(dtos, dto.VoteDto{
			Id:           vote.Id,
			Title:        vote.Title,
			Participants: vote.Participants,
			Status:       vote.Status,
		})
	}
	return dtos
}

func (dao *VoteDao_) SearchVote(userId uint32, keyword string) (Votes, error) {
	var votes = make(Votes, 0)
	result := dao.db.Model(&Vote{}).Where("publisher_id = ? and title like ?", userId, "%"+keyword+"%").Find(&votes)
	return votes, result.Error
}

func (dao *VoteDao_) UpdateVote(id int, updates map[string]interface{}) error {
	return dao.db.Model(&Vote{}).Where("id = ?", id).Updates(updates).Error
}

//	--------------------------------------

func (util *VoteUtil_) GetOptionSplit() string {
	return VoteOptionSplit
}

func (util *VoteUtil_) GetAnsSplit() string {
	return AnsSplit
}

func (util *VoteUtil_) InitAnsCount(options []string) string {
	var NonEmptyOptions = make([]string, 0)
	for _, option := range options {
		if option != "" {
			NonEmptyOptions = append(NonEmptyOptions, option+"-0")
		}
	}
	return strings.Join(NonEmptyOptions, util.GetAnsSplit())
}

func (util *VoteUtil_) countStrToMap(str string) (map[string]int, []string) {
	//	a-1;b-2;c-123
	var optionNames = make([]string, 0)
	var optionMap = make(map[string]int)
	options := strings.Split(str, util.GetAnsSplit())
	for _, option := range options {
		if option == "" {
			continue
		}
		arr := strings.Split(option, "-")
		if len(arr) != 2 {
			continue
		}
		optionName := arr[0]
		optionNum, err := strconv.Atoi(arr[1])
		if err != nil || optionNum < 0 {
			continue
		}
		optionMap[optionName] = optionNum
		optionNames = append(optionNames, optionName)
	}
	return optionMap, optionNames
}

func (util *VoteUtil_) mapToCountStr(map_ map[string]int, optionNames []string) string {
	var ret = ""
	for _, optionName := range optionNames {
		ret += fmt.Sprintf("%s-%d", optionName, map_[optionName])
		ret += util.GetAnsSplit()
	}
	return ret[:len(ret)-1]
}

func (util *VoteUtil_) addAnsCount(countStr string, optionName string) string {
	optionMap, optionNames := util.countStrToMap(countStr)
	if _, exist := optionMap[optionName]; exist {
		optionMap[optionName]++
	}
	return util.mapToCountStr(optionMap, optionNames)
}

func (util *VoteUtil_) subAnsCount(countStr string, optionName string) string {
	optionMap, optionNames := util.countStrToMap(countStr)
	if _, exist := optionMap[optionName]; exist && optionMap[optionName] > 0 {
		optionMap[optionName]--
	}
	return util.mapToCountStr(optionMap, optionNames)
}

const (
	Add = 1
	Sub = 2
)

func (util *VoteUtil_) DealWithAnsCount(countStr string, optionName string, way int) string {
	if way == Add {
		return util.addAnsCount(countStr, optionName)
	} else if way == Sub {
		return util.subAnsCount(countStr, optionName)
	}
	return countStr
}

func (util *VoteUtil_) ParseAnsCount(countStr string) []dto.Ans {
	var dtos = make([]dto.Ans, 0)
	optionMap, optionNames := util.countStrToMap(countStr)
	for _, optionName := range optionNames {
		if optionName == "" {
			continue
		}
		var ans = dto.Ans{
			OptionName:  optionName,
			OptionCount: optionMap[optionName],
		}
		dtos = append(dtos, ans)
	}
	return dtos
}

func (util *VoteUtil_) Options(options string) []string {
	var arr = make([]string, 0)
	optionArr := strings.Split(options, util.GetOptionSplit())
	for _, option := range optionArr {
		if option == "" {
			continue
		}
		arr = append(arr, option)
	}
	return arr
}

func (util *VoteUtil_) UpdateStatus(status int) int {
	if status == 0 {
		return 1
	} else {
		return 0
	}
}

func (util *VoteUtil_) MergeOptionsToCountStr(ansCountStr string, selectedIdx []int, imc int) string {
	ansCountMap, optionNames := util.countStrToMap(ansCountStr)
	if len(optionNames) < 2 {
		return ansCountStr
	}
	for i := 0; i < len(selectedIdx); i++ {
		if selectedIdx[i] == 0 {
			continue
		}
		if i > len(optionNames)-1 {
			break
		}
		ansCountMap[optionNames[i]] = ansCountMap[optionNames[i]] + 1
		if imc == 0 {
			break
		}
	}
	return util.mapToCountStr(ansCountMap, optionNames)
}
