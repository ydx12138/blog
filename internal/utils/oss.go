package utils

import (
	"blog/config"
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func UploadToOss(body multipart.File, path string, filename string) (string, error) {
	ossConfig := config.Cfg.OssConfig
	if ossConfig.AccessKeyId == "" || ossConfig.AccessKeySecret == "" || ossConfig.Bucket == "" || ossConfig.Endpoint == "" {
		return "", errors.New("oss config is incomplete")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			ossConfig.AccessKeyId,
			ossConfig.AccessKeySecret,
			"",
		)).
		WithRegion(ossConfig.Endpoint)
	client := oss.NewClient(cfg)

	key := path + filename
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(ossConfig.Bucket),
		Key:    oss.Ptr(key),
		Body:   body,
	}
	if _, err := client.PutObject(context.TODO(), request); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", ossConfig.Bucket, ossConfig.Endpoint, key), nil
}
