package testDriver

import (
	csv2 "encoding/csv"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/charset"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
)

func TestDriverAllMsg(c *cli.Context) error {
	username, password, server, path := c.String("username"), c.String("password"), c.String("server"), c.String("select")
	verbose := c.Bool("verbose")
	message.CharsetReader = charset.Reader
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

			if mbox, err := c.Select(path, true); err != nil {
				return err
			} else {
				s := new(imap.SeqSet)
				s.AddRange(1, mbox.Messages)

				msg := make(chan *imap.Message)
				go c.Fetch(s, []imap.FetchItem{imap.FetchEnvelope}, msg)

				csv := make([][]string, 0)
				csv = append(csv, []string{"num","SUBJECT", "Addr", "Mail"})
				for m := range msg {
					env := m.Envelope
					num := strconv.Itoa(int(m.SeqNum))
					subject :=  env.Subject
					for _, addr := range env.From { csv = append(csv, []string{num, subject, "FROM", ToString(addr)}) }
					for _, addr := range env.To { csv = append(csv, []string{num,subject, "TO", ToString(addr)}) }
					for _, addr := range env.Cc { csv = append(csv, []string{num,subject, "CC", ToString(addr)}) }
					for _, addr := range env.Bcc { csv = append(csv, []string{num,subject, "BCC", ToString(addr)}) }
				}

				out, _ := os.Create("D:/temp/overview.csv")
				defer out.Close()
				w := csv2.NewWriter(out)
				w.UseCRLF = true
				w.Comma = ';'
				w.WriteAll(csv)
				out.Close()
			}
		}
	}

	return nil
}

func ToString(addr *imap.Address) string {
	return fmt.Sprintf("%s %s <%s>", addr.MailboxName, addr.HostName, addr.PersonalName)
}
