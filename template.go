package tailMail

import (
	"bytes"
	"log"
	"text/template"
	"time"
)

// 文件模版用的格式化时间输出的函数。
func FormatNow(format string) string {
	return time.Now().Format(format)
}

// 文件名模版替换函数。
func FormatFileName(tfilename string) string {
	tmpl := template.New("t1")                                  // 模版名字
	tmpl = tmpl.Funcs(template.FuncMap{"formatNow": FormatNow}) // 模版使用的函数

	tmpl, err := tmpl.Parse(tfilename) // 模版
	if err != nil {
		log.Println("监控文件名模版配置错误：001 ", tfilename)
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, "")
	if err != nil {
		log.Println("监控文件名模版配置错误：002 ", tfilename)
		panic(err)
	}
	return doc.String()
}
