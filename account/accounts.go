package account

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"mailAssistant/cntl"
	"mailAssistant/logging"
	"os"
	"path/filepath"
	"strings"
)

var (
	importPath    string
	importError   error
	accountWalker func(*fsnotify.Watcher) filepath.WalkFunc
)

func init() {
	SetAccountImportPath("")
	SetAccountWalker(nil)
}

// ImportAccounts reads the account Yaml files and register file watcher
func ImportAccounts() (*Accounts, error) {
	accounts := newAccounts()
	accounts.accountDir = importPath
	if err := accounts.loadFromDisk(); err == nil && importError == nil {
		go accounts.startWatcher()
		return &accounts, nil
	} else if importError != nil {
		return nil, importError
	} else {
		return nil, err
	}
}

func newAccounts() Accounts {
	return Accounts{Account: make(map[string]Account), files: make(map[string]string)}
}

// Accounts represents the collection of account Yaml's
type Accounts struct {
	accountDir string
	Account    map[string]Account
	files      map[string]string
}

// GetAccount search in memory the named account
// @Nillable
func (a Accounts) GetAccount(name string) *Account {
	if value, found := a.Account[name]; found {
		return &value
	}
	return nil
}
// HasAccount is checking if a named account exists
func (a Accounts) HasAccount(name string) bool {
	if _, found := a.Account[name]; found {
		return true
	}
	return false
}

func (a *Accounts) loadFromDisk() error {
	files, err := ioutil.ReadDir(a.accountDir)
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			a.importAccount(a.accountDir, file.Name(), fsnotify.Create)
		}
		return nil
	}
	return err
}

func (a Accounts) importAccount(path string, filename string, op fsnotify.Op) {
	if op == fsnotify.Create {
		account, err := parseYaml(path, filename)
		if err != nil {
			logging.NewLogger().Panic("ERROR", err)
		} else if account == nil || account.IsEmpty() {
			return
		}

		a.Account[account.Name] = account.convert()
		a.files[account.fileName] = account.Name
	} else if op == fsnotify.Write {
		a.importAccount(path, filename, fsnotify.Remove)
		a.importAccount(path, filename, fsnotify.Create)
	} else if op == fsnotify.Remove {
		// log.Println("Remove ", filename)
		fileName := strings.ToLower(filename)
		accName := a.files[fileName]
		delete(a.Account, accName)
		delete(a.files, fileName)
	} // Remove
}

func (a *Accounts) startWatcher() {
	watcher, _ := fsnotify.NewWatcher()
	defer func(){ _ = watcher.Close()}()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				a.importAccount("", event.Name, event.Op)
			case err := <-watcher.Errors:
				if err != nil {
					logging.NewLogger().Severe("go func watcher - ERROR", err)
				}
			} // select
		} // for
	}() // go func

	if err := filepath.Walk(a.accountDir, accountWalker(watcher)); err != nil {
		logging.NewLogger().Panic("startWalker - ERROR", err)
	} else {
		cntl.WaitForNotify()
	}
}

// SetAccountWalker is used for unit tests
func SetAccountWalker(walker func(*fsnotify.Watcher) filepath.WalkFunc) {
	if walker == nil {
		accountWalker = accountWatchDir
	} else {
		accountWalker = walker
	}
}

// SetAccountImportPath is used for unit tests
func SetAccountImportPath(path string) {
	if path != "" {
		importPath, importError = path, nil
	} else {
		importPath, importError = filepath.Rel(".", "resources/config/accounts/")
	}
}

func accountWatchDir(watcher *fsnotify.Watcher) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			_ = watcher.Add(path)
		}
		return nil
	}
}
