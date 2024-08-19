package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	logFileName  = "BGSystemLog.log"
	maxFileSize  = 5 * 1024 * 1024 // 5 MB
	backupSuffix = ".gz"
)

type Logger struct {
	file *os.File
}

func (l *Logger) Println(message string) {
	fileInfo, err := l.file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	if fileInfo.Size() >= maxFileSize {
		l.rotateFile()
	}

	_, err = l.file.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message))
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

func (l *Logger) rotateFile() {
	oldFileName := l.file.Name()
	newFileName := fmt.Sprintf("%s.%d", oldFileName, time.Now().Unix())

	err := os.Rename(oldFileName, newFileName)
	if err != nil {
		fmt.Printf("Error renaming log file: %v\n", err)
		return
	}

	go compressFile(newFileName)

	l.file.Close()
	l.file, err = os.OpenFile(oldFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening new log file: %v\n", err)
	}
}

func compressFile(fileName string) {
	inFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file to compress: %v\n", err)
		return
	}
	defer inFile.Close()

	outFile, err := os.Create(fileName + backupSuffix)
	if err != nil {
		fmt.Printf("Error creating compressed file: %v\n", err)
		return
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, inFile)
	if err != nil {
		fmt.Printf("Error compressing file: %v\n", err)
		return
	}

	err = os.Remove(fileName)
	if err != nil {
		fmt.Printf("Error removing old log file: %v\n", err)
	}
}
