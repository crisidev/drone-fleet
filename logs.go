package main

import (
	"os"

	"github.com/op/go-logging"
)

func logSetup(debug bool) {
	logFormat := logging.MustStringFormatter(`%{color}%{level:-7.7s} %{shortfunc:-15.15s} %{color:reset} %{message}`)
	outBackend := logging.NewLogBackend(os.Stderr, "", 0)
	outBackendFormatter := logging.NewBackendFormatter(outBackend, logFormat)
	outBackendLeveled := logging.AddModuleLevel(outBackendFormatter)
	if debug {
		outBackendLeveled.SetLevel(logging.DEBUG, "")
	} else {
		outBackendLeveled.SetLevel(logging.INFO, "")
	}
	logging.SetBackend(outBackendLeveled)
}
