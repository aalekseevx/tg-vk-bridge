package bridge

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Start() {
	initRedis()
	initDatabase()

	// Start web api in new routine
	go mainApi()

	// Start telegram polling in current routine
	mainBot()
}
