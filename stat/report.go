package stat

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/BurntSushi/toml"
)

// 读昨天统计结果数据，以便发送统计邮件。
func ReadYesterdayStatValue(dir, serverName string) (b bool, report *StatisticsReportDaily) {
	yesterday := time.Now().AddDate(0, 0, -1)
	yfn := path.Join(dir, fmt.Sprintf("stat_%s.toml", yesterday.Format("20060102")))

	_, err := ioutil.ReadFile(yfn)
	if err != nil {
		// 文件不存在
		return false, nil
	}
	_, err = toml.DecodeFile(yfn, &report)
	if err != nil {
		return false, nil
	}
	return true, report
}

// 把需要发送的信息变成html格式的文本。
func (info *StatisticsReportDaily) GetMailRportHtml(configDir string) (err error, body string) {
	templateFileName := path.Join(configDir, "templateStat.html")
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	t = template.Must(t, err)

	bt := bytes.NewBuffer([]byte{})

	err = t.Execute(bt, info)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	return nil, bt.String()
}
