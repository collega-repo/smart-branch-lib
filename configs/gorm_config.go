package configs

import (
	"database/sql"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

var DB *gorm.DB

func NewDBConn() {
	conf := commons.Configs.Datasource.DB
	dsn := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=%v search_path=%s`,
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Database, conf.Sslmode, conf.Schema)
	if strings.ToUpper(conf.Driver) == "PGX" {
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?search_path=%s&sslmode=%v",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Database, conf.Schema, conf.Sslmode)
	}

	db, err := sql.Open(conf.Driver, dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(conf.MaxIdle)
	db.SetConnMaxIdleTime(conf.MaxIdleTime * time.Minute)
	db.SetMaxOpenConns(conf.MaxOpen)
	db.SetConnMaxLifetime(conf.MaxLifetime * time.Minute)

	if conf.Ping {
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	}

	if DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}),
		&gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
			DisableForeignKeyConstraintWhenMigrating: true,
		}); err != nil {
		panic(err)
	}

	if conf.Debug {
		DB = DB.Debug()
	}
	DB.Callback().Query()
}
