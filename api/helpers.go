package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"github.com/sirupsen/logrus"
)

var (
	projectPath, _ = user.Current()
	homePath       = projectPath.HomeDir
	pipelinePath   = filepath.Join(homePath + "/pipeline")
	logPath        = filepath.Join(pipelinePath + "/logfile.log")
)
// Initialize Logrus logging
func (s *server) InitLogrusPipeline() {
	// Create directory to store logs
	if !DirectoryExists(pipelinePath) {
		err := os.Mkdir(pipelinePath, 0777)
		CheckError(err)

		saveFile, err := os.Create(logPath)
		CheckError(err)
		defer saveFile.Close()

		err = os.Chmod(logPath, 0777)
		CheckError(err)
	}
}

// 	Logrus errors
func CheckError(err error) {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05" // Do not change
	Formatter.FullTimestamp = true
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, fileError := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if fileError != nil {
		fmt.Println(fileError.Error())
		// Cannot open log file. Logging to stderr
		if err != nil {
			logrus.Info(err)
			fmt.Println(err)
		}
	} else {
		if err != nil {
			logrus.SetOutput(f)
			logrus.Info(err)
		}
	}
}
// Verify environment variable is set
func env(key, fallbackValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}
	return v
}
// Verify environment variable is set from .env file
func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s required on environment variables", key))
	}
	return v
}
// Verify environment variable is set
func intEnv(key string, fallbackValue int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallbackValue
	}
	return i
}
// FileExists checks if a file exists at a given path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
// DirectoryExists checks if a directory exists at a given path
func DirectoryExists(directory string) bool {
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}