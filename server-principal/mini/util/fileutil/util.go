package fileutil

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"mini/util"
	"os"
	"path/filepath"
)

const (
	Photo        = "image"
	Unknown      = "unknown"
	UrlPrefix    = "https://pvs.81jcpd.cn/"
	TempFilePath = "temp/pvs/"
)

var ImageExtensions = []string{
	".jpg", ".png", ".jpeg", ".gif", ".bmp", ".tif", ".tiff",
	".webp", ".helc", ".helf", ".jp2", ".j2k", ".svg",
}

// GetFileType 根据文件名获取文件类型
func GetFileType(filename string) (string, string) {
	ext := filepath.Ext(filename)
	if util.StringsContain(ImageExtensions, ext) {
		return Photo, ext
	}
	return Unknown, ext
}

func CreateTempFile(file multipart.File, filepath string) error {
	// 创建本地文件
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	// 将上传的文件内容复制到本地文件中
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

// GetFileMD5 计算文件的 MD5 哈希值
func GetFileMD5(filePath string) (string, error) {
	// 1.打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	// 2.创建 MD5 哈希对象
	hash := md5.New()
	// 3.将文件内容传入哈希对象
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	// 4.计算哈希值并转为字符串
	hashInBytes := hash.Sum(nil)
	md5String := hex.EncodeToString(hashInBytes)
	return md5String, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
