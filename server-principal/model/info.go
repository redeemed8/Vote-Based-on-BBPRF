package model

import (
	"gorm.io/gorm"
	"mini_pkcs/config"
)

type PrivKeyInfoDao_ struct{ db *gorm.DB }

var PrivKeyInfoDao PrivKeyInfoDao_

func InitInfo() {
	PrivKeyInfoDao.db = config.DB
	PrivKeyInfoDao.CreateTable()
}

type PrivKeyInfo struct {
	Id            int `gorm:"primarykey"`
	TokenMd5      string
	ClientPrivKey string
}

func (pkinfo *PrivKeyInfo) TableName() string {
	return "7831_pkinfo"
}

func (dao *PrivKeyInfoDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&PrivKeyInfo{})
}

func (dao *PrivKeyInfoDao_) CreatePrivKeyInfo(privKeyInfo *PrivKeyInfo) error {
	return dao.db.Model(&PrivKeyInfo{}).Create(privKeyInfo).Error
}

func (dao *PrivKeyInfoDao_) GetPrivKeyInfoByToken(tokenMd5 string) (string, error) {
	var r string
	result := dao.db.Model(&PrivKeyInfo{}).Select("client_priv_key").Where("token_md5 = ?", tokenMd5).First(&r)
	return r, result.Error
}
