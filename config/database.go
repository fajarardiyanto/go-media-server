package config

import (
	databaseInterface "github.com/fajarardiyanto/flt-go-database/interfaces"
	databaseLib "github.com/fajarardiyanto/flt-go-database/lib"
	"github.com/fajarardiyanto/go-media-server/internal/model"
)

var (
	database databaseInterface.SQL
	rdb      databaseInterface.Redis
	rabbitMQ databaseInterface.RabbitMQ
)

func Config() {
	db := databaseLib.NewLib()
	db.Init(GetLogger())

	InitMysql(db)
	InitRedis(db)
	InitRabbitMQ(db)
}

func InitMysql(db databaseInterface.Database) {
	database = db.LoadSQLDatabase(model.GetConfig().Database.Mysql)

	if err := database.LoadSQL(); err != nil {
		logger.Error(err).Quit()
	}
}

func InitRedis(db databaseInterface.Database) {
	rdb = db.LoadRedisDatabase(model.GetConfig().Database.Redis)

	if err := rdb.Init(); err != nil {
		logger.Error(err)
		return
	}
}

func InitRabbitMQ(db databaseInterface.Database) {
	rabbitMQ = db.LoadRabbitMQ(model.GetConfig().Version, model.GetConfig().Database.RabbitMQ)
}

func GetDBConn() databaseInterface.SQL {
	return database
}

func GetRdbConn() databaseInterface.Redis {
	return rdb
}

func GetRabbitMQ() databaseInterface.RabbitMQ {
	return rabbitMQ
}
