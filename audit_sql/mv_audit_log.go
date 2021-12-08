package audit_sql

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
该函数的作用将符合要求日志文件从auditDir目录中mv到dataDir，
要求如下：
	1、以audit开头
	2、不以log结尾
	3、不能是目录
 */
func mvAuditLog(auditDir, dataDir string) (auditFileSlice []string , err error) {
	dstAuditFileSlice := []string{}

	info, err := ioutil.ReadDir(auditDir)
	if err != nil {
		return dstAuditFileSlice, err
	}
	// 声明一个切片用于存放符合条件的日志文件
	fileSlice := []string{}
	for _, element := range info {
		if !element.IsDir() {
			// 说明不是目录
			name := element.Name()
			if strings.HasPrefix(name, "audit") {
				// 说明该文件以audit开头
				if !strings.HasSuffix(name, "log") {
					// 说明该文件以audit开头并且以log结尾，符合要求，将该文件加入到切片中
					fileSlice = append(fileSlice, name)
				}
			}
		}
	}

	// mv文件
	for _, name := range fileSlice {
		fileName := fmt.Sprintf("%s/%s", auditDir, name)
		f1, err := os.Open(fileName)
		if err != nil {
			return dstAuditFileSlice, err
		}
		defer f1.Close()
		dstFile := fmt.Sprintf("%s/%s", dataDir, name)
		f2, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return dstAuditFileSlice, err
		}
		defer f2.Close()
		_, err1 := io.Copy(f2, f1)
		if err1 != nil {
			return dstAuditFileSlice, err1
		}
		// 说明数据mv到了dataDir目录，删除该文件
		if err := moveFile(fileName); err != nil {
			return dstAuditFileSlice, err
		}
		dstAuditFileSlice = append(dstAuditFileSlice, dstFile)
	}
	return dstAuditFileSlice, nil
}

func moveFile(file string) error {
	return os.Remove(file)
}
