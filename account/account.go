package account

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/emersion/go-imap/client"
	"mailAssistant/logging"
	"mailAssistant/utils"
	"testing"
)

type factory func(addr string, tlsConfig *tls.Config) (IClient, error)

var (
	clientTLS factory
)

// NewAccountForTest is used for unit tests
func NewAccountForTest(t *testing.T, name, user, pw, host string, skipVerify bool) Account {
	if t == nil {
		panic(errors.New("invalid caller"))
	}
	return Account{name, user, pw, host, 20000, false, skipVerify}
}

func init() {
	SetClientFactory(nil)
}

// SetClientFactory is for unit tests
func SetClientFactory(function factory) {
	if function == nil {
		clientTLS = func(addr string, tlsConfig *tls.Config) (c IClient, err error) {
			orgClient, orgErr := client.DialTLS(addr, tlsConfig)
			return NewClientPromise(orgClient), orgErr
		}
	} else {
		clientTLS = function
	}
}

func tlsConfig(insecure bool) *tls.Config {
	cnf := tls.Config{}
	cnf.InsecureSkipVerify = insecure
	return &cnf
}

// Account represents a mail account
type Account struct {
	name       string
	username   string
	password   string
	hostname   string
	port       int
	debug      bool
	skipVerify bool
}

// DialAndLoginPromise connect to IMAP server and if successfully call the callback
func (account Account) DialAndLoginPromise(callback func(*ImapPromise)) {
	connection, err := clientTLS(account.hostnamePort(), tlsConfig(account.skipVerify))
	if err != nil {
		panic(err)
	}

	if account.debug {
		connection.SetDebug(logging.NewNamedLogWriter("${project}.account." + account.name))
	}
	defer utils.Defer(connection.Logout)
	if err := connection.Login(account.username, account.password); err != nil {
		panic(err)
	}
	callback(newImapPromise(connection))
}

func (account Account) hostnamePort() string {
	return fmt.Sprintf("%s:%d", account.hostname, account.port)
}
