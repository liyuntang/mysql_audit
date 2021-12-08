package audit_physics

type Mysql_audit_schema_size struct {
	Product, Cluster_name, Host_name string
	Port int
	Db_name string
	Db_size int64
}

type Mysql_audit_table_size struct {
	Product, Cluster_name, Host_name string
	Port int
	Db_name, Table_name string
	Table_rows int64
}