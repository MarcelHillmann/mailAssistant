package testDriver

import (
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mailAssistant/conditions"
	"os"
	"strings"
)

func TestTreiber(c *cli.Context) error {
	username, password, server := c.String("username"), c.String("password"), c.String("server")
	file := c.Path("file")
	verbose,sVerbose := c.Bool("verbose"), c.Bool("sverbose")

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
				go c.Fetch(s, []imap.FetchItem{
					imap.FetchBody,
					imap.FetchBodyStructure,
					imap.FetchEnvelope,
					imap.FetchFlags,
					imap.FetchInternalDate,
					imap.FetchRFC822Header,
					imap.FetchRFC822Size,
					imap.FetchRFC822Text,
					imap.FetchRFC822}, msg)

				for m := range msg {
					env := m.Envelope
					line := fmt.Sprintf("%s\t", env.Subject)
					for _, addr := range env.From { line += fmt.Sprintf("F: %s ", ToString(addr)) }
					for _, addr := range env.To { line += fmt.Sprintf("T: %s ", ToString(addr)) }
					log.Print(line)

					bodyStructure := m.BodyStructure
					log.Printf("%d\n%d\n%d\n%s\n%s\n%s\n%t\n%s\n%s\n%#v\n%s\n%s\n%s\n%s\n%d\n%#v\n%#v",
						m.SeqNum,m.Size,m.Uid, bodyStructure.Description,
						bodyStructure.Disposition, bodyStructure.Encoding, bodyStructure.Extended,
						bodyStructure.Id, bodyStructure.Language, bodyStructure.Lines, bodyStructure.Location,
						bodyStructure.MD5, bodyStructure.MIMESubType, bodyStructure.MIMEType, bodyStructure.Size,
						bodyStructure.DispositionParams, bodyStructure.Params)

					if sVerbose {
						superVerbose(m)
					}
				}
			}
		}
	}

	return nil
}


func superVerbose(m *imap.Message) {
	section := new(imap.BodySectionName)
	for true {
		literal := m.GetBody(section)
		if reader, err := mail.CreateReader(literal); err != nil {
			break
		}else{
			for true{
				part, err := reader.NextPart()
				if err != nil {
					break
				}
				switch part.Header.(type){
				case *mail.AttachmentHeader:
					att := part.Header.(*mail.AttachmentHeader)
					fields := att.Fields()
					for fields.Next() {
						log.Printf("%s\t%s", fields.Key(), fields.Value())
					}
				case *mail.InlineHeader:
					inline := part.Header.(*mail.InlineHeader)
					fields := inline.Fields()
					for fields.Next() {
						log.Printf("%s\t%s", fields.Key(), fields.Value())
					}
				case *mail.Header:
					h := part.Header.(*mail.Header)
					fields := h.Fields()
					for fields.Next() {
						log.Printf("%s\t%s", fields.Key(), fields.Value())
					}
				}
				if buffer, err := ioutil.ReadAll(part.Body); err == nil {
					content := string(buffer)

					for true {
						start1 := strings.Index(content, "Nura")
						start2 := strings.Index(content, "nura")
						start3 := strings.Index(content, "NURA")

						start := 0
						if start1 == -1 && start2 == -1 && start3 == -1{
							break
						} else if start1 > -1 {
							start = start1
						} else if start2 > -1{
							start = start2
						}else {
							start = start3
						}

						log.Printf("%#v", part.Header)
						x := content[start-10:]
						if len(x) > 20 {
							x = x[0:20]
						}
						log.Print(x)
						content = content[start+4:]
					}
				}

			}

		}
	}
}