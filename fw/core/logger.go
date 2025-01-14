/* YaNFD - Yet another NDN Forwarding Daemon
 *
 * Copyright (C) 2020 Eric Newberry.
 *
 * This file is licensed under the terms of the MIT License, as found in LICENSE.md.
 */

package core

import (
	"os"

	"github.com/named-data/ndnd/std/log"
)

var Log = log.Default()
var logFileObj *os.File

// OpenLogger initializes the logger.
func OpenLogger(filename string) {
	// open file if filename is not empty
	if filename == "" {
		logFileObj = os.Stderr
	} else {
		var err error
		logFileObj, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}

	// create new logger
	Log = log.NewText(logFileObj)

	// set log level
	level, err := log.ParseLevel(GetConfig().Core.LogLevel)
	if err != nil {
		panic(err)
	}
	Log.SetLevel(level)
}

// ShutdownLogger shuts down the logger.
func CloseLogger() {
	if logFileObj != nil {
		logFileObj.Close()
	}
}
