package actions

import (
	"mailAssistant/logging"
	"reflect"
	"runtime"
)

type jobCallBack func(job Job, waitGroup *int32)

var regLog = logging.NewNamedLogger("${project}.actions.registry")

var (
	actions = make(map[string]jobCallBack, 0)
)

func register(name string, callback jobCallBack) {
	regLog.Debug("register ", name)
	if fcc, found := actions[name]; found {
		callbackName := runtime.FuncForPC(reflect.ValueOf(callback).Pointer()).Name()
		regLog.Severe("duplicate action name", name, " => ", callbackName)
		fccName := runtime.FuncForPC(reflect.ValueOf(fcc).Pointer()).Name()
		regLog.Severe("duplicate action name", name, " => ", fccName)
	}
	actions[name] = callback
}
