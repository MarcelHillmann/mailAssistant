package cmd

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/commands"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mailAssistant/conditions"
	"os"
	"path/filepath"
	"strings"
)

func RunConfigCheck(c *cli.Context) error {
	log.SetFlags(log.Ltime)
	log.Println("Run config check")

	configDir := c.String("config")
	return runRecursive(configDir, configDir)
}

func runRecursive(base, dir string) error {
	if file, err := os.OpenFile(dir, os.O_RDONLY, 0); err == nil || os.IsExist(err) {
		file.Close()
		return runFile("", dir)
	}
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory not exists %s\n\t%s", dir, err)
	} else if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := runRecursive(base, path); err != nil {
				return err
			}
		} else if err = runFile(base, file.Name()); err != nil {
			return err
		}
	}

	return nil
}

func runFile(base, file string) error {
	if content, err := ioutil.ReadFile(file); err != nil {
		return err
	} else {
		log.Println("----- ----- ----- ----- ----- ----- ----- ----- ----- ----- ----- ")
		var condition conditions.Condition
		var yamlContent map[string]interface{}

		err := yaml.Unmarshal(content, &yamlContent)
		if err != nil {
			return err
		}
		args := yamlContent["args"].([]interface{})
		for _, arg := range args {
			item := arg.(map[interface{}]interface{})
			if search, ok := item["search"]; ok {
				condition = conditions.ParseYaml(search)
				break
			}
		}

		cnf := strings.TrimPrefix(base, base)
		if condition == nil || condition.String() == "" {
			return fmt.Errorf("check config file '%s'", cnf)
		}
		log.Printf("%s: %s\n", cnf, condition.String())
		criteria := imap.NewSearchCriteria()
		_ = criteria.ParseWithCharset(condition.Get(), nil)
		s := commands.Search{Charset: "UTF-8", Criteria: criteria}
		cmd := s.Command()
		_ = cmd.WriteTo(&imap.Writer{Writer: os.Stderr, AllowAsyncLiterals: false})
		log.Println("+++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++")
		return nil
	}
}
