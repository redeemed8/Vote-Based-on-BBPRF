package models

import (
	"gorm.io/gorm"
	"mini/config"
)

type HDao_ struct{ db *gorm.DB }

var HDao HDao_

func InitH() {
	HDao.db = config.DB
	HDao.CreateTable()
}

type H struct {
	Name string `gorm:"primarykey"`
	Data int
}

func (h *H) TableName() string {
	return "3471_h"
}

func (dao *HDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&H{})
}

func (dao *HDao_) CreateH(h *H) error {
	return dao.db.Model(&H{}).Create(h).Error
}

func (dao *HDao_) GetHByName(name string) (H, error) {
	var h H
	result := dao.db.Model(&H{}).Where("name = ?", name).First(&h)
	return h, result.Error
}

func (dao *HDao_) AddDataByName(name string) error {
	return dao.db.Model(&H{}).Where("name = ? and data = 1", name).UpdateColumn("data", gorm.Expr("data + ?", 1)).Error
}
