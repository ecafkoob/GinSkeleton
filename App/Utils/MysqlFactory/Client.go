package MysqlFactory

import (
	"GinSkeleton/App/Global/Errors"
	"GinSkeleton/App/Utils/Helper"
	"database/sql"
	"log"
	"time"
)

// 初始化数据库驱动
func Init_sql_driver() *sql.DB {
	configFac := Helper.CreateYamlFactory()
	DbType := configFac.GetString("DbType")
	Host := configFac.GetString("Mysql.Host")
	Port := configFac.GetString("Mysql.Port")
	User := configFac.GetString("Mysql.User")
	Pass := configFac.GetString("Mysql.Pass")
	DataBase := configFac.GetString("Mysql.DataBase")
	SetMaxIdleConns := configFac.GetInt("Mysql.SetMaxIdleConns")
	SetMaxOpenConns := configFac.GetInt("Mysql.SetMaxOpenConns")
	SetConnMaxLifetime := configFac.GetDuration("Mysql.SetConnMaxLifetime")
	db, err := sql.Open(DbType, string(User)+":"+Pass+"@tcp("+Host+":"+Port+")/"+DataBase+"?parseTime=true")
	if err != nil {
		log.Fatal(Errors.Errors_Db_SqlDriverInitFail)
	}
	db.SetMaxIdleConns(SetMaxIdleConns)
	db.SetMaxOpenConns(SetMaxOpenConns)
	db.SetConnMaxLifetime(SetConnMaxLifetime * time.Second)
	return db
}

// 从连接池获取一个连接
func GetOneEffectivePing() *sql.DB {
	configFac := Helper.CreateYamlFactory()
	max_retry_times := configFac.GetInt("Mysql.PingFailRetryTimes")
	// ping 失败允许重试
	v_db_driver := Init_sql_driver()
	for i := 1; i <= max_retry_times; i++ {
		if err := v_db_driver.Ping(); err != nil { //  获取一个连接失败，进行重试
			v_db_driver = Init_sql_driver()
			time.Sleep(time.Second * 1)
			if i == max_retry_times {
				log.Fatal(Errors.Errors_Db_GetConnFail)
			}
		} else {
			break
		}
	}
	return v_db_driver
}
