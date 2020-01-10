package conditions

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCase(t *testing.T){
	client, _ := client.DialTLS("imap.udag.de:993",nil)
	_ = client.Login("mahillmannde-0001","Unlock#6074")
	defer client.Logout()
/*
	boxes := make(chan *imap.MailboxInfo)
	go client.List("","*Mare*", boxes)

	for box := range boxes {
		log.Printf("'%s' '%s' -> %#v", box.Name, box.Delimiter, box.Attributes)
	}
	return */

	client.SetDebug(os.Stderr)
	item := make(map[interface{}]interface{})
	content, err := ioutil.ReadFile("b:/mailAssistant/config/rules/archiv/AWS.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(content, &item); err != nil{
		panic(err)
	}

	findArgs := item["args"].([]interface{})
	var searchArg interface{}

	for i := range findArgs {
		findArg, ok := findArgs[i].(map[interface{}]interface{})
		if ok {
			if v, f := findArg["search"]; f {
				searchArg = v
				break
			}
		}
	}

	cond := ParseYaml(searchArg)
	args := cond.Get()
	search := imap.NewSearchCriteria()
	if xerr := search.ParseWithCharset(args,nil); xerr != nil {
		panic(xerr)
	}

	_,_ =  client.Select("INBOX",true)
	seqNum, err := client.Search(search)
	log.Print(err)
	log.Printf("%#v", seqNum)

	s := new(imap.SeqSet)
	s.AddNum(seqNum...)
	if len(s.Set) == 0 {
		return
	}
	msg := make(chan *imap.Message)
	go client.Fetch(s,[]imap.FetchItem{imap.FetchFlags, imap.FetchEnvelope}, msg)

	for m := range msg {
		log.Printf("%d: %s\t%#v", m.SeqNum, m.Envelope.Subject, m.Flags)
	}
}