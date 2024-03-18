package database

import (
	"cpn-quiz-api-file-manage-go/logger"
	"database/sql"

	config "github.com/spf13/viper"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Database struct {
	DB  *gorm.DB
	Log logger.PatternLogger
}

func (d Database) GetConnectionDB() *gorm.DB {
	if d.DB == nil {
		d.DB = d.initConnect()
	}
	return d.DB
}

func (d Database) initConnect() *gorm.DB {
	var err error
	conString := config.GetString("cpm.sqlserver.connection.string")
	d.DB, err = gorm.Open(sqlserver.Open(conString), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("")
	} else {
		d.Log.Info("", "database connection success")
	}

	return d.DB
}

func LoadConfig() {

	conString := config.GetString("cpm.sqlserver.connection.string")
	db, err := gorm.Open(sqlserver.Open(conString), &gorm.Config{})

	if err != nil {
		panic("database connection failed : " + err.Error())
	}

	result := db.Table("CPM.CPM_CONFIG").Select("[KEY],[VALUE]")
	rows, err := result.Rows()

	if err == nil {
		for rows.Next() {
			var key, value string
			if err = rows.Scan(
				&key,
				&value,
			); err != nil {
				panic("Error Process Scan LoadMapping " + err.Error())
			}

			config.Set(key, value)
		}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
}
