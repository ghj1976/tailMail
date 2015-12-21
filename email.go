// 郭红俊 20140217 更新版本。
// 比之前版本性能更佳
package tailMail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
)

const (

	// 邮件中的分割符
	marker = "ACUSTOMUNIQUGUOHONGJUNEBOUNDARY"
)

// 发送html格式的邮件
func SendHtmlMail(mailServer SmtpMailServerEntity, subject, body string, toMailArr []string) {

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

/*
给指定的用户发送邮件
attachmentFilePath 为 nil 就是不发送
这里没有提供登陆验证功能

技术参考：
http://stackoverflow.com/questions/4018709/how-to-create-an-email-with-embedded-images-that-is-compatible-with-the-most-mai
*/
func SendMail(mailServer, fromMail, subject, body, attachmentFilePath string, toMailArr []string) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("From: \"%s\"<%s>\r\n", fromMail, fromMail))
	buffer.WriteString("To: ")
	// 注意，这里分隔符是 ,
	for _, mail := range toMailArr {
		buffer.WriteString(fmt.Sprintf("\"%s\"<%s>,", mail, mail))
	}
	buffer.WriteString("\r\n")
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", marker))
	buffer.WriteString(fmt.Sprintf("--%s\r\n", marker))

	buffer.WriteString("\r\n")
	buffer.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buffer.WriteString("Content-Transfer-Encoding:8bit\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(body)
	buffer.WriteString(fmt.Sprintf("--%s\r\n", marker))

	// 有附件的情况
	if len(attachmentFilePath) > 0 {

		buffer.WriteString("\r\n")
		buffer.WriteString(fmt.Sprintf("Content-Type: application/image; name=\"%s\"\r\n", attachmentFilePath))
		buffer.WriteString("Content-Transfer-Encoding:base64\r\n")
		buffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachmentFilePath))
		buffer.WriteString("Content-ID:im001\r\n")
		buffer.WriteString("\r\n")

		// 准备内嵌资源图片文件
		content, _ := ioutil.ReadFile(".\\tmp\\" + attachmentFilePath)
		encoded := base64.StdEncoding.EncodeToString(content)
		lineMaxLength := 500 //split the encoded file in lines of some max length (1000? 1024? I read 1024 somewhere, but hit a max of 1000 once, so I aim lower just in case)
		nbrLines := len(encoded) / lineMaxLength
		for i := 0; i < nbrLines; i++ {
			buffer.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n") //\n converted to \r\n by smtp pacakge
		}
		buffer.WriteString(encoded[nbrLines*lineMaxLength:])

		buffer.WriteString(fmt.Sprintf("--%s--\r\n", marker))
	}

	// fmt.Println("0000")

	// 发送邮件
	err := smtp.SendMail(mailServer, nil, fromMail, toMailArr, buffer.Bytes())
	if err != nil {
		log.Println(err)
	}
}

// 发送 ssl 认证的邮件
// 相关代码参考： https://gist.github.com/chrisgillis/10888032
func SendSSLMail(mailServer SmtpMailServerEntity, subject, body, attachmentFilePath string, toMailArr []mail.Address) {

	var hasAttach = false
	if len(attachmentFilePath) > 0 {
		hasAttach = true
	}
	// 准备要发送的内容（HTML格式）
	var buffer bytes.Buffer
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buffer.WriteString("Content-Transfer-Encoding:8bit\r\n")

	buffer.WriteString(fmt.Sprintf("From: \"%s\"<%s>\r\n", encodeMailString(mailServer.SendMailUserMail.Name), mailServer.SendMailUserMail.Address))
	buffer.WriteString("To: ")
	// 注意，这里分隔符是 ,
	for _, mail := range toMailArr {
		buffer.WriteString(fmt.Sprintf("\"%s\"<%s> , ", encodeMailString(mail.Name), mail.Address))
	}
	buffer.WriteString("\r\n")

	buffer.WriteString(fmt.Sprintf("Subject: %s \r\n", encodeMailString(subject)))

	if hasAttach {
		buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", marker))
		buffer.WriteString(fmt.Sprintf("--%s\r\n", marker))
	}

	buffer.WriteString("\r\n")
	buffer.WriteString(body)

	// 有附件的情况
	if hasAttach {
		buffer.WriteString(fmt.Sprintf("--%s\r\n", marker))

		buffer.WriteString("\r\n")
		buffer.WriteString(fmt.Sprintf("Content-Type: application/image; name=\"%s\"\r\n", attachmentFilePath))
		buffer.WriteString("Content-Transfer-Encoding:base64\r\n")
		buffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachmentFilePath))
		buffer.WriteString("Content-ID:im001\r\n")
		buffer.WriteString("\r\n")

		// 准备内嵌资源图片文件
		content, _ := ioutil.ReadFile(".\\tmp\\" + attachmentFilePath)
		encoded := base64.StdEncoding.EncodeToString(content)
		lineMaxLength := 500 //split the encoded file in lines of some max length (1000? 1024? I read 1024 somewhere, but hit a max of 1000 once, so I aim lower just in case)
		nbrLines := len(encoded) / lineMaxLength
		for i := 0; i < nbrLines; i++ {
			buffer.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n") //\n converted to \r\n by smtp pacakge
		}
		buffer.WriteString(encoded[nbrLines*lineMaxLength:])

		buffer.WriteString(fmt.Sprintf("--%s--\r\n", marker))
	} else {
		buffer.WriteString("\r\n")
	}

	log.Println("组装要发送的内容。 ok")

	// Connect to the SMTP Server

	auth := smtp.PlainAuth(mailServer.ServerAddress, mailServer.LoginUser, mailServer.LoginPassword, mailServer.ServerAddress)
	log.Println("...smtp.PlainAuth... ok")

	servername := fmt.Sprintf("%s:%d", mailServer.ServerAddress, mailServer.ServerAddressPort)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mailServer.ServerAddress,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}
	log.Println("tls.Dial ok")

	c, err := smtp.NewClient(conn, mailServer.ServerAddress)
	if err != nil {
		log.Panic(err)
	}
	log.Println("smtp.NewClient ok")

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}
	log.Println("c.Auth ok")

	// To && From
	if err = c.Mail(mailServer.LoginUser); err != nil {
		log.Panic(err)
	}
	log.Println("c.Mail(mailServer.LoginUser) ok")

	for _, to := range toMailArr {
		if err = c.Rcpt(to.Address); err != nil {
			log.Panic(err)
		}
	}
	log.Println("c.Rcpt ok")

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}
	log.Println("c.Data ok")

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		log.Panic(err)
	}
	log.Println("w.Write ok")

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
	log.Println("w.Close ok")

	c.Quit()
}

// http://www.mamicode.com/info-detail-893793.html
// 标题乱码解决技术
func encodeMailString(str string) string {
	return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(str)))
}
