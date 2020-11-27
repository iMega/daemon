package cmd

import "github.com/sirupsen/logrus"

type config struct {
	Name   string
	Logger *logger
}

type logger struct {
	Channel   string
	BuildID   string
	Formatter logrus.JSONFormatter
}
