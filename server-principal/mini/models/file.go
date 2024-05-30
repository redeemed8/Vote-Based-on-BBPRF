package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mime/multipart"
	"mini/config"
	"mini/util"
	"mini/util/fileutil"
	"mini/util/timeutil"
	"time"
)

const (
	Photo     = "image"
	UrlPrefix = "https://pvs.81jcpd.cn/"
)

var VideoExtensions = []string{
	"mp4", "avi", "mkv", "wmv", "flv",
	"mov", "webm", "3gp", "ogv", "mpg", "mpeg",
}

var ImageExtensions = []string{
	"jpg", "png", "jpeg", "gif", "bmp", "tif", "tiff",
	"webp", "helc", "helf", "jp2", "j2k", "svg",
}

type OssFileDao_ struct{ db *gorm.DB }
type OssFileUtil_ struct{}

var OssFileDao OssFileDao_
var OssFileUtil OssFileUtil_

func InitOssFile() {
	OssFileDao.db = config.DB
	OssFileDao.CreateTable()
}

type OssFile struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FileId         string `gorm:"primarykey"` //	文件id，采用12位随机串
	FileType       string //	文件类型
	FileExtension  string //	文件后缀名，不加点
	FileMD5        string //	文件的 md5值
	FileObjectName string //	文件的云路径，如 abc/efg/123.jpg
	FileUrl        string //	文件访问的 云 url
	FileSymbol     string //	文件的标识，值就是id，也就是他在临时文件夹下的默认名字
	LocalPath      string //	文件的本地保存名
}

func (table *OssFile) TableName() string {
	return "7346_ossfile"
}

func (dao *OssFileDao_) CreateTable() {
	_ = dao.db.AutoMigrate(&OssFile{})
}

func (dao *OssFileDao_) CreateOssFile(ossFile *OssFile) error {
	return dao.db.Model(&OssFile{}).Create(ossFile).Error
}

// ---------------

func (fileUtil *OssFileUtil_) FileToOssFile(file multipart.File, fileHeader *multipart.FileHeader) (*OssFile, error) {
	//	获取文件类型
	fileType, extension := fileutil.GetFileType(fileHeader.Filename)
	if fileType == fileutil.Unknown {
		return nil, errors.New("未知的文件类型")
	}
	//	创建本地文件的路径
	id := util.MakeRandStr(12)

	localPath := fileutil.TempFilePath + id + extension
	//	创建临时文件
	err1 := fileutil.CreateTempFile(file, localPath)
	if err1 != nil {
		fmt.Println("--------------")
		fmt.Println("创建本地临时文件失败 , err = ", err1)
		fmt.Println("--------------")
		return nil, errors.New("创建本地临时文件失败")
	}
	//	MD5
	md5, err2 := fileutil.GetFileMD5(localPath)
	if err2 != nil {
		fmt.Println("--------------")
		fmt.Println("获取本地文件的md5值出错 , err = ", err2)
		fmt.Println("--------------")
		return nil, errors.New("获取本地文件的md5值出错")
	}
	//  objectName
	year, month, day := timeutil.GetYMDStr() //	获取当前日期
	prefix := md5[:2]
	//	2020/02/15/ec/ecsdjijhdisjdisdj.mp4
	objectName := year + "/" + month + "/" + day + "/" + prefix + "/" + md5 + extension

	return &OssFile{
		FileId:         id,
		FileType:       fileType,
		FileExtension:  extension,
		FileMD5:        md5,
		FileObjectName: objectName,
		FileUrl:        UrlPrefix + objectName,
		FileSymbol:     id,
		LocalPath:      localPath,
	}, nil
}
