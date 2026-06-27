package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/google/uuid"
)

/*
Go SDK V2 客户端初始化配置说明：

1. 签名版本：Go SDK V2 默认使用 V4 签名，提供更高的安全性
2. Region配置：初始化 Client 时，您需要指定阿里云通用 Region ID 作为发起请求地域的标识
3. Endpoint配置：
   - 可以通过 Endpoint 参数，自定义服务请求的访问域名
   - 当不指定时，SDK 默认根据 Region 信息，构造公网访问域名
4. 协议配置：
   - SDK 构造访问域名时默认采用 HTTPS 协议
   - 如需采用 HTTP 协议，请在指定域名时指定为 HTTP
*/

func UploadToOss() {
	// 方式一：只填写Region（推荐）
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion("cn-beijing") // 填写Bucket所在地域

	// 方式二：同时填写Region和Endpoint
	// cfg := oss.LoadDefaultConfig().
	//     WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
	//     WithRegion("<region-id>").                                // 填写Bucket所在地域
	//     WithEndpoint("<endpoint>")     // 填写Bucket所在地域对应的公网Endpoint

	// 创建OSS客户端
	client := oss.NewClient(cfg)

	// 定义要上传的字符串内容
	body := strings.NewReader("hi oss")
	//路径
	yourObjectKey := fmt.Sprintf(
		"blog/%s/%s.jpg",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
	)

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr("blog-ydx"),    // 存储空间名称
		Key:    oss.Ptr(yourObjectKey), // 对象名称
		Body:   body,                   // 要上传的字符串内容
	}

	// 发送上传对象的请求
	result, err := client.PutObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put object %v", err)
	}

	// 打印上传对象的结果
	log.Printf("Status: %#v\n", result.Status)
	log.Printf("RequestId: %#v\n", result.ResultCommon.Headers.Get("X-Oss-Request-Id"))
	log.Printf("ETag: %#v\n", *result.ETag)
}
