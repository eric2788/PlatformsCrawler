package logging

import (
	"context"
	"github.com/sirupsen/logrus"
)

func GetCrawlerLogger(name string) *logrus.Entry {
	return logrus.WithField("crawler", name)
}

func GetMainLogger() *logrus.Entry {
	return logrus.WithContext(context.Background())
}
