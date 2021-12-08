package audit_sql

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func post(allSqlTpl AllSqlTpl) error {
	//fmt.Println("上传文件")
	//return nil
	req := new(ReplaceAllSqlTplByProductAndMd5Request)
	req.Extra = new(Extra)
	req.Extra.SkipPublish = true
	req.AllSqlTpl = &allSqlTpl
	req.Key = make(map[string]string, 0)
	req.Key["product"] = allSqlTpl.Product
	req.Key["md5"] = allSqlTpl.Md5

	jsonReq, _ := json.Marshal(map[string]interface{}{
		"service":	"go.microv2.srv.SlowSql",
		"endpoint": "AllSqlTplSRV.StoreAuditData",
		"request":	req,
	})
	body := strings.NewReader(string(jsonReq))
	curlreq, err := http.NewRequest("POST", "http://micro.service.niceprivate.com/rpc", body)
	if err != nil {
		return err
	}
	curlreq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(curlreq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		//byte数组直接转成string，优化内存
		str := (*string)(unsafe.Pointer(&respBytes))
		return errors.New(fmt.Sprintf("call failed, code:", resp.StatusCode, ", body:", *str))
	}
	return nil
}
