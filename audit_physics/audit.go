package audit_physics

import (
	"fmt"
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
	2、判断是否是offline角色，如果不是offline角色则直接提示报错退出程序
	4、根据端口拼接出data目录
	5、到库中获取所有业务schema、table信息
	6、拼接schema目录路径，计算schema目录空间大小
	7、开启并发线程扫描table行数


 */
func doit(configration *tomlConfig.AUDIT, loggs *log.Logger)  {
	// 获取实例信息
	mysqlInfo := database.GetMySQLInfo(configration)
	if mysqlInfo.Err != nil {
		// 说明获取实例信息有误，打印错误，退出程序
		loggs.Println(mysqlInfo.Err)
		return
	}

	// 从mysql中获取schema、table信息
	endPoint := fmt.Sprintf("%s:%d", mysqlInfo.HostName, mysqlInfo.Port)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", configration.Database.User, configration.Database.Passwd, endPoint, "information_schema", configration.Database.Charset)
	schemaSlice, tableSlice, err := getDbTableInfo(dataSource)
	if err != nil {
		loggs.Println(err)
		return
	}

	// 统计table行数
	dataSlice, err := audit_table(dataSource, mysqlInfo, tableSlice, configration.System.Thread, loggs)
	if err != nil {
		loggs.Println("sorry, audit table is bad, err is", err)
		return
	}


	endPoint = fmt.Sprintf("%s:%d", configration.Database.Address, configration.Database.Port)
	dataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", configration.Database.User, configration.Database.Passwd, endPoint, configration.Database.Schema, configration.Database.Charset)

	// table行数数据入库
	engine, err := database.GetEngine(dataSource)
	if err != nil {
		loggs.Println("sorry, open mysql is bad, err is", err)
		return
	}
	for _, data := range dataSlice {
		_, err := engine.Insert(data)
		if err != nil {
			loggs.Println("sorry, insert data to mysql is bad, err is", err)
		}
	}

	engine.Close()

	// 统计schema大小
	audit_schema(dataSource, mysqlInfo, schemaSlice, loggs)


}