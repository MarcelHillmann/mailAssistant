package rules

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/require"
	"mailAssistant/cntl"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const fileName = "..\\testdata\\rules\\foobar.yml"
const fileNameUx = "../testdata/rules/foobar.yml"
const fileNameSub = "..\\testdata\\rules\\subdir\\barfoo.yml"
const fileNameSubUx = "../testdata/rules/subdir/barfoo.yml"

func TestRules(t *testing.T) {
	t.Run("startWatcher", func(t *testing.T) {
		t.Run("OK", rulesStartWatcherOk)
		t.Run("Error", rulesStartWatcherError)
	})
	t.Run("loadFromDisk", rulesLoadFromDisk)
	t.Run("importRule", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			t.Run("OK", rulesImportRuleCreate)
			t.Run("Empty", rulesImportRuleCreateEmpty)
			t.Run("NotWellFormed", rulesImportRuleCreateNotWellFormed)
			t.Run("Again", rulesImportRuleCreateAgain)
			t.Run("NotExisting", rulesImportRuleCreateNotExisting)
		})
		t.Run("Modify", rulesImportRuleModify)
		t.Run("Delete", rulesImportRuleDelete)
	})
	t.Run("rulesWatchDir", func(t *testing.T) {
		t.Run("OK", rulesWatchDirOk)
		t.Run("Error", rulesWatchDirError)
		t.Run("IsDir", rulesWatchDirIsDir)
	})
}

func TestRules_GetLogger_Init(t *testing.T) {
	logRules = nil
	r := Rules{}
	log := r.getLogger()
	require.NotNil(t, log)
	require.Equal(t, "mailAssistant.rule.rules", log.Name())
	logRules = nil
}

func TestRules_GetLogger_ReReInit(t *testing.T) {
	logRules = nil
	r := Rules{}
	log := r.getLogger()
	require.NotNil(t, log)
	require.Equal(t, "mailAssistant.rule.rules", log.Name())
	log2 := r.getLogger()
	require.NotNil(t, log2)
	require.Equal(t, "mailAssistant.rule.rules", log2.Name())
	logRules = nil
}

func rulesImportRuleCreate(t *testing.T) {
	r := newRules(nil)
	require.Len(t, r.Rules, 0, "rules pre")
	require.Len(t, r.files, 0, "files pre")

	r.importRule("..", "testdata/rules/fooBar.yml", fsnotify.Create)
	require.Len(t, r.Rules, 1, "rules post")
	require.Len(t, r.files, 1, "files post")

	require.Condition(t, func() bool {
		return r.files[fileName] == "testcase" || r.files[fileNameUx] == "testcase"
	}, "r.[filename] == testcase")
	tc := r.Rules["testcase"]
	require.NotNil(t, tc)
	require.Equal(t, "foo.bar", tc.GetString("mail_account"))
	require.Equal(t, "foo@bar.local", tc.GetString("mail_from"))
	require.Equal(t, "INBOX", tc.GetString("path"))
	require.False(t, tc.HasArg("notExisting"))
	tc2 := r.Rules["testcase2"]
	require.NotNil(t, tc2)
	require.Nil(t, tc2.Args)
	require.Equal(t, "", tc2.name)
	require.Equal(t, "", tc2.schedule)
	require.Equal(t, "", tc2.action)
	require.Nil(t, tc2.clock)
	require.False(t, tc2.disabled)
}

func rulesImportRuleCreateEmpty(t *testing.T) {
	r := newRules(nil)
	require.Len(t, r.Rules, 0)
	require.Len(t, r.files, 0)

	r.importRule("..", "testdata/rules_invalid/empty.yml", fsnotify.Create)
	require.Len(t, r.Rules, 0)
	require.Len(t, r.files, 0)
}

func rulesImportRuleCreateNotWellFormed(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		require.EqualError(t, err.(error), "ERROR  yaml: line 7: could not find expected ':'")
	}()

	r := newRules(nil)
	require.Len(t, r.Rules, 0)
	require.Len(t, r.files, 0)

	r.importRule("..", "testdata/rules_invalid/error.yml", fsnotify.Create)
}

func rulesImportRuleCreateAgain(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		sErr := err.(error).Error()
		require.Condition(t, func() bool {
			return sErr == fileName || sErr == fileNameUx
		}, sErr)
	}()
	r := newRules(nil)
	r.files[fileName] = ""
	r.files[fileNameUx] = ""
	r.importRule("", "../testdata/rules/fooBar.yml", fsnotify.Create)
}

