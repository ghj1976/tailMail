package config

import (
	"fmt"
	"net/mail"
)

// 初始化配置文件
func InitConfigFile() {

	// 需要监控的文件
	fileName1 := "/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/test/11.log"
	fileName2 := "/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/test/22_{{formatNow \"2006-01-02\"}}.log"

	// 写配置信息
	ci := ConfigInfo{
		configFileType: "toml",
		configFileName: "config.toml",
		Config: TailConfigCollectionEntity{
			MailServer: SmtpMailServerEntity{
				ServerAddress:     "smtp.exmail.qq.com",
				ServerAddressPort: 465,
				NeedLogin:         true,
				LoginUser:         "guohongjun@bbb.com",
				LoginPassword:     "*******",
				SendMailUserMail:  mail.Address{Name: "郭红俊", Address: "guohongjun@bbb.com"},
			},
			Stat: StatConfig{
				Enable:     true,
				ServerName: "10.162.222.210",
			},
			ConfigArr: []TailConfigEntity{
				{
					FileName:            fileName1,
					FileNameUseTemplate: false,
					Subject:             "异常监控报告，服务器：61.235",
					Remark:              "",
					ToMailArr: []mail.Address{{Name: "ghj1976", Address: "ghj1976@aaa.com"},
						{Name: "郭红俊", Address: "guohongjun@bbb.com"}}},
				{
					FileName:            fileName2,
					FileNameUseTemplate: true,
					Subject:             "测试邮件标题",
					Remark:              "",
					ToMailArr:           []mail.Address{{Name: "ghj1976", Address: "ghj1976@aaa.com"}}}}},
	}
	err := ci.WriteConfig()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化配置文件完成！")
	}

}
