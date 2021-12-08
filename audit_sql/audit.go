package audit_sql

import (
	"log"
	"mysql_audit/database"
	"mysql_audit/tomlConfig"
)

func StartAuditSql(configration *tomlConfig.AUDIT, loggs *log.Logger)  {
	defer func() {
		info := recover()
		if info != nil {
			loggs.Println("程序意外退出，info is", info)
		}
	}()
	doit(configration, loggs)
}

/*
开始审计，流程如下：
	1、获取主机名、从配置中读取端口，根据这两个条件从mysql_cluster表中获取product、cluster、role信息
	2、根据端口拼接出audit_log目录
	3、从audit_log目录中复制文件到dataDir目录，判断条件如下：
		1）以audit开头
		2）不以log结尾
		3) 不能是目录
	4、对dataDir目录的文件进行处理


 */
func doit(configration *tomlConfig.AUDIT, loggs *log.Logger)  {
	// 获取实例信息
	mysqlInfo := database.GetMySQLInfo(configration)
	if mysqlInfo.Err != nil {
		// 说明获取实例信息有误，打印错误，退出程序
		loggs.Println(mysqlInfo.Err)
		return
	}
	// 说明获取信息正常，遍历auditDir目录，将审计日志mv到dataDir
	auditFileSlice, err := mvAuditLog(mysqlInfo.AuditDir, configration.System.DataDir)
	if err != nil {
		loggs.Println(err)
		return
	}
	// 判断是否有文件需要处理
	if len(auditFileSlice) <= 0 {
		loggs.Println("没有audit log需要处理")
		return
	}

	// 处理文件
	ana(auditFileSlice, mysqlInfo, configration.System.Retry, loggs)
}