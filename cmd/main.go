package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ghj1976/tailMail"
	"log"
	"net/mail"
	"os"
	"path"
	"strings"
	"time"
)

var (
	outPutLog  = flag.Bool("o", false, "是否把提示输出到log文件？默认是 执行目录下log目录下，每天一个文件。")
	configPath = flag.String("p", "", "配置文件、进度文件所在目录。默认是当前目录。")
	initConfig = flag.Bool("i", false, "初始化配置文件")
)

func main() {
	// 打印当前时间
	log.Println("##begin##")
	// 读取命令参数
	flag.Parse()
	//fmt.Println(flag.Args())

	tailMail.OutPutLog = *outPutLog

	// 执行目录路径
	configDir := strings.TrimSpace(*configPath)
	if configDir == "" || configDir == "." {
		// 当前目录
		configDir = "./"
	} else {
		// 选定的目录
		configDir = path.Clean(configDir)
	}

	// 准备日志输出。
	tailMail.InitLogFile(configDir)
	defer tailMail.LoggerFinish()

	// 加载配置文件
	tailMail.InitConfigFile(configDir, "toml")

	// 命令行参数时重建配置文件
	if *initConfig {
		initConfigFile()
		os.Exit(-1)
		return
		// saveProgressInfoInit()
	}

	// 开始循环分析文件
	work(configDir)

	// 打印结束时间
	log.Println("==end==")
}

func work(configDir string) {

	log.Println("")
	log.Println(".... 开始 ....")

	log.Println("读取配置信息中...")
	// 读取配置文件
	configArr, err := tailMail.ReadConfig()
	if err != nil {
		log.Println(err)
	}

	// 读取进度文件
	progressMap, err := tailMail.ReadProgress()
	if err != nil {
		log.Println(err)
	}

	// 组装需要发送的信息数组
	var tailInfoMap map[string]tailMail.TailInfoEntity
	tailInfoMap = make(map[string]tailMail.TailInfoEntity, len(configArr.ConfigArr))

	for _, conf := range configArr.ConfigArr {

		fn := conf.FileName
		// 文件名使用模版机制。
		if conf.FileNameUseTemplate {
			// 模版替换
			fn = tailMail.FormatFileName(fn)
		}
		tailInfo := tailMail.TailInfoEntity{
			FileName:       fn,
			MonitorTime:    time.Now(),
			IncrementalTxt: new(bytes.Buffer), // 读取出来的增量内容
			MailBodyHtml:   new(bytes.Buffer), // 要发送的邮件内容正文body
			HasNewInfo:     true,              //是不是有新的需要发送的信息
			Config:         conf,
			LastFileSize:   progressMap.ProgressMap[fn],
			MailServer:     configArr.MailServer,
		}
		tailInfoMap[fn] = tailInfo
	}

	// 遍历每个配置文件，并处理数据
	for _, info := range tailInfoMap {
		log.Println("正在分析文件：", info.FileName)

		tailFileMail(&info, configDir)

		if info.HasNewInfo {
			log.Println("更新进度信息配置文件。")
			// 更新进度信息
			progressMap.ProgressMap[info.FileName] = info.LastFileSize

			err = tailMail.WriteProgress(&progressMap)
			if err != nil {
				log.Println("写进度配置文件异常：")
				log.Println(err)
			}
		}
	}
}

// 监控一个文件，如果这个文件有内容新增，根据配置信息，发送邮件给相关人
func tailFileMail(tailInfo *tailMail.TailInfoEntity, configDir string) {

	// 获取文件新增的内容文本，用 buffer 是为了提高性能。
	var err error
	tailInfo.HasNewInfo, tailInfo.LastFileSize, err = tailMail.Tail(tailInfo.FileName, tailInfo.IncrementalTxt, tailInfo.LastFileSize)
	if err != nil {
		log.Println("err:", err)
	}

	if tailInfo.HasNewInfo {
		// 使用模板获取邮件正文内容
		err = tailInfo.GetMailHtml(configDir)
		if err != nil {
			log.Println("err:", err)
		}

		// 清除已经没用的 buffer
		tailInfo.IncrementalTxt.Reset()

		log.Println("开始发送邮件！", tailInfo.Config.ToMailArr)

		// 发送邮件
		//tailMail.SendHtmlMail(tailInfo.MailServer, tailInfo.Config.Subject, tailInfo.MailBodyHtml.String(), tailInfo.Config.ToMailArr)

		tailMail.SendSSLMail(tailInfo.MailServer, tailInfo.Config.Subject, tailInfo.MailBodyHtml.String(), "", tailInfo.Config.ToMailArr)
		log.Println("..邮件发送完成！...")

	} else {
		log.Println("这个文件没有发送变化！")
	}

}

// 初始化配置文件
func initConfigFile() {

	// 需要监控的文件
	fileName1 := "/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/test/11.log"
	fileName2 := "/Users/ghj1976/project/mygocode/src/github.com/ghj1976/tailMail/test/22_{{formatNow \"2006-01-02\"}}.log"

	// 写配置信息
	configArr := tailMail.TailConfigCollectionEntity{
		MailServer: tailMail.SmtpMailServerEntity{
			ServerAddress:     "smtp.exmail.qq.com",
			ServerAddressPort: 465,
			NeedLogin:         true,
			LoginUser:         "guohongjun@bbb.com",
			LoginPassword:     "*******",
			SendMailUserMail:  mail.Address{Name: "郭红俊", Address: "guohongjun@bbb.com"},
		},
		ConfigArr: []tailMail.TailConfigEntity{
			{
				FileName:            fileName1,
				FileNameUseTemplate: false,
				Subject:             "异常监控报告，服务器：61.235",
				Remark:              "",
				ToMailArr:           []mail.Address{{Name: "ghj1976", Address: "ghj1976@aaa.com"}, {Name: "郭红俊", Address: "guohongjun@bbb.com"}}},
			{
				FileName:            fileName2,
				FileNameUseTemplate: true,
				Subject:             "测试邮件标题",
				Remark:              "",
				ToMailArr:           []mail.Address{{Name: "ghj1976", Address: "ghj1976@aaa.com"}}}}}

	err := tailMail.WriteConfig(&configArr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("初始化配置文件完成！")
	}

}

// 进度文件初始化测试，目前没有意义了。
func saveProgressInfoInit() {

	progressArr := tailMail.TailProgressCollectionEntity{
		ProgressMap: map[string]int64{
			"E:\\tmp\\dnsdata\\auth.acc-2014-01-18-64.35.log": int64(0),
			"E:\\tmp\\dnsdata\\138.18.stdout-2014_01_22.log":  int64(0),
		},
	}

	err := tailMail.WriteProgress(&progressArr)
	if err != nil {
		fmt.Println(err)
	}
}
