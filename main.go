package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

func initLog(dev bool) {
	const (
		TimeFormatMy = "2006-01-02 15:04:05"
	)

	if dev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.TimeFieldFormat = TimeFormatMy
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // console colored
		log.Logger = log.With().Caller().Logger()
	} else {
		log.Logger = log.Output(&lumberjack.Logger{
			Filename:   "log/can.log", // log file
			MaxSize:    50,            // magebytes
			MaxBackups: 3,             // counts of backup logs
			MaxAge:     7,             // days
			Compress:   false,         // disabled by default
		})
		log.Logger = log.With().Timestamp().Logger()
	}

}

func main() {
	dev := flag.Bool("dev", false, "set log level to debug")
	port := flag.String("port", "8080", "-port specify port")
	flag.Parse()

	initLog(*dev)

	// start http server
	http.ListenAndServe(":"+*port, nil)
}
