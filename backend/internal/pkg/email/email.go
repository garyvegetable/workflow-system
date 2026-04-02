package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

type EmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

// Configure 更新邮件服务配置
func (s *EmailService) Configure(host string, port int, user, password, from string) {
	s.host = host
	s.port = port
	s.user = user
	s.password = password
	s.from = from
}

// IsConfigured 检查是否已配置
func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.user != "" && s.password != ""
}

// Send 发送邮件
func (s *EmailService) Send(to, subject, body string) error {
	if !s.IsConfigured() {
		log.Printf("[Email Mock] Email: %s, Title: %s\nContent: %s\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=utf-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg.Bytes())
	if err != nil {
		log.Printf("[Email Error] Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("[Email] Sent to %s: %s", to, subject)
	return nil
}

// SendEmployeeCredentials 发送员工账号凭证
func (s *EmailService) SendEmployeeCredentials(toEmail, username, password string) error {
	if toEmail == "" {
		log.Println("[Email] No email address provided, skipping")
		return nil
	}

	subject := "您的账号已创建"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
<h2>账号信息</h2>
<p>您的账号已创建成功，以下是登录信息：</p>
<table style="border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>用户名</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>密码</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
</table>
<p>请登录后尽快修改密码。</p>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`, username, password)

	return s.Send(toEmail, subject, body)
}

// SendApprovalNotification 发送审批通知
func (s *EmailService) SendApprovalNotification(approverEmail, applicantName, workflowName, title string) error {
	if approverEmail == "" {
		return nil
	}

	subject := fmt.Sprintf("【待审批】%s - 来自 %s", title, applicantName)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
<h2>您有新的审批任务</h2>
<p>您有一个新的待审批任务：</p>
<table style="border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>流程名称</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>申请人</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>标题</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
</table>
<p>请登录系统进行审批。</p>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`, workflowName, applicantName, title)

	return s.Send(approverEmail, subject, body)
}

// SendResultNotification 发送审批结果通知
func (s *EmailService) SendResultNotification(applicantEmail, workflowName, title, result string) error {
	if applicantEmail == "" {
		return nil
	}

	subject := fmt.Sprintf("【审批结果】%s - %s", title, result)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
<h2>您的申请已处理</h2>
<p>您的申请已审批完成：</p>
<table style="border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>流程名称</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>标题</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>结果</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
</table>
<p>请登录系统查看详情。</p>
<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
<p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`, workflowName, title, result)

	return s.Send(applicantEmail, subject, body)
}

// ParseTemplate 解析邮件模板
func ParseTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
