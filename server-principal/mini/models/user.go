package models

import (
	"errors"
	"gorm.io/gorm"
	"mini/config"
	"time"
)

type UserInfoDao_ struct{ db *gorm.DB }

var UserInfoDao UserInfoDao_

func InitUserInfo() {
	UserInfoDao.db = config.DB
	UserInfoDao.CreateTable()
}

type UserInfo struct {
	Id        uint32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Account   string
	Password  string
	Phone     string
	UUID      string
	AvatarUrl string `gorm:"default:''"`
	Nickname  string `gorm:"default:''"`
}

func (user *UserInfo) TableName() string {
	return "5236_user"
}

func (dao *UserInfoDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&UserInfo{})
}

func (dao *UserInfoDao_) CreateUser(user *UserInfo) error {
	return dao.db.Model(&UserInfo{}).Create(user).Error
}

func (dao *UserInfoDao_) UpdateUser(id uint32, updates UserInfo) error {
	return dao.db.Model(&UserInfo{}).Where("id = ?", id).Updates(updates).Error
}

func (dao *UserInfoDao_) GetUserById(id uint32) (UserInfo, error) {
	var queryUser UserInfo
	result := dao.db.Model(&UserInfo{}).Where("id = ?", id).First(&queryUser)
	return queryUser, result.Error
}

func (dao *UserInfoDao_) GetUserByAccount(account string) (UserInfo, error) {
	var queryUser UserInfo
	result := dao.db.Model(&UserInfo{}).Where("account = ?", account).First(&queryUser)
	return queryUser, result.Error
}

func (dao *UserInfoDao_) GetUserByPhone(phone string) (UserInfo, error) {
	var queryUser UserInfo
	result := dao.db.Model(&UserInfo{}).Where("phone = ?", phone).First(&queryUser)
	return queryUser, result.Error
}

func (dao *UserInfoDao_) Identifying(id uint32, phone string, uuid string) (bool, int, error) {
	var user UserInfo
	result := dao.db.Model(&UserInfo{}).Where("id = ? and phone = ? and uuid = ?", id, phone, uuid).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, 500, errors.New("服务器异常，请稍后再试")
	}
	if user.Id < 1 || user.UUID == "" || user.Phone == "" {
		return false, 401, errors.New("请先登录")
	}
	return true, 200, nil
}
