package tailMail

import (
	"bytes"
	"html/template"
	"log"
	"net/mail"
	"path"
	"time"
)

// 文件监控配置文件中每一个监控项的配置
type TailConfigEntity struct {
	FileName  string         // 文件名
	Subject   string         // 发送邮件的主题
	Remark    string         // 备注，用于区分发送的内容
	ToMailArr []mail.Address // 需要发送的邮件地址

}

// 需要发送的内容信息
type TailInfoEntity struct {
	FileName    string    // 文件名
	MonitorTime time.Time // 监控的时间点

	IncrementalTxt *bytes.Buffer // 读出来的增量内容
	MailBodyHtml   *bytes.Buffer // 待发送的邮件body
	HasNewInfo     bool          //是不是有新的需要发送的信息

	Config TailConfigEntity // 原始配置信息

	LastFileSize int64 // 已经发送到那个位置了，上次文件的大小

	CurPos int64 // 本次截取的文件位置， 如果本次内容太多， 不会截取全部。

	MailServer SmtpMailServerEntity // 邮件服务器地址

}

// 邮件发送服务器的配置信息
type SmtpMailServerEntity struct {
	ServerAddress     string       // 邮件服务器地址，不包含端口号，比如 ： 101.11.154.4， stmp.163.com 这样的字符串
	ServerAddressPort int          // 邮件服务器端口号，默认 25
	NeedLogin         bool         // 邮箱需要登陆么？ true 需要登陆， false 不用
	LoginUser         string       // 邮箱登陆名
	LoginPassword     string       // 邮箱登陆密码
	SendMailUserMail  mail.Address // 发送邮件者邮箱，回复邮件时需要回复到这个地址

}

// 文件监控配置文件
type TailConfigCollectionEntity struct {
	MailServer SmtpMailServerEntity // 邮件发送服务器的配置信息
	ConfigArr  []TailConfigEntity   // 多个文件进度集合
}

// 文件监控进度配置文件
type TailProgressCollectionEntity struct {
	ProgressMap map[string]int64 // 多个文件进度集合
}

// 把需要发送的信息变成html格式的文本。
func (info *TailInfoEntity) GetMailHtml(configDir string) (err error) {
	templateFileName := path.Join(configDir, "template.html")
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	t = template.Must(t, err)

	//fmt.Println(t.Name())

	//t, _ := tetemplatemplate.New("name").Parse("src {{.}} ee")

	//var buffer bytes.Buffer
	//buffer := new(bytes.Buffer)

	//fmt.Println(info)
	//t.ExecuteTemplate(&buffer, "template.html", info)

	err = t.Execute(info.MailBodyHtml, info)
	if err != nil {
		log.Println(err)
		return err
	}
	return
}
