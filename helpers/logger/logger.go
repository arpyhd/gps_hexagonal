package logger

import (
	//	"flag"

	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file
	var logpath = "gps_gpshexagonal.log"

	//	flag.Parse()
	//var file, err1 = os.Create(logpath)
	file, err1 := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	//   Log.Println("LogFile : " + logpath)
}
