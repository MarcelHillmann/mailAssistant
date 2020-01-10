package testDriver

import (
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mailAssistant/conditions"
	"os"
)

func TestTreiber(c *cli.Context) error {
	username, password, server := c.String("username"), c.String("password"), c.String("server")
	file := c.Path("file")
	verbose := c.Bool("verbose")

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	path := "INBOX"
	items := make(map[interface{}]interface{})
	sItems := make([]interface{},0)
	if err := yaml.Unmarshal(content, &items); err != nil {
		return err
	} else if val, ok := items["args"]; ok {
		for _, x := range val.([]interface{}) {
			y := x.(map[interface{}]interface{})
			if z, ok := y["search"]; ok {
				sItems = z.([]interface{})
				break
			} else if z, ok := y["path"]; ok {
				path = z.(string)
			}
		}
	}

	criteria := imap.NewSearchCriteria()
	search := conditions.ParseYaml(sItems)
	if c, err := client.DialTLS(server, nil); err != nil {
		return err
	} else {

		defer c.Close()
		if err := c.Login(username, password); err != nil {
			return err
		} else {
			defer c.Logout()
			if verbose {
				c.SetDebug(os.Stderr)
			}

			if _, err := c.Select(path, true); err != nil {
				return err
			} else if err := criteria.ParseWithCharset(search.Get(), nil); err != nil {
				return err
			} else if seqNum, err := c.Search(criteria); err != nil {
				return err
			}else if len(seqNum) == 0 {
				return errors.New("nothing found!!")
			} else {
				s := new(imap.SeqSet)
				s.AddNum(seqNum...)
				msg := make(chan *imap.Message)
				go c.Fetch(s, []imap.FetchItem{imap.FetchEnvelope}, msg)

				for m := range msg {
					env := m.Envelope
					line := fmt.Sprintf("%s\t", env.Subject)
					for _, addr := range env.From { line += fmt.Sprintf("F: %s ", ToString(addr)) }
					for _, addr := range env.To { line += fmt.Sprintf("T: %s ", ToString(addr)) }
					log.Print(line)
				}
			}
		}
	}

	return nil
}
