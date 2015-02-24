package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	l, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(l)
}
