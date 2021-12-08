package audit_physics

import (
	"fmt"
	"log"
	"mysql_audit/database"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup


func audit_table(dataSource string, mysqlInfo *database.MySQLInfo, tableSlice []string, thread int, loggs *log.Logger) (dataSlice []*Mysql_audit_table_size, e error) {
	engine, err := database.GetEngine(dataSource)
	if err != nil {
		return nil, err
	}
	defer engine.Close()

	tableAuditSlice := []*Mysql_audit_table_size{}
	// 声明两个chnnel
	taskChannel := make(chan *Mysql_audit_table_size)
	resultChannel := make(chan *Mysql_audit_table_size)

	// 从resultChannel中读取数据，并将数据放入tableAuditSlice切片
	go func() {
		for result := range resultChannel {
			tableAuditSlice = append(tableAuditSlice, result)
		}
	}()

	// 开始并发审计数据表

	for i:=1;i<=thread;i++ {
		wg.Add(1)

		go func(wait *sync.WaitGroup) {
			defer wait.Done()
			for task := range taskChannel {
				tableName := fmt.Sprintf("%s.%s", task.Db_name, task.Table_name)
				sql := fmt.Sprintf("select count(*) from %s;", tableName)
				result, err := engine.QueryString(sql)
				if err != nil {
					loggs.Println("sorry, exec sql", sql, "is bad, err is", err)
					return
				}
				sizeString := result[0]["count(*)"]
				size , err := strconv.ParseInt(sizeString, 10, 64)
				if err != nil {
					loggs.Println("sorry, transfer table size of", tableName, "is bad, err is", err)
					return
				}
				task.Table_rows = size
				resultChannel <- task
			}

		}(&wg)
	}
	// 分发任务
	count := len(tableSlice)
	num := 1
	for _, tableName := range tableSlice {
		data := &Mysql_audit_table_size{
			Product: mysqlInfo.Product,
			Cluster_name: mysqlInfo.ClusterName,
			Host_name: mysqlInfo.HostName,
			Port: mysqlInfo.Port,
			Db_name: strings.Split(tableName, ".")[0],
			Table_name: strings.Split(tableName, ".")[1],
		}
		loggs.Println(fmt.Sprintf("一共需要审计%d张数据表，现在开始审计第%d张数据表", count, num))
		taskChannel <- data
		num+=1

	}


	close(taskChannel)
	wg.Wait()
	close(resultChannel)

	return tableAuditSlice, nil

}
