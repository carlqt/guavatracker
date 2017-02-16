package main

import (
	"io"
	"log"
	"os"

	"github.com/pkg/browser"
)

var logger *log.Logger
var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("server.log", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(logFile, os.Stdout)
	logger = log.New(mw, "", log.Ldate|log.Ltime)
}

func main() {
	defer logFile.Close()
	logger.Println("Opening google...")
	browser.OpenURL("http://google.com")
}
