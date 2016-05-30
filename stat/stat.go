package stat

import (
	"fmt"
	"path"
	"time"
)

type StatInfo struct {
	statFileName string                 // toml 格式的统计数据文件
	DayStat      *StatisticsReportDaily // 日统计
}

// 构造统计类
func NewStatInfo(dir string) *StatInfo {
	fn := path.Join(dir, fmt.Sprintf("stat_%s.toml", time.Now().Format("20060102")))

	si := &StatInfo{
		statFileName: fn,
	}
	return si

}
