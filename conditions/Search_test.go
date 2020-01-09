package conditions

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
	"testing"
)

func TestCase(t *testing.T){
	client, _ := client.DialTLS("imap.udag.de:993",nil)
	_ = client.Login("mahillmannde-0001","Unlock#6074")
	defer client.Logout()

	boxes := make(chan *imap.MailboxInfo)
	go client.List("","*Mare*", boxes)

	for box := range boxes {
		log.Printf("'%s' '%s' -> %#v", box.Name, box.Delimiter, box.Attributes)
	}

	return

	client.SetDebug(os.Stderr)
	item := make([]interface{},0)
	/*	item = append(item, map[interface{}]interface{}{"field":"older","value":"24h"})
		item = append(item, map[interface{}]interface{}{"field":"unseen"})
		item = append(item, map[interface{}]interface{}{"field":"from","value":"mailer@dzone.com"})
		item = append(item, map[interface{}]interface{}{"field":"or","value":[]interface{}{
			map[interface{}]interface{}{"field":"from","value":"noreply@dzone.com"},
			map[interface{}]interface{}{"field":"from","value":"privacy@dzone.com"},
		}})
	*/
	item = append(item, map[interface{}]interface{}{"field":"before","value":"11-Dec-2019"})
	item = append(item, map[interface{}]interface{}{"field":"flagged"})
	item = append(item, map[interface{}]interface{}{"field":"from","value":"notebooksbilliger.de"})
	/*
		item = append(item, map[interface{}]interface{}{"field":"or","value":[]interface{}{
			map[interface{}]interface{}{"field":"from","value":"team@notebooksbilliger.de"},
			map[interface{}]interface{}{"field":"from","value":"team@notebooksbilliger.de"},
			map[interface{}]interface{}{"field":"from","value":"produktempfehlung@notebooksbilliger.de"},
			map[interface{}]interface{}{"field":"from","value":"service@notebooksbilliger.de"},
		}})
	*/
	cond := ParseYaml(item)
	args := cond.Get()
	search := imap.NewSearchCriteria()
	search.ParseWithCharset(args,nil)

	_,_ =  client.Select("INBOX.TEST",true)
	seqNum, err := client.Search(search)
	log.Print(err)
	log.Printf("%#v", seqNum)

	s := new(imap.SeqSet)
	s.AddNum(seqNum...)
	msg := make(chan *imap.Message)
	go client.Fetch(s,[]imap.FetchItem{imap.FetchFlags, imap.FetchEnvelope}, msg)

	for m := range msg {
		log.Printf("%d: %s\t%#v", m.SeqNum, m.Envelope.Subject, m.Flags)
	}
}