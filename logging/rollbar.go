package logging

import (
	"fmt"
	"os"

	"github.com/stvp/rollbar"
)

func Critical(err error, fields ...*rollbar.Field) {
	fmt.Println(rollbar.CRIT + ": " + err.Error())
	if os.Getenv("ENV") == "production" {
		rollbar.Error(rollbar.CRIT, err, fields...)
	}
}

func Error(err error, fields ...*rollbar.Field) {
	fmt.Println(rollbar.ERR + ": " + err.Error())
	if os.Getenv("ENV") == "production" {
		rollbar.Error(rollbar.ERR, err, fields...)
	}
}

func Info(err error, fields ...*rollbar.Field) {
	fmt.Println(rollbar.INFO + ": " + err.Error())
	if os.Getenv("ENV") == "production" {
		rollbar.Error(rollbar.INFO, err, fields...)
	}
}
