package account

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"mailAssistant/logging"
	"os"
	"path/filepath"
	"strings"
)

var parserReadAll func(r io.Reader) ([]byte, error) = ioutil.ReadAll

func parseYaml(path string, file string) (*accountAux, error) {
	joinPath := filepath.Join(path, file)
	if path == "" {
		joinPath = file
	}
	filename := strings.ToLower(file)
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		osFile, err := os.Open(joinPath)
		if err == nil {
			defer func() { _ = osFile.Close() }()
			content, err := parserReadAll(osFile)
			if err == nil {
				account := &accountAux{fileName: file}
				err := yaml.Unmarshal(content, account)
				return account, err
			}
			logging.NewLogger().Panic(err)
			return nil, nil
		}
		logging.NewLogger().Severe(err)
		return nil, err
	}
	return nil, nil
}
