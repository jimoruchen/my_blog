package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"web_app/settings"
)

// SendVerificationCode 使用原生 net/smtp 实现 (Port 465/SSL)
func SendVerificationCode(toEmail, code string) error {
	cfg := settings.Conf.EmailConfig

	// 1. 准备邮件内容
	// 注意：QQ 邮箱要求 Header 中的 From 必须与认证账号完全一致
	// Subject 和 Body 之间必须有一个空行
	header := make(map[string]string)
	header["From"] = cfg.Username
	header["To"] = toEmail
	header["Subject"] = "【my_blog】邮箱验证码"
	header["Content-Type"] = "text/plain; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	body := fmt.Sprintf("\r\n您好！\n\n您正在注册 my_blog 账号，验证码为：%s\n\n该验证码 5 分钟内有效，请勿泄露给他人。\n\n(如非本人操作，请忽略)", code)
	message += body

	// 2. 建立 TLS 连接 (直接连接 465 端口)
	// 这里使用了和测试脚本一样的 InsecureSkipVerify: true 以防止本地证书报错
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cfg.Host,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接邮件服务器失败: %w", err)
	}
	defer conn.Close()

	// 3. 创建 SMTP 客户端
	c, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer c.Quit()

	// 4. 认证 (Auth)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败 (请检查密码/授权码): %w", err)
	}

	// 5. 发送邮件 (Mail -> Rcpt -> Data)
	if err = c.Mail(cfg.Username); err != nil {
		return fmt.Errorf("发送发件人信息失败: %w", err)
	}
	if err = c.Rcpt(toEmail); err != nil {
		return fmt.Errorf("发送收件人信息失败: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入流失败: %w", err)
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("关闭写入流失败: %w", err)
	}

	return nil
}
