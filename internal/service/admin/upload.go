package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zap.L().Error("UploadImage 获取文件失败:" + err.Error())
		response.ErrWithMsg(code.BadRequest, c)
		return
	}
	defer file.Close()

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}

	// 限制大小 10MB
	if header.Size > 10*1024*1024 {
		response.ErrWithMsg(code.BadRequest, c)
		return
	}

	// 创建uploads目录
	uploadDir := "uploads"
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		zap.L().Error("UploadImage 创建目录失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(header.Filename, ext), ext)
	savePath := filepath.Join(uploadDir, filename)

	// 保存文件
	dst, err := os.Create(savePath)
	if err != nil {
		zap.L().Error("UploadImage 创建文件失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		zap.L().Error("UploadImage 写入文件失败:" + err.Error())
		response.ErrWithMsg(code.InternalError, c)
		return
	}

	// 返回URL
	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "上传成功",
		Data: map[string]string{
			"url": url,
		},
	})
}
