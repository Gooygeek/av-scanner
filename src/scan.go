package main

import (
	"fmt"
	"os"

	"github.com/dutchcoders/go-clamd"
)

func scanFileWithClamAV(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	clamd := clamd.NewClamd("tcp://127.0.0.1:3310")
	response, err := clamd.ScanStream(file, make(chan bool))
	if err != nil {
		return "", "", err
	}

	var scanResult string
	var scanLog string
	for res := range response {
		scanLog += res.Description + "\n"
		scanResult = res.Status
		logger.Debug(fmt.Sprintf("ClamAV scan result:\nRaw:\n%s\nDescription:\n%s\nPath:\n%s\nHash:\n%s\nSize:\n%d\nStatus:\n%s",
			res.Raw, res.Description, res.Path, res.Hash, res.Size, res.Status))
	}

	return scanResult, scanLog, nil
}
