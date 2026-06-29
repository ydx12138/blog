package admin

import (
	"blog/config"
	"blog/internal/utils"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 允许的图片类型
var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	//接收图片
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zap.L().Error("UploadImage 获取文件失败:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			zap.L().Error("UploadImage" + err.Error())
		}
	}(file)

	// 校验扩展名，校验扩展名是不是allowedExts里的一种
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		zap.L().Error("图片扩展名错误")
		response.ErrWithMsg(code.BadRequest, c)
		return
	}

	// 限制大小 10MB
	if header.Size > 10*1024*1024 {
		zap.L().Error("图片大小应小于等于10MB")
		response.ErrWithMsg(code.BadRequest, c)
		return
	}

	/*// 创建uploads目录
	uploadDir := "uploads"
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		zap.L().Error("UploadImage 创建目录失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}*/

	// 生成唯一文件名xxx.yyy
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(header.Filename, ext), ext)
	url := utils.UploadToOss(file, config.Cfg.OssConfig.Image_path, filename)
	//拼接目录和文件名
	//savePath := filepath.Join(uploadDir, filename)

	/*//保存文件，文件存在就替换，不存在就创建
	dst, err := os.Create(savePath)
	if err != nil {
		zap.L().Error("UploadImage 创建文件失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			zap.L().Error("UploadImage" + err.Error())
		}
	}(dst)
	//之前创造的时空白文件，现在把图片复制过去
	_, err = io.Copy(dst, file)
	if err != nil {
		zap.L().Error("UploadImage 写入文件失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}*/

	// 返回URL
	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "上传成功",
		Data: map[string]string{
			"url": url,
		},
	})
}
