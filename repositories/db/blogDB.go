package db

import (
	"database/sql"

	"github.com/L-oris/yabb/foreign/env"
	"github.com/L-oris/yabb/logger"
	_ "github.com/go-sql-driver/mysql"
)

// BlogDB holds the connection to Blog DB
var BlogDB *sql.DB

func init() {
	var err error
	BlogDB, err = sql.Open("mysql", env.Vars.DB+"?parseTime=true")
	if err != nil {
		logger.Log.Fatal("db connection error: ", err)
	}

	if err := BlogDB.Ping(); err != nil {
		logger.Log.Fatal("db ping error: ", err)
	}

	logger.Log.Debug("blogDB connected")
}
