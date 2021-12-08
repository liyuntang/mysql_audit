package main

import (
	"flag"
	"fmt"
	"mysql_audit/audit_physics"
	"mysql_audit/audit_sql"
	"mysql_audit/common"
	"mysql_audit/tomlConfig"
	"os"
	"time"
)

var conf string
var audit_type string
var run_time string

func init()  {
	flag.StringVar(&conf, "c", "conf/audit.toml", "配置文件")
	flag.StringVar(&audit_type, "t", "audit_sql", "审计类型，包含audit_sql、audit_physics两种,audit_sql表示审计sql语句,audit_physics表示审计库、表行数、大小等")
	flag.StringVar(&run_time, "time", "12:00", "开始时间，只作用于audit_physics类型的审计，格式为hh:mm,24小时制")
}

func main()  {
	flag.Parse()

	// 解析配置文件
	configration := tomlConfig.TomlConfig(conf)

	// 获取日志句柄
	loggs := common.WriteLog(configration.System.LogFile)

	// 判读审计方式，

	switch audit_type {
	case "audit_sql":
		//loggs.Println("审计方式为audit_sql")
		for {
			audit_sql.StartAuditSql(configration, loggs)
			//os.Exit(0)
			time.Sleep(configration.System.IntervalTime * time.Second)
		}
	case "audit_physics":
		//loggs.Println("审计方式为audit_physics")
		for {
			// 检查时间
			if checkTime(run_time) {
				audit_physics.StartAuditSql(configration, loggs)
			}
			time.Sleep(1 * time.Minute)
		}
	default:
		loggs.Println("审计方式输入错误,当前仅支持sql、table两种,sql表示审计sql语句,table表示审计库、表行数、大小等")
		os.Exit(0)
	}

}

func checkTime(runTime string) bool {
	t, err := time.Parse("15:04", run_time)
	if err != nil {
		fmt.Println("sorry, parse run time is bad, err is", err)
		os.Exit(0)
	}

	h := t.Hour()
	m := t.Minute()
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()

	if h == hour && m == minute {
		fmt.Println("runtime is", runTime, "now is", now.Format("15:04"),"符合运行时间,开始审计")
		return true
	}
	fmt.Println("runtime is", runTime, "now is", now.Format("15:04"),"不符合运行时间")
	return false
}