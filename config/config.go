// 这里是对配置相关的封装。
package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ghj1976/tailMail"
	"github.com/ghj1976/tailMail/progress"
)

type ConfigInfo struct {
	configFileType string                      // 配置文件类型
	configFileName string                      // 配置文件名
	Config         *TailConfigCollectionEntity // 配置信息
}

// 构造一个 配置类
func NewConfigInfo(ft, dir string) *ConfigInfo {

	var fn string
	if ft == "json" {
		fn = path.Join(dir, "config.json")
	} else {
		ft = "toml" // 默认 toml 格式的配置文件
		fn = path.Join(dir, "config.toml")
	}

	ci := &ConfigInfo{
		configFileType: ft,
		configFileName: fn,
	}
	return ci
}

// 读配置文件
func (ci *ConfigInfo) ReadConfig() error {
	log.Println("config fileName:", ci.configFileName)
	if ci.configFileType == "json" {
		txt, err := ioutil.ReadFile(ci.configFileName)
		if err != nil {
			log.Printf("File error: %v\n", err)
			os.Exit(1)
		}
		return json.Unmarshal(txt, ci.Config)

	} else {
		if ci.Config == nil {
			ci.Config = &TailConfigCollectionEntity{}
		}
		_, err := toml.DecodeFile(ci.configFileName, ci.Config)
		return err
	}
	return nil
}

// 写配置文件
func (ci *ConfigInfo) WriteConfig() error {
	if ci.configFileType == "json" {
		txt, err := json.Marshal(ci.Config)
		if err != nil {
			log.Println("json err:", err)
			return err
		}
		err = ioutil.WriteFile(ci.configFileName, txt, 0644)
		if err != nil {
			log.Println("json err:", err)
		}
		return err
	} else {
		return tailMail.WriteTOMLFile(ci.configFileName, ci.Config)
	}

}

// 装配、准备待处理的数据Map
func (ci *ConfigInfo) PrepareWork(pi *progress.ProcessInfo) map[string]TailInfoEntity {
	// 组装需要发送的信息数组
	var tailInfoMap map[string]TailInfoEntity
	tailInfoMap = make(map[string]TailInfoEntity, len(ci.Config.ConfigArr))

	for _, conf := range ci.Config.ConfigArr {
		fn := conf.FileName
		if conf.FileNameUseTemplate { // 文件名使用模版机制。
			fn = tailMail.FormatFileName(fn) // 模版替换
		}
		tailInfo := TailInfoEntity{
			FileName:       fn,
			MonitorTime:    time.Now(),
			IncrementalTxt: new(bytes.Buffer), // 读取出来的增量内容
			MailBodyHtml:   new(bytes.Buffer), // 要发送的邮件内容正文body
			HasNewInfo:     true,              //是不是有新的需要发送的信息
			LastFileSize:   pi.GetLastFileSize(fn),
			Config:         conf,
		}
		tailInfoMap[fn] = tailInfo
	}
	return tailInfoMap
}
