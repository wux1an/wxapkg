package util

import (
	"github.com/fatih/color"
	"os"
)

func Fatal(err error) {
	if err == nil {
		return
	}

	color.Red("%v", err)
	os.Exit(0)
}
