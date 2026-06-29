package utils

import (
	"blog/config"
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func UploadToOss(body multipart.File, path string, filename string) string {
	//ID ，serect_key , Bucket
	cfg := oss.LoadDefaultConfig().
		//WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()). //从环境变量读取
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider( //从配置文件读取
						config.Cfg.OssConfig.AccessKeyId,
						config.Cfg.OssConfig.AccessKeySecret, "")).
		WithRegion("cn-beijing") // 填写Bucket所在地域

	// 创建OSS客户端
	client := oss.NewClient(cfg)

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr("blog-ydx"),      // 存储空间名称
		Key:    oss.Ptr(path + filename), // 对象名称
		Body:   body,                     // 要上传的字符串内容
	}

	// 发送上传对象的请求
	_, err := client.PutObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put object %v", err)
		return ""
	}

	// 打印上传对象的结果
	url := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", "blog-ydx", "cn-beijing", path+filename)
	return url
}
