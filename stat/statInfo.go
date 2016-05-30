// 统计报表相关实体类、功能封装。
package stat

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ghj1976/tailMail"
)

// 基于某个文件的统计信息
type StatisticsReportFile struct {
	FileName string // 文件名
	Subject  string // 监控邮件标题
	Num      int    // 次数
}

// 每个邮箱收到的报告统计
type StatisticsReportEMail struct {
	EMailName    string                          // 邮箱名称
	EMailAddress string                          // 邮箱地址
	Num          int                             // 今天本邮箱合计的报警数量
	FileStatMap  map[string]StatisticsReportFile // 每个邮箱的统计报告
}

// 每日统计报告类
type StatisticsReportDaily struct {
	CurrDay      time.Time                        // 哪一天的统计报告，小时、分、秒对这里没有意识，随便写。
	ServerName   string                           // 服务器名称
	Num          int                              // 今天本服务器的合计的报警数量
	EmailStatMap map[string]StatisticsReportEMail // 每个邮箱的统计报告
}

// 恢复初始设置。
func reset(stat *StatisticsReportDaily) {
	if stat == nil {
		return
	}
	stat.CurrDay = time.Now()
	stat.Num = 0
	stat.EmailStatMap = make(map[string]StatisticsReportEMail)
}

// 获得一个当前的统计时间
func (si *StatInfo) GetCurrDayStat(serverName string) {
	cn := time.Now()
	if si.DayStat != nil {
		ln := si.DayStat.CurrDay
		if cn.Year() == ln.Year() && cn.Month() == ln.Month() && cn.Day() == ln.Day() {
			return
		} else {
			// 不是当天，先保存
			si.WriteDayStatValue()
			// 再初始化成当天
			si.DayStat = &StatisticsReportDaily{
				ServerName: serverName,
			}
			reset(si.DayStat)
		}

	} else {
		// 读配置文件
		si.readDayStatValue(serverName)
	}
}

// 读统计结果数据文件， 如果文件不存在，则返回一个默认数据的，以便可以回写
func (si *StatInfo) readDayStatValue(serverName string) error {
	_, err := ioutil.ReadFile(si.statFileName)
	if err != nil && os.IsNotExist(err) {
		// 文件不存在，不影响使用
		si.DayStat = &StatisticsReportDaily{
			ServerName: serverName,
		}
		reset(si.DayStat)
		return nil
	}

	_, err = toml.DecodeFile(si.statFileName, &si.DayStat)
	return err
}

// 更新写 统计结果
func (si *StatInfo) WriteDayStatValue() (err error) {
	return tailMail.WriteTOMLFile(si.statFileName, si.DayStat)
}
