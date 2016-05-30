package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"

	"github.com/ghj1976/tailMail/config"
)

// 发送html格式的邮件
func SendHtmlMail(mailServer config.SmtpMailServerEntity, subject, body string, toMailArr []string) {

	// 首先验证邮件登陆信息
	var auth smtp.Auth
	if mailServer.NeedLogin {
		auth = smtp.PlainAuth(mailServer.ServerAddress, mailServer.LoginUser, mailServer.LoginPassword, mailServer.ServerAddress)
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("From: \"%s\" <%s>\r\n", mailServer.SendMailUserMail, mailServer.SendMailUserMail))
	buffer.WriteString("To: ")
	// 注意，这里分隔符是 ,
	for _, mail := range toMailArr {
		buffer.WriteString(fmt.Sprintf("\"%s\"<%s>,", mail, mail))
	}
	buffer.WriteString("\r\n")
	buffer.WriteString(fmt.Sprintf("Subject: %s \r\n", subject))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type:text/html; charset=UTF-8\r\n")
	buffer.WriteString("Content-Transfer-Encoding:8bit\r\n")
	buffer.WriteString(body)
	buffer.WriteString("\r\n")

	serverAddress := fmt.Sprintf("%s:%d", mailServer.ServerAddress, mailServer.ServerAddressPort)
	fmt.Println(serverAddress)

	var err error
	// 发送邮件
	if mailServer.NeedLogin {
		err = smtp.SendMail(serverAddress, auth, mailServer.SendMailUserMail.String(), toMailArr, buffer.Bytes())

	} else {
		err = smtp.SendMail(serverAddress, nil, mailServer.SendMailUserMail.String(), toMailArr, buffer.Bytes())
	}

	if err != nil {
		log.Println("发送邮件错误：", err)
	} else {
		log.Println("send mail finish!")
	}
}
