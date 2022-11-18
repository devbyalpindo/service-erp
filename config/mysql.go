package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMySQL() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	connectionString := CONFIG["MYSQL_USER"] + ":" + CONFIG["MYSQL_PASS"] + "@tcp(" + CONFIG["MYSQL_HOST"] + ":" + CONFIG["MYSQL_PORT"] + ")/" + CONFIG["MYSQL_SCHEMA"] + "?parseTime=true"
	mysqlConn, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Println("Error connect to MySQL : ", err.Error())
		return nil, err
	}

	log.Print("MYSQL connection success")
	return mysqlConn, nil
}
