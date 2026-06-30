package utils

import (
	"blog/config"
	"crypto/rand"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/big"
)

func SendEmailToQQ(to, subject string) error {
	msg := gomail.NewMessage()

	// 发件人
	msg.SetAddressHeader("From", config.Cfg.MailConfig.Username, config.Cfg.MailConfig.FromName)

	// 收件人
	msg.SetHeader("To", to)

	// 邮件标题
	msg.SetHeader("Subject", subject)

	// HTML 内容
	msg.SetBody("text/html", bodyByHTML(GenerateCode()))

	dialer := gomail.NewDialer(
		config.Cfg.MailConfig.Host,
		config.Cfg.MailConfig.Port,
		config.Cfg.MailConfig.Username,
		config.Cfg.MailConfig.Password,
	)

	return dialer.DialAndSend(msg)
}

// 生成6位验证码
func GenerateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

func bodyByHTML(verificationCode string) string {
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>邮箱验证码</title>
</head>
<body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,'Microsoft YaHei',sans-serif;">
    <table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0"
                       style="background:#ffffff;border-radius:10px;padding:40px;box-shadow:0 2px 8px rgba(0,0,0,.08);">

                    <tr>
                        <td style="font-size:26px;font-weight:bold;color:#333;">
                            YDX Blog
                        </td>
                    </tr>

                    <tr>
                        <td style="padding-top:24px;font-size:16px;color:#333;">
                            您好，
                        </td>
                    </tr>

                    <tr>
                        <td style="padding-top:12px;font-size:15px;color:#666;line-height:28px;">
                            您正在进行邮箱身份验证，请使用下面的验证码完成操作：
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding:30px 0;">
                            <div style="
                                display:inline-block;
                                padding:16px 36px;
                                background:#1677ff;
                                color:#fff;
                                font-size:32px;
                                font-weight:bold;
                                letter-spacing:8px;
                                border-radius:8px;">
                                %s
                            </div>
                        </td>
                    </tr>

                    <tr>
                        <td style="font-size:14px;color:#666;line-height:26px;">
                            • 验证码有效期：<strong>5 分钟</strong><br>
                            • 请勿将验证码泄露给任何人。<br>
                            • 如果这不是您的操作，请忽略本邮件。
                        </td>
                    </tr>

                    <tr>
                        <td style="padding-top:40px;border-top:1px solid #eee;font-size:13px;color:#999;">
                            此邮件由系统自动发送，请勿直接回复。<br>
                            © 2026 YDX Blog. All Rights Reserved.
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, verificationCode)
	return body
}
