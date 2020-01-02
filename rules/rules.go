package rules

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"mailAssistant/cntl"
	"mailAssistant/logging"
	"os"
//	"time"

	"io/ioutil"
	"mailAssistant/account"
	"path/filepath"
	"strings"
)

var (
//	never       time.Time
	logRules    *logging.Logger
	rulesDir,rulesDirErr = filepath.Abs("resources/config/rules/")
	rulesWalker func(*fsnotify.Watcher) filepath.WalkFunc
)

func init() {
//	never, _ = time.Parse(time.RFC3339, "2999-12-31T23:59:59Z")
	setRulesWalker(nil)
}

func setRulesWalker(walker func(*fsnotify.Watcher) filepath.WalkFunc) {
	if walker == nil {
		rulesWalker = rulesWatchDir
	} else {
		rulesWalker = walker
	}
}

// ImportAndLaunch is reading rule yaml and launching all jobs
func ImportAndLaunch(accounts *account.Accounts) error {
	if rulesDirErr == nil {
		rules := newRules(accounts)
		rules.loadFromDisk(rules.rulesDir)
		go rules.startWatcher()
	}
	return rulesDirErr
} // ImportRules

func newRules(accounts *account.Accounts) Rules {
	if !strings.HasSuffix(rulesDir, string(filepath.Separator)) {
		rulesDir += string(filepath.Separator)
	}
	return Rules{make(map[string]string, 0), make(map[string]Rule, 0), accounts, rulesDir}
}

// Rules represents a collation of rule yaml's
type Rules struct {
	files    map[string]string
	Rules    map[string]Rule
	removed  map[string] bool
	accounts *account.Accounts
	rulesDir string
}

func (rules *Rules) startWatcher() {
	watcher, _ := fsnotify.NewWatcher()
	defer func(){
		if err := recover(); err == nil {
			_ = watcher.Close()
		}else{
			panic(err)
		}
	}()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasSuffix(event.Name, "~") {
					continue
				}
				rules.getLogger().Debug(event.Name)
				rules.importRule("", event.Name, event.Op)
			case err := <-watcher.Errors:
				rules.getLogger().Severe("watcher - ERROR ", err)
			} // select
		} // for
	}() // go func

	if err := filepath.Walk(rules.rulesDir, rulesWalker(watcher)); err != nil {
		rules.getLogger().Panic("startWalker", err)
	}

	cntl.WaitForNotify()
}

func (rules *Rules) loadFromDisk(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			rules.loadFromDisk(filepath.Join(path, file.Name()))
			continue
		}
		rules.getLogger().Debug("foreach %s/%s\n", path, file.Name())
		rules.importRule(path, file.Name(), fsnotify.Create)
	} // for files
}

func (rules Rules) importRule(path, file string, op fsnotify.Op) {
	if op == fsnotify.Create {
		if rule, err := parseYaml(rules.rulesDir, path, file); err != nil {
			rules.getLogger().Panic("ERROR ", err)
		} else if rule == nil || rule.IsEmpty() {
			// nothing to do
		} else {
			rules.getLogger().Debug("Create --> %s => %s \n", rule.Name, rule.fileName)
			r := rule.convert()
			rules.Rules[rule.Name] = r

			log.Print(">>> ",rule.fileName)

			if _, found := rules.files[rule.fileName]; found && ! rules.removed[rule.fileName]{
				rules.getLogger().Panic(rule.fileName)
			}
			delete(rules.removed, rule.fileName)
			rules.files[rule.fileName] = rule.Name
			r.Schedule(rules.accounts)
		}
	} else if op == fsnotify.Write {
		rules.importRule(path, file, fsnotify.Remove)
		rules.importRule(path, file, fsnotify.Create)
	} else if op == fsnotify.Remove {
		fileName := ""
		if path == "" {
			fileName = strings.ToLower(file)
		}else{
			fileName = strings.ToLower(filepath.Join(path, file))
		}
		rules.getLogger().Debug("Remove ", path,"=>", file,"->", fileName)
		ruleFileName := strings.TrimPrefix(fileName, strings.ToLower(rules.rulesDir))
		rules.getLogger().Debug("Remove <<<", ruleFileName)
		rules.removed[ruleFileName] = true
		ruleName := rules.files[ruleFileName]
		rules.Rules[ruleName].Stop()
		delete(rules.Rules, ruleName)
		delete(rules.files, fileName)
	} // Remove
}

func (rules Rules) getLogger() *logging.Logger {
	if logRules == nil {
		logRules = logging.NewNamedLogger("${project}.rule.rules")
	}
	return logRules
}

func rulesWatchDir(watcher *fsnotify.Watcher) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}else if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	}
}