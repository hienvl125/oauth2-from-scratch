package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hienvl125/oauth2-from-scratch/oauth2-client/config"
	"github.com/jmoiron/sqlx"
)

func NewMySqlxDB(conf *config.Config) (*sqlx.DB, error) {
	// DSN pattern: "my_db_user:my_db_password@(my_db_host:my_db_port)/my_db_name"
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
	)
	return sqlx.Connect("mysql", dsn)
}
