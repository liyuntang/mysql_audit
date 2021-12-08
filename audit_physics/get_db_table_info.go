package audit_physics

import "mysql_audit/database"

func getDbTableInfo(dataSource string) (dbSlice, tabSlice []string, e error) {
	engine, err := database.GetEngine(dataSource)
	if err != nil {
		return nil, nil, err
	}
	defer engine.Close()
	schemaSlice := []string{}
	tableSlice := []string{}
	// 获取shema信息
	sql1 := "select distinct table_schema from tables where table_schema not in ('mysql', 'information_schema', 'performance_schema', 'test', 'sys');"
	r1, err := engine.QueryString(sql1)
	if err != nil {
		return schemaSlice, tableSlice, err
	}
	for _, dict := range r1 {
		schemaSlice = append(schemaSlice, dict["table_schema"])
	}
	// 获取table信息
	sql2 := "select table_schema, table_name from tables where table_schema not in ('mysql', 'information_schema', 'performance_schema', 'test', 'sys');"
	r2, err := engine.QueryString(sql2)
	if err != nil {
		return schemaSlice, tableSlice, err
	}
	for _, dict := range r2 {
		tableSlice = append(tableSlice, dict["table_schema"]+"."+dict["table_name"])
	}
	return schemaSlice, tableSlice, nil
}
