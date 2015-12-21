package client

import "github.com/getlantern/tlsdialer"

// Getfiretweetversion returns the current build version string
func GetFireTweetVersion() string {
	return defaultClient.getFireTweetVersion()
}

// GoCallback is the supertype of callbacks passed to Go
type GoCallback interface {
	Do()
    Protect(fd int)
}

// RunClientProxy creates a new client at the given address.
func RunClientProxy(listenAddr, appName string, ready GoCallback) error {
	go func() {
		defaultClient = newClient(listenAddr, appName)
		defaultClient.serveHTTP()
        tlsdialer.Callback = ready.Protect
		ready.Do()
	}()
	return nil
}

// StopClientProxy stops the proxy.
func StopClientProxy() error {
	defaultClient.stop()
	return nil
}
