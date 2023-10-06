package main

import (
	client "gps_hexagonal/adapters/mysql"
	"gps_hexagonal/adapters/socket"
	"gps_hexagonal/domain"
	"gps_hexagonal/helpers/logger"
)

func init() {
	socket.SERVER_HOST = "192.168.1.3"
	socket.SERVER_PORT = 33333
}

func main() {
	logger.Log.Printf("Start")

	pool, err := client.NewConnectionPool("gps:gps123@tcp(127.0.0.1:3306)/gps_receiver", 3)
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