func rulesImportRuleCreateNotExisting(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		sErr := err.(error).Error()
		require.Condition(t, func() bool {
			return sErr == "open ..\\testdata\\rules\\dont_exist.yml: The system cannot find the file specified." ||
				sErr == "open ../testdata/rules/dont_exist.yml: no such file or directory"
		}, sErr)
	}()
	r := newRules(nil)
	r.files[fileName] = ""
	r.files[fileNameUx] = ""
	r.importRule("..", "testdata/rules/dont_exist.yml", fsnotify.Create)
}

func rulesImportRuleModify(t *testing.T) {
	r := newRules(nil)
	require.Len(t, r.Rules, 0)
	require.Len(t, r.files, 0)
	r.importRule("..", "testdata/rules/fooBar.yml", fsnotify.Write)
	require.Len(t, r.Rules, 1)
	require.Len(t, r.files, 1)
}

func rulesImportRuleDelete(t *testing.T) {
	r := newRules(nil)
	if strings.ToLower(runtime.GOOS) == "windows" {
		r.files[fileName] = "testcase"
	} else {
		r.files[fileNameUx] = "testcase"
	}
	r.Rules["testcase"] = Rule{}

	require.Len(t, r.Rules, 1)
	require.Len(t, r.files, 1)
	r.importRule("..", "testdata/rules/fooBar.yml", fsnotify.Remove)
	require.Len(t, r.Rules, 0)
	require.Len(t, r.files, 0)
}

func rulesLoadFromDisk(t *testing.T) {
	r := newRules(nil)
	r.loadFromDisk("../testdata/rules/")
	require.Len(t, r.Rules, 2)
	require.Len(t, r.files, 2)
	require.Condition(t, func() bool {
		return r.files[fileName] == "testcase" || r.files[fileNameUx] == "testcase"

	}, "r.files == testcase")
	require.Condition(t, func() bool {
		return r.files[fileNameSub] == "sub testcase" || r.files[fileNameSubUx] == "sub testcase"
	}, "sub testcase")
	require.NotNil(t, r.Rules["testcase"])
	require.NotNil(t, r.Rules["sub testcase"])
}

func rulesStartWatcherOk(t *testing.T) {
	defer setRulesWalker(nil)

	setRulesWalker(func(watcher *fsnotify.Watcher) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			watcher.Events <- fsnotify.Event{Name: "text.yaml~", Op: fsnotify.Create}
			_ = watcher.Close()
			return nil
		}
	})

	var started, stopped bool
	r := newRules(nil)
	go func() {
		started = true
		r.startWatcher()
		stopped = true
	}()

	time.Sleep(100 * time.Millisecond)
	require.Equal(t, 1, cntl.ToNotify(), "ToNotify 1")
	cntl.Notify()
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, 0, cntl.ToNotify(), "ToNotify 0")
	require.True(t, started, "started?")
	require.True(t, stopped, "stopped?")
}

func rulesStartWatcherError(t *testing.T) {
	defer func() {
		err := recover()
		require.EqualError(t, err.(error), "startWalker walker must fail")
		setRulesWalker(nil)
	}()

	setRulesWalker(func(watcher *fsnotify.Watcher) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			return errors.New("walker must fail")
		}
	})
	r := newRules(nil)
	r.startWatcher()
	require.Fail(t, "never call")
}

type MockFileInfo struct {
	dir bool
}

func (m MockFileInfo) Name() string {
	return ""
}
func (m MockFileInfo) Size() int64 {
	return 0
}
func (m MockFileInfo) Mode() os.FileMode {
	return 0
}
func (m MockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (m MockFileInfo) IsDir() bool {
	return m.dir
}
func (m MockFileInfo) Sys() interface{} {
	return nil
}
func rulesWatchDirOk(t *testing.T) {
	w, _ := fsnotify.NewWatcher()
	require.Nil(t, rulesWatchDir(w)("", MockFileInfo{false}, nil))
}

func rulesWatchDirError(t *testing.T) {
	w, _ := fsnotify.NewWatcher()
	require.EqualError(t, rulesWatchDir(w)("", MockFileInfo{false}, errors.New("fail")), "fail")
}

func rulesWatchDirIsDir(t *testing.T) {
	w, _ := fsnotify.NewWatcher()
	require.Nil(t, rulesWatchDir(w)("", MockFileInfo{true}, nil))
}
