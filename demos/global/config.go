package global

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	RdbClient *redis.Client
	// MysqlClient *gorm.DB
)

func init() {
	RdbClient = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "redis123",
			DB:       0,
		},
	)

	c := context.Background()
	if err := RdbClient.Ping(c).Err(); err != nil {
		panic(err)
	}
}

// func init() {
// 	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
// 	dsn := "root:Mysql123@tcp(127.0.0.1:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	MysqlClient = db
// }
