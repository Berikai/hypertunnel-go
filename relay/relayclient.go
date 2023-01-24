package relay

import (
	"time"
)

func (c *RelayClient) createClient() *RelayClient {
	client := &RelayClient{opts: c.opts, options: c.options}
	//client.on("pair", c.onClientPair)
	//client.on("close", c.onClientClose)
	//client.on("bytes", c.onClientBytes)
	return client
}

func (c *RelayClient) onClientPair() {
	c.pairCount += 1
	c.client = c.createClient()
}

func (c *RelayClient) onClientClose() {
	c.client = nil

	if c.opts.Retry {
		time.AfterFunc(5*time.Second, func() {
			if c.endCalled {
				return
			}
			c.client = c.createClient()
		})
	}
}

func (c *RelayClient) onClientBytes(tx, rx int) {
	c.bytes["tx"] += tx
	c.bytes["rx"] += rx
}

func (c *RelayClient) End() {
	c.endCalled = true
	c.client.relaySocket.Close()
}
