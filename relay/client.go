package relay

import (
	"crypto/tls"
	"fmt"
	"net"
)

type RelayClient struct {
	opts      ClientOptions
	pairCount int
	bytes     map[string]int
	client    *RelayClient
	endCalled bool

	options       *RelayOptions
	serviceSocket net.Conn
	bufferData    bool
	buffer        []byte

	relaySocket net.Conn
}

type ClientOptions struct {
	Host      string
	Port      int
	RelayHost string
	RelayPort int
	Options   map[string]interface{}
	Retry     bool
}

type RelayOptions struct {
	Tls                bool
	RejectUnauthorized bool
	Secret             string
}

func NewClient(opts *ClientOptions, options *RelayOptions) *RelayClient {
	c := &RelayClient{
		opts:       *opts,
		options:    options,
		bufferData: true,
	}
	if options.Tls {
		c.relaySocket = c.createSecureRelaySocket()
	} else {
		c.relaySocket = c.createRelaySocket()
	}
	c.onRelaySocketConnect()
	return c
}

func (c *RelayClient) createRelaySocket() net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.opts.RelayHost, c.opts.RelayPort))
	if err != nil {
		return nil
	}
	return conn
}

func (c *RelayClient) createSecureRelaySocket() net.Conn {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", c.opts.RelayHost, c.opts.RelayPort), &tls.Config{
		InsecureSkipVerify: !c.options.RejectUnauthorized,
	})
	if err != nil {
		return nil
	}
	c.onRelaySocketConnect()
	return conn
}

func (c *RelayClient) onRelaySocketConnect() {
	c.authorize()
}

func (c *RelayClient) authorize() {
	if c.options.Secret != "" {
		c.relaySocket.Write([]byte(c.options.Secret))
	}
}
