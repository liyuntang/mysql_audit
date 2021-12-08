package audit_sql

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lifenglin/micro-library/library"
	"log"
	"mysql_audit/database"
	"os"
	"strings"
	"time"
)

type Extra struct {
	UseCache  bool  `json:"use_cache"`
	Expire   uint32 `json:"expire"`
	RefreshCache bool  `json:"refresh_cache"`
	Cluster   string `json:"cluster"`
	OnlyCache  bool  `json:"only_cache"`
	UseLocalCache bool  `json:"use_local_cache"`
	SkipPublish  bool  `json:"skip_publish"`
}

type ReplaceAllSqlTplByProductAndMd5Request struct {
	AllSqlTpl *AllSqlTpl        `json:"all_sql_tpl"`
	Key       map[string]string `json:"key"`
	Extra   *Extra   `json:"extra"`
}

type AllSqlTpl struct {
	Product      string `protobuf:"bytes,2,opt,name=product,proto3" json:"product"`
	Cluster_name string `protobuf:"bytes,3,opt,name=cluster_name,proto3" json:"cluster_name"`
	Role         string	`protobuf:"bytes,4,opt,name=role,proto3" json:"role"`
	Host_name    string	`protobuf:"bytes,5,opt,name=host_name,proto3" json:"hostname"`
	Port         int	`protobuf:"bytes,6,opt,name=port,proto3" json:"port"`
	Md5          string `protobuf:"bytes,7,opt,name=md5,proto3" json:"md5"`
	Sql_text	string	`protobuf:"bytes,8,opt,name=sql_text,proto3" json:"sql_text"`
	Tpl          string `protobuf:"bytes,9,opt,name=tpl,proto3" json:"tpl"`
	Backtrace    string `protobuf:"bytes,10,opt,name=backtrace,proto3" json:"backtrace"`
	Ip           string `protobuf:"bytes,11,opt,name=ip,proto3" json:"ip"`
	Db           string `protobuf:"bytes,12,opt,name=db,proto3" json:"db"`
	User         string `protobuf:"bytes,13,opt,name=user,proto3" json:"user"`
	Num          int	`protobuf:"bytes,14,opt,name=num,proto3" json:"num"`
}

type auditRecord struct {
	Audit_record recordData
}

type recordData struct {
	Timestamp, Command_class, Sqltext, User, Ip, Db string
}

func ana(auditLogFileSlice []string, info *database.MySQLInfo, retry int, loggs *log.Logger) {
	allSqlTpl := new(AllSqlTpl)
	allSqlTpl.Product = info.Product
	allSqlTpl.Cluster_name = info.ClusterName
	allSqlTpl.Host_name = info.HostName
	allSqlTpl.Role = info.Role
	allSqlTpl.Port = info.Port
	mapSqlTpl := map[string]AllSqlTpl{}
	tmpSlice := []string{}
	year, oweek := time.Now().ISOWeek()
	week := fmt.Sprintf("%d_%d", year, oweek)
	var record *auditRecord

	for _, auditLogFile := range auditLogFileSlice {
		file, err := os.Open(auditLogFile)
		if err != nil {
			loggs.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		
		for scanner.Scan() {
			buf := scanner.Bytes()
			json.Unmarshal(buf, &record)
			
			switch strings.ToLower(record.Audit_record.Command_class) {
			case "set_option":
				continue
			case "error":
				continue
			case "show_status":
				continue
			case "show_slave_status":
				continue
			case "show_variables":
				continue
			case "commit":
				continue
			default:
			}
			allSqlTpl.User = strings.Split(strings.Split(record.Audit_record.User, " ")[0], "[")[0]
			allSqlTpl.Db = record.Audit_record.Db

			allSqlTpl.Ip = record.Audit_record.Ip
			// sql语句
			allSqlTpl.Sql_text = record.Audit_record.Sqltext

			// sql 指纹
			allSqlTpl.Tpl, _ = library.TransferSQLToTpl(allSqlTpl.Sql_text)
			strMd5 := allSqlTpl.Product + "_" + allSqlTpl.Cluster_name + "_" + allSqlTpl.Db + "_" + allSqlTpl.User + "_" + allSqlTpl.Tpl + "_" + week
			tmpString := md5.Sum([]byte(strMd5))
			allSqlTpl.Md5 = hex.EncodeToString(tmpString[:])
			if data, ok := mapSqlTpl[allSqlTpl.Md5]; !ok {
				// 说明不存在该md5值
				allSqlTpl.Num = 1
				mapSqlTpl[allSqlTpl.Md5] = *allSqlTpl
				tmpSlice = append(tmpSlice, allSqlTpl.Md5)
			} else {
				// 说明该sql指纹已经存在了，数量加1
				data.Num +=1
				mapSqlTpl[allSqlTpl.Md5] = data
			}
		}

	}

	// 删除文件
	for _, auditlog := range auditLogFileSlice {
		if err := moveFile(auditlog); err != nil {
			loggs.Println("删除", auditlog, "失败")
		} else {
			loggs.Println("删除", auditlog, "成功")
		}
	}

	// 上传数据到微服务
	rows := len(mapSqlTpl)
	rowNum := 1
	for _, allSqlTpl := range mapSqlTpl {
		err := post(allSqlTpl)
		if err != nil {
			loggs.Println(fmt.Sprintf("一共需要上传%d条数据，第%d条数据上传失败，进入重试流程", rows, rowNum))
		} else {
			//loggs.Println(fmt.Sprintf("一共需要上传%d条数据，第%d条数据上传成功", rows, rowNum))
			rowNum += 1
			continue
		}
		// 重试 这是不是有点娘儿们
		for i:=1;i<=retry;i++ {
			err := post(allSqlTpl)
			if err != nil {
				loggs.Println(fmt.Sprintf("一共需要上传%d条数据，第%d条数据上传失败，进行第%d次重试", rows, rowNum, i))
				time.Sleep(20 * time.Millisecond)
			} else {
				//loggs.Println(fmt.Sprintf("一共需要上传%d条数据，第%d条数据上传成功", rows, rowNum))
				rowNum += 1
				continue
			}
		}
		loggs.Println(fmt.Sprintf("一共需要上传%d条数据，第%d条数据进行%d次重试仍然失败,请排查下游服务状态", rows, rowNum, retry))
		rowNum += 1
	}
}
