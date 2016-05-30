// 文件 tailMail 处理进度的记录和读取相关的功能封装。
// 进度文件 是系统自动产生的， 不需要配置，
// 进度文件只考虑 toml 格式的， 不支持json格式的。
package progress

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/ghj1976/tailMail"
)

// 文件监控进度配置文件
type TailProgressCollectionEntity struct {
	ProgressMap map[string]int64 // 多个文件进度集合
}

type ProcessInfo struct {
	fileName string
	process  TailProgressCollectionEntity
}

// 传入路径信息，构造进度类，
// 注意，这里传入的是路径信息，而不是文件名
func NewProcessInfo(dir string) *ProcessInfo {
	return &ProcessInfo{
		fileName: path.Join(dir, "progress.toml"),
	}
}

// 读进度文件
func (pi *ProcessInfo) ReadProgress() error {
	pi.process = TailProgressCollectionEntity{
		ProgressMap: map[string]int64{},
	}
	if len(pi.fileName) <= 0 {
		log.Panic("进度文件未定义。")
	}
	_, err := ioutil.ReadFile(pi.fileName)
	if err != nil && os.IsNotExist(err) {
		// 文件不存在，不影响使用
		return nil
	}

	_, err = toml.DecodeFile(pi.fileName, &pi.process)

	// 如果文件读取出错，不能用 nil 指针。
	if pi.process.ProgressMap == nil {
		err = nil
	}
	return nil
}

// 更新写进度文件
func (pi *ProcessInfo) WriteProgress() error {
	if len(pi.fileName) <= 0 {
		log.Panic("进度文件未定义。")
	}
	return tailMail.WriteTOMLFile(pi.fileName, pi.process)
}

// 获得某个文件的最后尺寸
func (pi *ProcessInfo) GetLastFileSize(fn string) int64 {
	pos, b := pi.process.ProgressMap[fn]
	if b {
		return pos
	} else {
		return int64(0)
	}
}

// 设置文件的最后尺寸
func (pi *ProcessInfo) SetLastFileSize(fn string, pos int64) {
	pi.process.ProgressMap[fn] = pos
}

// 更新数据，并更新文件
func (pi *ProcessInfo) UpdateFile(fn string, pos int64) {
	pi.SetLastFileSize(fn, pos)
	// 处理完一个文件，写一次进度，避免中间发生异常时，全部回退。
	err := pi.WriteProgress()
	if err != nil {
		log.Println("写进度配置文件异常：")
		log.Println(err)
	}
}

// 进度文件初始化测试，目前没有意义了。
func saveProgressInfoInit() {
	pi := &ProcessInfo{
		fileName: "11.toml",
		process: TailProgressCollectionEntity{
			ProgressMap: map[string]int64{
				"E:\\tmp\\dnsdata\\auth.acc-2014-01-18-64.35.log": int64(0),
				"E:\\tmp\\dnsdata\\138.18.stdout-2014_01_22.log":  int64(0),
			},
		},
	}

	err := pi.WriteProgress()
	if err != nil {
		fmt.Println(err)
	}
}
