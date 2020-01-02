package actions

import (
	"github.com/emersion/go-imap"
	"mailAssistant/account"
	"mailAssistant/arguments"
	"mailAssistant/logging"
	"mailAssistant/planning"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// NewJob is a job factory
func NewJob(jobName, name string, args map[string]interface{}, accounts *account.Accounts) Job {
	log := logging.NewNamedLogger("${project}.actions")
	fcc, ok := actions[name]
	if !ok {
		log.Severe("unknown action ", name, "for", jobName)
		fcc = newDummy
	}
	loggerName := runtime.FuncForPC(reflect.ValueOf(fcc).Pointer()).Name()
	log.Info("action ", loggerName, "for", jobName)

	semaphore[jobName] = semaphoreNull()
	return Job{Args: arguments.NewArgs(args), log: logging.NewNamedLogger(loggerName), callback: fcc, accounts: accounts, jobName: jobName}
}

func semaphoreNull() *int32 {
	result := Released
	return &result
}

var semaphore = make(map[string]*int32)

// Job represents a job for scheduling
type Job struct {
	*arguments.Args
	log      *logging.Logger
	callback jobCallBack
	accounts *account.Accounts
	jobName  string
	saveTo   string
}

// Run is called by clockwerk framework
func (j Job) Run() {
	j.log.Enter()
	j.callback(j, semaphore[j.jobName])
	j.log.Leave()
}

// GetAccount is checking and returning the searched account
func (j Job) GetAccount(name string) *account.Account {
	if ! j.accounts.HasAccount(name)  {
		j.log.Severe(name,"is not defined")
	}
	return j.accounts.GetAccount(name)
}

func (j *Job) getSaveTo() string {
	if j.saveTo != "" {
		return j.saveTo
	}
	saveTo := ""
	if runtime.GOOS == "windows" {
		saveTo = j.GetString("saveToWin")
	}
	if saveTo == "" {
		saveTo = j.GetString("saveTo")
	}

	j.saveTo = saveTo
	return j.saveTo
}

const mailPrefix = "mail_"

func (j Job) getSearchParameter() [][]interface{} {
	result := make([][]interface{}, 0, 0)
	for _, key := range j.GetArgKeys() {
		lKey := strings.ToLower(key)
		switch {
			case lKey == "mail_account":
				// ignore always
			case lKey == "mail_older":
				value := j.GetString(key)
				duration := planning.ParseSchedule(value)
				before := time.Now().Unix() - int64(duration.Seconds())
				arg := append(make([]interface{}, 0), "before", time.Unix(before, 0).Format(imap.DateLayout))
				result = append(result, arg)
			case lKey == "mail_or" ||lKey == "mail_not":
				orList := j.GetList(key)
				arg := append(make([]interface{}, 0), strings.TrimPrefix(lKey, mailPrefix))
				arg = parseRecursive(arg, orList)
				result = append(result, arg)
			case strings.HasPrefix(lKey, mailPrefix):
				arg := append(make([]interface{}, 0), strings.TrimPrefix(lKey, mailPrefix), j.GetArg(key))
				result = append(result, arg)
			default:
				//
		} // switch
	}
	return result
}

func parseRecursive(arg []interface{}, list []interface{}) []interface{}{
	for _, value := range list {
		item := value.(map[string]interface{})
		arg = append(arg, item["field"])
		kind := reflect.TypeOf(item["value"]).Kind()
		if kind == reflect.String {
			arg = append(arg, item["value"])
		}else if kind == reflect.Slice || kind == reflect.Array {
			return parseRecursive(arg, item["value"].([]interface{}))
		}
	}
	return arg
}

// GetLogger is returning the job logger
func (j Job) GetLogger() *logging.Logger {
	return j.log
}

