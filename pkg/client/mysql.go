package client

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewMySQLClient() *gorm.DB {
	//file := logger2.CreateFileWriter(conf.Log.LogFilePath("slow-sql.log"))

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
		//	SlowThreshold:             5 * time.Millisecond,
		//	LogLevel:                  logger.Warn,
		//	IgnoreRecordNotFoundError: true,
		//}),
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?&parseTime=True&loc=Local",
			"root",
			"123456",
			"0.0.0.0",
			3306,
			"k",
		),
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), gormConfig)

	if err != nil {
		panic(fmt.Errorf("mysql connect error :%v", err))
	}

	if db.Error != nil {
		panic(fmt.Errorf("database error :%v", err))
	}

	_, _ = db.DB()

	fmt.Println("successful to connect db")
	return db
}
