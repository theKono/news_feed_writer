package mysql

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/theKono/orchid/cfg"
)

var (
	// DBSessions specifies the 2 MySQL connections in terms of sharding.
	DBSessions = []*gorm.DB{
		ConnectMysqlDb(cfg.MysqlMain),
		ConnectMysqlDb(cfg.MysqlShard),
	}
)

// ConnectMysqlDb connects to MySQL instance.
//
// If fail to make connection, ConnectMysqlDb panics.
// If MYSQL_DEBUG is set to 1, the logging will display more
// information.
func ConnectMysqlDb(dataSourceName string) *gorm.DB {
	db, err := gorm.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatalf("Cannot connect to MySQL: %v\n%v", dataSourceName, err)
	}

	if cfg.MysqlDebug {
		db.LogMode(true)
	}

	return db
}
