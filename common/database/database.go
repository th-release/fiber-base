package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"cth.release/common"
	_ "github.com/lib/pq"
)

func InitDatabase() *sql.DB {
	config := common.GetConfig()

	dbPort, err := strconv.Atoi(config.DatabasePort)

	if err != nil {
		return nil
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DatabaseHost, dbPort, config.DatabaseUser, config.DatabasePass, config.DatabaseName)

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		return nil
	}

	return db
}
