package database

import (
	"cpn-quiz-api-mailer-go/logger"
	"database/sql"

	config "github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//=>for debug sql
	_logger "gorm.io/gorm/logger"
)

type Database struct {
	DB  *gorm.DB
	Log logger.PatternLogger
}

func (d Database) GetConnectionDB() *gorm.DB {
	//=>singleton concept for get db instance.
	if d.DB == nil {
		d.DB = d.initConnect()
	}
	return d.DB
}

func (d Database) initConnect() *gorm.DB {
	var err error
	dsn := config.GetString("cpm.sqlserver.connection.string")
	//=>Re-connection to database.
	d.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("database connection failed : " + err.Error())
	} else {
		d.Log.Info("", "database connection success")
	}

	return d.DB
}

func LoadConfig() {
	//=>Get connection string from config
	dsn := config.GetString("cpm.sqlserver.connection.string")

	//=>Open connection to database for get application config from table.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: _logger.Default.LogMode(_logger.Info),
	})

	if err != nil {
		panic("database connection failed : " + err.Error())
	}

	//=>Query to cq_config table
	result := db.Table("public.cq_config").Select([]string{"key", "value"})
	rows, err := result.Rows()

	if err == nil {
		for rows.Next() {
			//=>Scan key,value from data
			var key, value string
			if err = rows.Scan(
				&key,
				&value,
			); err != nil {
				panic("Error process scan to mapping config: " + err.Error())
			}

			//=>set to config
			config.Set(key, value)
		}
	}

	//=>close if finish work
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
}
