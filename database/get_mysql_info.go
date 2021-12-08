package database

import (
	"errors"
	"fmt"
	"mysql_audit/tomlConfig"
	"os"
	"strings"
)

type MySQLInfo struct {
	Product, ClusterName, Role string
	HostName string
	AuditDir, DataDir string
	Port int
	Err error
}

func GetMySQLInfo(configration *tomlConfig.AUDIT) *MySQLInfo {
	// 获取主机名及端口
	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]
	//hostName = "s-mysql-core04"
	port := configration.System.Port
	info := &MySQLInfo{}
	endPoint := fmt.Sprintf("%s:%d", configration.Database.Address, configration.Database.Port)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", configration.Database.User, configration.Database.Passwd, endPoint, configration.Database.Schema, configration.Database.Charset)
	engine, err := GetEngine(dataSource)

	if err != nil {
		info.Err = err
		return info
	}
	defer engine.Close()
	sql := fmt.Sprintf("select product, cluster_name, role from mysql_cluster where hostname='%s' and port=%d;", hostName, port)
	//fmt.Println(sql)
	resultSlice, err := engine.QueryString(sql)
	if err != nil {
		info.Err = err
		return info
	}
	// 判断mysql_cluster表中的数据是否正确，
	num := len(resultSlice)
	if num != 1 {
		info.Err = errors.New(fmt.Sprintf("sorry, 主机%s,端口%d对应了%d个mysql实例", hostName, port, num))
		return info
	}
	// 说明mysql_cluster表里的数据没问题，解析数据
	for _, dict := range resultSlice {
		info.Product = dict["product"]
		info.ClusterName = dict["cluster_name"]
		info.Role = dict["role"]
		info.HostName = hostName
		info.Port = port
		info.AuditDir = fmt.Sprintf("/home/mysql/mysql_%d/audit_log", port)
		info.DataDir = fmt.Sprintf("/home/mysql/mysql_%d/data", port)
		//info.AuditDir = "/Users/liyuntang/audit_log"
		//info.DataDir = "/Users/liyuntang/data"
		info.Err = nil
	}
	return info
}

