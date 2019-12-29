package cntl

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log"
	"mailAssistant/logging"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGroupCntl_Notify(t *testing.T){
	done = make(chan bool)
	got:=0
	go func(){
		for <-done {
			got++
		}
	}()

	count=2
	Notify()
	close(done)
	time.Sleep(100 *time.Millisecond)
	require.Equal(t, 2, got)
}

func TestGroupCntl_WaitForNotify(t *testing.T){
	count=0
	done = make(chan bool)

	go WaitForNotify()
	time.Sleep(100 * time.Millisecond)

	require.Equal(t, 1, count)
	close(done)
}

func TestGroupCntl_ToNotify(t *testing.T){
	defer func(){
		count=0
	}()

	done = make(chan bool)
	count = 4
	require.Equal(t, 4, ToNotify())
}

func TestGroupCntl_WaitForOsNotify(t *testing.T) {
	defer func() {
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	}()

	done = make(chan bool)
	byt := bytes.NewBufferString("")
	log.SetFlags(0)
	log.SetOutput(byt)
	logging.SetLevel("mailAssistant.cntl","all")
	go func(){
		time.Sleep(time.Second)
		osChannel <- os.Interrupt
	}()
	WaitForOsNotify(os.Interrupt, syscall.SIGTERM)

	require.Equal(t, 0, ToNotify())
	WaitForNotify()
	require.Equal(t, "SEVERE  [mailAssistant.cntl#waitforosnotify.func1] Got signal:  interrupt\n", byt.String())
}