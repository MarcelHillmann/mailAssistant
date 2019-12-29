package cntl

import (
	"mailAssistant/logging"
	"os"
	"os/signal"
)

var (
	done      = make(chan bool)
	osChannel chan os.Signal
	count     = 0
	groupLog  *logging.Logger
)

// Notify is sending so much as needed bool to a channel
func Notify() {
	// close(done)
	for i := 0; i < count; i++ {
		done <- true
	}
	count = 0
}

// WaitForNotify is blocking until a notify is received
func WaitForNotify() {
	count++
	<-done
}

// ToNotify is returning the number of waiting threads
func ToNotify() int {
	return count
}

// WaitForOsNotify is blocking until a os signal is received and inform all threads
func WaitForOsNotify(signals ...os.Signal) {
	osChannel = make(chan os.Signal, len(signals))
	signal.Notify(osChannel, signals...)
	go func() {
		signal := <-osChannel
		getLogger().Severe("Got signal: ", signal)
		Notify()
	}()
}

func getLogger() *logging.Logger {
	if groupLog == nil {
		groupLog = logging.NewLogger()
	}
	return groupLog
}
