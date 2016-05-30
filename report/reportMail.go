package report

import (
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/ghj1976/tailMail/config"
	"github.com/ghj1976/tailMail/email"
	"github.com/ghj1976/tailMail/stat"
)

// 发送昨日报告邮件
func SendReportMail(configDir, serverName string, mailServer config.SmtpMailServerEntity, toMailArr []mail.Address) {
	b, report := stat.ReadYesterdayStatValue(configDir, serverName)
	if !b {
		log.Println("没有昨天的报表数据！")
		return
	}

	err, body := report.GetMailRportHtml(configDir)
	if err != nil {
		log.Println(err)
		return
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	subject := fmt.Sprintf("监控日报-%s-%d.%d.%d", serverName,
		yesterday.Year(), yesterday.Month(), yesterday.Day())

	email.SendSSLMail(mailServer, subject, body, "", toMailArr)

}
