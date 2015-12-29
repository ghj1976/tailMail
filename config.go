package tailMail

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	jsonConfigFileName   string // json 格式 配置文件
	jsonProgressFileName string // json 格式 进度文件
	tomlConfigFileName   string // toml 格式 配置文件
	tomlProgressFileName string // toml 格式 进度文件
	configType           string // 配置类型，只能是 json 或 toml
)

// 初始化， 配置及进度文件准备好
func InitConfigFile(configDir, config_type string) {
	if config_type == "json" {
		configType = "json"
		jsonConfigFileName = path.Join(configDir, "config.json")
		jsonProgressFileName = path.Join(configDir, "progress.json")
	} else {
		configType = "toml"
		tomlConfigFileName = path.Join(configDir, "config.toml")
		tomlProgressFileName = path.Join(configDir, "progress.toml")
	}
}

// 读配置文件
func ReadConfig() (configArr TailConfigCollectionEntity, err error) {
	if configType == "json" {
		txt, err := ioutil.ReadFile(jsonConfigFileName)
		if err != nil {
			log.Printf("File error: %v\n", err)
			os.Exit(1)
		}
		//fmt.Printf("%s\n", string(txt))
		err = json.Unmarshal(txt, &configArr)
	} else {
		_, err = toml.DecodeFile(tomlConfigFileName, &configArr)
	}
	return
}

// 写配置文件
func WriteConfig(configArr *TailConfigCollectionEntity) (err error) {
	if configType == "json" {
		txt, err := json.Marshal(configArr)
		if err != nil {
			log.Println("json err:", err)
			return err
		}
		err = ioutil.WriteFile(jsonConfigFileName, txt, 0644)
		if err != nil {
			log.Println("json err:", err)
		}
		return err
	} else {
		return writeTOMLFile(tomlConfigFileName, configArr)
	}

}

// 读进度文件
func ReadProgress() (progressArr TailProgressCollectionEntity, err error) {
	if configType == "json" {
		txt, err := ioutil.ReadFile(jsonProgressFileName)
		if err != nil && os.IsNotExist(err) {
			// 文件不存在，不影响使用
			progressArr = TailProgressCollectionEntity{
				ProgressMap: map[string]int64{},
			}
			return progressArr, nil
		}

		if err != nil {
			log.Printf("File error: %v\n", err)
			os.Exit(1)
		}
		//fmt.Printf("%s\n", string(txt))
		json.Unmarshal(txt, &progressArr)
	} else {

		_, err := ioutil.ReadFile(tomlProgressFileName)
		if err != nil && os.IsNotExist(err) {
			// 文件不存在，不影响使用
			progressArr = TailProgressCollectionEntity{
				ProgressMap: map[string]int64{},
			}
			return progressArr, nil
		}

		_, err = toml.DecodeFile(tomlProgressFileName, &progressArr)
	}

	// 如果文件读取出错，不能用 nil 指针。
	if progressArr.ProgressMap == nil {
		err = nil
		progressArr = TailProgressCollectionEntity{
			ProgressMap: map[string]int64{},
		}

	}
	return
}

// 更新写进度文件
func WriteProgress(progressArr *TailProgressCollectionEntity) (err error) {
	if configType == "json" {
		txt, err := json.Marshal(progressArr)
		if err != nil {
			log.Println("json err:", err)
		}
		err = ioutil.WriteFile(jsonProgressFileName, txt, 0644)
		if err != nil {
			log.Println("json err:", err)
		}
		return err
	} else {
		return writeTOMLFile(tomlProgressFileName, progressArr)

	}
}

// 把一个 对象 toml 序列化后写到指定文件
func writeTOMLFile(filename string, obj interface{}) error {
	fo, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fo.Close()

	// var firstBuffer bytes.Buffer
	e := toml.NewEncoder(fo)
	err = e.Encode(obj)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}

}
