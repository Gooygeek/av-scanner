package main

import (
	"net/http"
	"strings"
)

var notifyer Notifyer

type Notifyer struct {
	endpoint string
	topic    string
}

func (n *Notifyer) genericNotify(title, tag, priority, message string) {
	if n.endpoint == "" {
		logger.Debug("Notification endpoint not specified. Skipping notification.")
		return
	}
	req, err := http.NewRequest("POST", n.endpoint+"/"+n.topic, strings.NewReader(message))
	if err != nil {
		logger.Error("Failed to create request: " + err.Error())
	}

	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority)
	req.Header.Set("Tag", tag)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to notify: " + resp.Status)
	}
	logger.Debug("Notification sent successfully")
}

func (n *Notifyer) Debug(message string) {
	title := "DEBUG"
	tag := "bug"
	priority := "low"
	n.genericNotify(title, tag, priority, message)
}

func (n *Notifyer) Info(title, message string) {
	tag := "information_source"
	priority := "default"
	n.genericNotify(title, tag, priority, message)
}

func (n *Notifyer) OK(title, message string) {
	tag := "white_check_mark"
	priority := "default"
	n.genericNotify(title, tag, priority, message)
}

func (n *Notifyer) Warn(title, message string) {
	tag := "warning"
	priority := "default"
	n.genericNotify(title, tag, priority, message)
}

func (n *Notifyer) Detection(title, message string) {
	tag := "bangbang"
	priority := "urgent"
	n.genericNotify(title, tag, priority, message)
}

func (n *Notifyer) Error(message string) {
	title := "ERROR"
	tag := "x"
	priority := "high"
	n.genericNotify(title, tag, priority, message)
}
