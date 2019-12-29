package account

import (
	"bou.ke/monkey"
	"crypto/tls"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestAccount(t *testing.T) {
	t.Run("DialAndLoginPromise", func(t *testing.T) {
		t.Run("OK", dialAndLoginPromiseOk)
		t.Run("WithDebug", dialAndLoginPromiseWithDebug)
		t.Run("Connection", func(t *testing.T) {
			t.Run("Failed", dialAndLoginPromiseConnectionFailed)
		})
		t.Run("Login", func(t *testing.T) {
			t.Run("Failed", dialAndLoginPromiseLoginFailed)
		})
		t.Run("Logout", func(t *testing.T) {
			t.Run("Failed", dialAndLoginPromiseLogoutFailed)
		})
	})
}

func dialAndLoginPromiseOk(t *testing.T) {
	defer SetClientFactory(nil)

	logOut := false
	SetClientFactory(func(addr string, tlsConfig *tls.Config) (IClient, error) {
		require.Equal(t, "localhost:1000", addr)
		require.NotNil(t, tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		client := NewMockClient()
		client.LoginCallback = func(username, password string) error {
			require.Equal(t, "foo", username)
			require.Equal(t, "bar", password)
			return nil
		}
		client.SetDebugCallback = func(w io.Writer) {
			require.Fail(t, "debug is not enabled")
		}
		client.LogoutCallback = func() error {
			logOut = true
			return nil
		}
		return client, nil
	})

	acc := Account{"foo_bar", "foo", "bar", "localhost", 1000, false, true}
	acc.DialAndLoginPromise(func(promise *ImapPromise) {
		require.NotNil(t, promise)
	})

	require.True(t, logOut)
}

func dialAndLoginPromiseWithDebug(t *testing.T) {
	defer monkey.UnpatchAll()
	defer SetClientFactory(nil)
	logOut, useDebug := false, false
	SetClientFactory(func(addr string, tlsConfig *tls.Config) (IClient, error) {
		require.Equal(t, "localhost:1000", addr)
		require.NotNil(t, tlsConfig)
		require.False(t, tlsConfig.InsecureSkipVerify)

		client := NewMockClient()
		client.LoginCallback = func(username, password string) error {
			require.Equal(t, "foo", username)
			require.Equal(t, "bar", password)
			return nil
		}
		client.SetDebugCallback = func(w io.Writer) {
			useDebug = true
		}
		client.LogoutCallback = func() error {
			logOut = true
			return nil
		}
		return client, nil
	})

	acc := Account{"foo_bar", "foo", "bar", "localhost", 1000, true, false}
	acc.DialAndLoginPromise(func(promise *ImapPromise) {
		require.NotNil(t, promise)
	})

	require.True(t, useDebug)
	require.True(t, logOut)
}

func dialAndLoginPromiseConnectionFailed(t *testing.T) {
	defer SetClientFactory(nil)

	SetClientFactory(func(addr string, tlsConfig *tls.Config) (IClient, error) {
		require.Equal(t, "localhost:1000", addr)
		require.NotNil(t, tlsConfig)
		require.False(t, tlsConfig.InsecureSkipVerify)

		return nil, errors.New("panic")
	})

	defer func() {
		err := recover()
		require.NotNil(t, err, "no panic catched")
	}()

	acc := Account{"foo_bar", "foo", "bar", "localhost", 1000, false, false}
	acc.DialAndLoginPromise(func(promise *ImapPromise) {
		require.Fail(t, "never called")
	})
}

func dialAndLoginPromiseLoginFailed(t *testing.T) {
	defer SetClientFactory(nil)
	SetClientFactory(func(addr string, tlsConfig *tls.Config) (IClient, error) {
		require.Equal(t, "localhost:1000", addr)
		require.NotNil(t, tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		client := NewMockClient()
		client.LoginCallback = func(username, password string) error {
			require.Equal(t, "foo", username)
			require.Equal(t, "bar", password)
			return errors.New("panic")
		}
		return client, nil
	})

	defer func() {
		err := recover()
		require.NotNil(t, err, "no panic catched")
	}()

	acc := Account{"foo_bar", "foo", "bar", "localhost", 1000, false, true}
	acc.DialAndLoginPromise(func(promise *ImapPromise) {
		require.Fail(t, "never called")
	})
}

func dialAndLoginPromiseLogoutFailed(t *testing.T) {
	defer SetClientFactory(nil)
	SetClientFactory(func(addr string, tlsConfig *tls.Config) (IClient, error) {
		require.Equal(t, "localhost:1000", addr)
		require.NotNil(t, tlsConfig)
		require.True(t, tlsConfig.InsecureSkipVerify)

		client := NewMockClient()
		client.LoginCallback = func(username, password string) error {
			require.Equal(t, "foo", username)
			require.Equal(t, "bar", password)
			return nil
		}
		client.LogoutCallback = func() error {
			return errors.New("panic")
		}
		return client, nil
	})

	defer func() {
		err := recover()
		require.Nil(t, err, "panic catched")
	}()

	acc := Account{"foo_bar", "foo", "bar", "localhost", 1000, false, true}
	acc.DialAndLoginPromise(func(promise *ImapPromise) {
		require.NotNil(t, promise)
	})
}

func TestNewAccountForTest(t *testing.T) {
	a:=NewAccountForTest(t, "","","","",true)
	require.Empty(t, a.name)
	require.Empty(t, a.username)
	require.Empty(t, a.password)
	require.Empty(t, a.hostname)
	require.False(t, a.debug)
	require.NotEmpty(t, a.skipVerify)

	defer func() {
		err := recover()
		require.NotNil(t,err)
		require.EqualError(t,err.(error),"invalid caller")
	}()
	NewAccountForTest(nil, "","","","",true)
}
