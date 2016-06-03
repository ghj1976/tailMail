package config

import (
	"bytes"
	"html/template"
	"log"
	"path"
	"time"

	"github.com/ghj1976/tailMail/stat"
	"github.com/ghj1976/tailMail/tail"
)

// 需要发送的内容信息
type TailInfoEntity struct {
	FileName    string    // 文件名
	MonitorTime time.Time // 监控的时间点

	IncrementalTxt *bytes.Buffer // 读出来的增量内容
	MailBodyHtml   *bytes.Buffer // 待发送的邮件body
	HasNewInfo     bool          //是不是有新的需要发送的信息

	LastFileSize int64 // 已经发送到那个位置了，上次文件的大小
	CurPos       int64 // 本次截取的文件位置， 如果本次内容太多， 不会截取全部。

	Config TailConfigEntity // 当前文件处理相关的配置信息
}

// 把需要发送的信息变成html格式的文本。
func (info *TailInfoEntity) getMailHtml(configDir string) (err error) {
	templateFileName := path.Join(configDir, "template.html")
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	t = template.Must(t, err)

	err = t.Execute(info.MailBodyHtml, info)
	if err != nil {
		log.Println(err)
		return err
	}
	return
}

// 监控一个文件，如果这个文件有内容新增，记录下来
func (tailInfo *TailInfoEntity) TailFile(configDir string) {

	// 获取文件新增的内容文本，用 buffer 是为了提高性能。
	var err error
	tailInfo.HasNewInfo, tailInfo.LastFileSize, err = tail.Tail(tailInfo.FileName, tailInfo.LastFileSize, tailInfo.IncrementalTxt)
	if err != nil {
		log.Println("err:", err)
	}

	if tailInfo.HasNewInfo {
		// 使用模板获取邮件正文内容
		err = tailInfo.getMailHtml(configDir)
		if err != nil {
			log.Println("err:", err)
		}

		// 清除已经没用的 buffer
		tailInfo.IncrementalTxt.Reset()

	} else {
		log.Println("这个文件没有发送变化！")
	}

}

func (info *TailInfoEntity) Stat(serverName string, si *stat.StatInfo) {
	// 需要统计时，记录发送统计信息。
	si.GetCurrDayStat(serverName)
	si.DayStat.Num++

	for _, mail := range info.Config.ToMailArr {
		mailStat, ok := si.DayStat.EmailStatMap[mail.Address]
		if !ok {
			mailStat = stat.StatisticsReportEMail{}
		}
		mailStat.EMailName = mail.Name
		mailStat.EMailAddress = mail.Address
		mailStat.Num++

		fileStat, ok := mailStat.FileStatMap[info.FileName]
		if !ok {
			fileStat = stat.StatisticsReportFile{}
		}
		fileStat.FileName = info.FileName
		fileStat.Num++
		fileStat.Subject = info.Config.Subject
		if mailStat.FileStatMap == nil {
			mailStat.FileStatMap = make(map[string]stat.StatisticsReportFile)
		}
		mailStat.FileStatMap[fileStat.FileName] = fileStat

		if si.DayStat.EmailStatMap == nil {
			si.DayStat.EmailStatMap = make(map[string]stat.StatisticsReportEMail)
		}
		si.DayStat.EmailStatMap[mail.Address] = mailStat
	}
	si.WriteDayStatValue()
}
