package util

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}
