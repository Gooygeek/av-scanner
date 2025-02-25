package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"strings"
)

func computeSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func getFileMimeType(filePath string) (string, error) {
	cmd := exec.Command("file", "--mime-type", "--mime-encoding", "-b", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	mimeType := strings.TrimSpace(string(output))
	return mimeType, nil
}
