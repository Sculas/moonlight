// Package global contains some global variables, such as loggers.
package global

import (
	"github.com/sirupsen/logrus"
)

var (
	// Logger is the default logger.
	Logger *logrus.Logger

	ServerLogger *logrus.Entry
	ClientLogger *logrus.Entry
)
