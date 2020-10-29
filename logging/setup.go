package logging

import (
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type logLevel int

// String shows current log level as String
func (l logLevel) String() string {
	switch l {
	case none:
		return "NONE"
	case all:
		return "ALL"
	case debug:
		return "DEBUG"
	case info:
		return "INFO"
	case warn:
		return "WARN"
	case severe:
		return "SEVERE"
	default:
		return "NOT EXISTS"
	}
}

var logCnf string

const (
	notExists logLevel = 0
	none               = notExists + 10
	all                = none + 10
	debug              = all + 10
	info               = debug + 10
	warn               = info + 10
	severe             = warn + 10
	global             = "global"
)

func init() {
	setLoggingPath("")
	loadLogging()
}

/*  SetLoggingPath is used for unit tests */
func setLoggingPath(targetPath string) {
	if targetPath == "" {
		logCnf, _ = filepath.Rel(".", "resources/config/logging.yml")
	} else {
		logCnf = targetPath
	}
}
func loadLogging() {
	importCnf(logCnf)
	go startWatcher(logCnf)
}

func startWatcher(filepath string) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(filepath)

	defer watcher.Close()
	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Write {
				importCnf(filepath)
			}
		case err := <-watcher.Errors:
			if err != nil {
				getLogger().Severe(err)
			}
		} // select
	} // for
}

func importCnf(filepath string) {
	config := &registryAux{}
	if osFile, err := os.Open(filepath); err != nil {
		getLogger().Severe(err)
	} else if err := unmarshal(osFile, config); err != nil {
		getLogger().Severe(err)
	} // os.Open
}

var setupReadAll func(r io.Reader) ([]byte, error) = ioutil.ReadAll

func unmarshal(file *os.File, config *registryAux) error {
	defer file.Close()
	if content, err := setupReadAll(file); err != nil {
		return err
	} else if err := yaml.Unmarshal(content, config); err != nil {
		return err
	}

	SetLevel("*", "")
	if config.HasNoLevel() {
		loggerRegistry[global] = none
	} else {
		loggerRegistry[global] = config.GetLevel()
	}
	if config.HasChildren() {
		childRecursive("", loggerRegistry, config.Children)
	}
	return nil
}

func childRecursive(path string, levels map[string]logLevel, children []*registryAux) {
	var _path = ""
	if path != "" {
		_path = path + "."
	}
	for _, child := range children {
		levels[_path+child.Name] = child.GetLevel()
		if child.HasChildren() {
			childRecursive(_path+child.Name, levels, child.Children)
		}
	}
}

func getLogger() Logger {
	return NewGlobalLogger()
}
