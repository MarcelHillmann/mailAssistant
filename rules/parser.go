package rules

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"mailAssistant/logging"
	"mailAssistant/utils"
	"os"
	"path/filepath"
	"strings"
)

var parserReadAll func(r io.Reader) ([]byte, error) = ioutil.ReadAll

func parseYaml(rulesDir, path, file string) (*ruleAux, error) {
	joinPath := filepath.Join(path, file)
	if path == "" {
		joinPath = file
	}

	ruleFileName := strings.Replace(joinPath, rulesDir, "", -1)
	filename := strings.ToLower(file)
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		osFile, err := os.Open(joinPath)
		defer utils.Closer(osFile)

		if err == nil {
			content, err := parserReadAll(osFile)
			if err == nil {
				rule := &ruleAux{fileName: strings.ToLower(ruleFileName)}
				err := yaml.Unmarshal(content, rule)
				return rule, err
			}
			logging.NewLogger().Panic("ReadAll", err)
		}
		logging.NewLogger().Panic(err)
	}
	return nil, nil
}
