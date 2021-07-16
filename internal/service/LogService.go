package service

import (
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	LogFile = "/var/log/parser.log"
	DefaultDetailLevel = 0
)

var (
	file *os.File
	fileOpenError error
	detailLevel = DefaultDetailLevel
)

func setDetailLevel()  {
	debugLevel, exists := os.LookupEnv("DEBUG_LEVEL")
	debugLevelInt, err := strconv.Atoi(debugLevel)
	if exists && err == nil {
		detailLevel = debugLevelInt
	}
}

func getLogFile() (*os.File, error) {
	var once sync.Once

	once.Do(func() {
		setDetailLevel()
		file, fileOpenError = os.OpenFile(LogFile, os.O_APPEND | os.O_CREATE, 0644)
	})

	return file, fileOpenError
}

func ToLog(message string) error {
	file, err := getLogFile()
	if err != nil {
		return err
	}

	iLog := log.New(file, "", log.LstdFlags)
	switch detailLevel {
	case 1:
		iLog.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	iLog.Println(message)

	return nil
}

func CloseLogFile() error {
	if file != nil {
		return file.Close()
	}

	return nil
}
