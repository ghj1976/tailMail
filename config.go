package tailMail

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var configFileName string
var progressFileName string

// 初始化， 配置及进度文件准备好
func InitConfigFile(configDir string) {
	configFileName = path.Join(configDir, "config.json")
	progressFileName = path.Join(configDir, "progress.json")
}

// 读配置文件
func ReadConfig() (configArr TailConfigCollectionEntity, err error) {
	txt, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(txt))
	json.Unmarshal(txt, &configArr)
	return
}

// 写配置文件
func WriteConfig(configArr *TailConfigCollectionEntity) (err error) {
	txt, err := json.Marshal(configArr)
	if err != nil {
		log.Println("json err:", err)
	}
	err = ioutil.WriteFile(configFileName, txt, 0644)
	if err != nil {
		log.Println("json err:", err)
	}
	return
}

// 读进度文件
func ReadProgress() (progressArr TailProgressCollectionEntity, err error) {
	txt, err := ioutil.ReadFile(progressFileName)
	if err != nil && os.IsNotExist(err) {
		// 文件不存在，不影响使用
		err = nil
		progressArr = TailProgressCollectionEntity{
			ProgressMap: map[string]int64{},
		}
		return
	}
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(txt))
	json.Unmarshal(txt, &progressArr)
	return
}

// 写进度文件
func WriteProgress(progressArr *TailProgressCollectionEntity) (err error) {
	txt, err := json.Marshal(progressArr)
	if err != nil {
		log.Println("json err:", err)
	}
	err = ioutil.WriteFile(progressFileName, txt, 0644)
	if err != nil {
		log.Println("json err:", err)
	}
	return
}
