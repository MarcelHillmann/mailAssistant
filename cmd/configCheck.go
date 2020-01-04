package cmd

import (
	"fmt"
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

func runRecursive(base , dir string) error {
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory not exists %s\n\t", dir, err)
	} else if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := runRecursive(base, path); err != nil {
				return err
			}
		} else if content, err := ioutil.ReadFile(path); err != nil {
			return err
		} else {
			log.Println("----- ----- ----- ----- ----- ----- ----- ----- ----- ----- ----- ")
			var condition conditions.Condition
			var yamlContent map[string]interface{}

			yaml.Unmarshal(content, &yamlContent)
			args := yamlContent["args"].([]interface{})
			for _ , arg := range args {
				item := arg.(map[interface{}]interface{})
				if search , ok := item["search"]; ok {
					condition = conditions.ParseYaml(search)
					break
				}
			}
			cnf := strings.TrimPrefix(path, base)
			if condition.String() == "" {
				return fmt.Errorf("check config file '%s'", cnf)
			}
			log.Printf("%s: %s\n", cnf, condition.String())
			log.Println("+++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++ +++++")
		}
	}

	return nil
}
