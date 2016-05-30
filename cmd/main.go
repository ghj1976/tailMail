package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ghj1976/tailMail/config"
	"github.com/ghj1976/tailMail/email"
	"github.com/ghj1976/tailMail/logger"
	"github.com/ghj1976/tailMail/progress"
	"github.com/ghj1976/tailMail/report"
)

var (
	outPutLog    = flag.Bool("o", false, "是否把提示输出到log文件？默认是 执行目录下log目录下，每天一个文件。")
	configPath   = flag.String("p", "", "配置文件、进度文件所在目录。默认是当前目录。")
	initConfig   = flag.Bool("i", false, "初始化配置文件")
	configType   = flag.String("ct", "toml", "配置文件的类型，默认 toml 格式， 也支持 json 格式")
	reportConfig = flag.Bool("r", false, "发送每日监控报告，如果发现这个参数被设置为true后，在被调用时，自动发送昨天的统计报告邮件。")

	ci *config.ConfigInfo // 当前的配置信息
)

func main() {
	log.Println("##begin##")
	flag.Parse() // 读取命令参数

	// 执行目录路径
	configDir := strings.TrimSpace(*configPath)
	if configDir == "" || configDir == "." {
		configDir = "./" // 当前目录
	} else {
		configDir = path.Clean(configDir) // 选定的目录
	}
	log.Println("configDir:", configDir)

	// 准备日志输出。
	logger.InitLogFile(*outPutLog, configDir)
	defer logger.LoggerFinish()

	// 如果需要重建配置文件
	if *initConfig {
		config.InitConfigFile()
		os.Exit(-1)
		return
	}

	// 加载、读取配置文件
	ci = config.NewConfigInfo(*configType, configDir)
	err := ci.ReadConfig()
	if err != nil {
		log.Println("读取配置文件错误，", err)
	}

	// 如果是要发送昨日报告邮件
	if *reportConfig {
		if ci.Config.Stat.Enable {
			// 发送昨日统计报告邮件
			report.SendReportMail(configDir, ci.Config.Stat.ServerName,
				ci.Config.MailServer, ci.Config.GetAllEmail())
		}
		os.Exit(-1)
		return
	}

	// 开始循环分析文件
	work(configDir)

	// 打印结束时间
	log.Println("==end==")
}

func work(configDir string) {

	// 读取进度文件
	pi := progress.NewProcessInfo(configDir)
	err := pi.ReadProgress()
	if err != nil {
		log.Println(err)
	}
	log.Println(".... 读取进度文件完成 ....")

	// 准备要 tailMail 处理的文件集合
	tailInfoMap := ci.PrepareWork(pi)

	// 遍历每个配置文件，并处理数据
	for _, info := range tailInfoMap {
		log.Println("正在分析文件：", info.FileName)
		// 检查文件，获得是否有更新等信息。
		info.TailFile(configDir)

		if info.HasNewInfo {
			// 发送邮件
			log.Println("开始发送邮件！", info.Config.ToMailArr)
			email.SendSSLMail(ci.Config.MailServer, info.Config.Subject,
				info.MailBodyHtml.String(), "", info.Config.ToMailArr)
			log.Println("..邮件发送完成！...")

			// 更新进度信息
			log.Println("更新进度信息配置文件。")
			pi.UpdateFile(info.FileName, info.LastFileSize)

			// 更新统计信息

		}
	}
}
