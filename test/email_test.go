package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <1151695676@qq.com>"
	e.To = []string{"1151695676@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("<h1>您的验证码：<b>123456</b><h1>")
	//xujzslhdysowhgia
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1151695676@qq.com", "xujzslhdysowhgia", "smtp.qq.com"))
	//err := e.SendWithTLS("smtp.qq.com:587",
	//	smtp.PlainAuth("", "2162768982@qq.com", "wqy610366", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}

}
