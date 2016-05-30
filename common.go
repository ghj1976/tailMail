package tailMail

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// 把一个 对象 toml 序列化后写到指定文件
func WriteTOMLFile(filename string, obj interface{}) error {
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
