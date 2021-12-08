package audit_physics

import (
	"fmt"
	"io/ioutil"
	"log"
	"mysql_audit/database"
)

func audit_schema(dataSource string, mysqlInfo *database.MySQLInfo, schemaSlice []string, loggs *log.Logger)  {
	dataSlice := []*Mysql_audit_schema_size{}
	for _, db := range schemaSlice {
		dbDir := fmt.Sprintf("%s/%s", mysqlInfo.DataDir, db)

		info, err := ioutil.ReadDir(dbDir)
		if err != nil {
			loggs.Println(fmt.Sprintf("sorry, open schema dir %s is bad, err is %v", dbDir, err))
			return
		}
		var size int64
		for _, file := range info {
			size+=file.Size()
		}
		dataSlice = append(dataSlice, &Mysql_audit_schema_size{
			Product: mysqlInfo.Product,
			Cluster_name: mysqlInfo.ClusterName,
			Host_name: mysqlInfo.HostName,
			Port: mysqlInfo.Port,
			Db_name: db,
			Db_size: size,
		})
	}
	// 入库
	engine, err := database.GetEngine(dataSource)
	if err != nil {
		loggs.Println("sorry, open mysql is bad, err is", err)
		return
	}
	defer engine.Close()

	rows, err := engine.Insert(dataSlice)
	if err != nil {
		loggs.Println("sorry, insert data to mysql is bad, err is", err)
		return
	}
	loggs.Println("insert data to mysql is ok, rows is", rows)
}
