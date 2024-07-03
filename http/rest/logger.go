package rest

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Entry {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999999999",
	})
	return logrus.NewEntry(log)
}
