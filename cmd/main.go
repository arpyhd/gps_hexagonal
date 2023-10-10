package main

import (
	client "gps_hexagonal/adapters/mysql"
	"gps_hexagonal/adapters/socket"
	"gps_hexagonal/domain"
	"gps_hexagonal/helpers/logger"
	"os"
	"strconv"
)

func init() {
	socket.SERVER_HOST = "0.0.0.0"
	port, err := strconv.Atoi(os.Getenv("GPS_PORT"))
	if err != nil {
		panic(err)
	}
	socket.SERVER_PORT = port

}

func main() {
	logger.Log.Printf("Start")

	pool, err := client.NewConnectionPool("root:"+os.Getenv("MYSQL_ROOT_PASSWORD")+"@tcp(mysql:3306)/"+os.Getenv("MYSQL_DATABASE"), 1)
	if err != nil {
		panic(err)
	}
	defer pool.Pool.Close()

	// acquire a connection from the pool
	db, err := pool.GetConnection()
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	gpsRepository := client.NewGpsAdapter(db)
	gps := domain.NewGpsMap()
	go socket.Worker(gps, gpsRepository)
	socket.Server()

}
