// 配置节相关类和方法
package config

import (
	"net/mail"
)

// 文件监控配置文件中每一个监控项的配置
type TailConfigEntity struct {
	FileName            string         // 文件名
	FileNameUseTemplate bool           // 文件名是否使用模版 郭红俊 20151222 增加新功能
	Subject             string         // 发送邮件的主题
	Remark              string         // 备注，用于区分发送的内容
	ToMailArr           []mail.Address // 需要发送的邮件地址

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

// 统计配置项
type StatConfig struct {
	Enable     bool   // 启动配置
	ServerName string // 服务器的唯一编号字符串，一般为服务器IP，或者名字。
}

// 文件监控配置文件
type TailConfigCollectionEntity struct {
	MailServer SmtpMailServerEntity // 邮件发送服务器的配置信息
	ConfigArr  []TailConfigEntity   // 多个文件进度集合
	Stat       StatConfig           // 统计配置
}

func (r *TailConfigCollectionEntity) GetAllEmail() []mail.Address {
	mailArr := []mail.Address{}
	for _, c := range r.ConfigArr {
		for _, m := range c.ToMailArr {
			appendMail(&mailArr, m)
		}
	}
	return mailArr
}

func appendMail(arr *[]mail.Address, ma mail.Address) {
	for _, m1 := range *arr {
		if m1.Address == ma.Address {
			return
		}
	}
	*arr = append(*arr, ma)
}
