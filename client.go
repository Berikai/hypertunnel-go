package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/berikai/hypertunnel-go/relay"
)

type Client struct {
	Port                int
	Host                string
	Server              string
	ServerParts         *url.URL
	Token               string
	DesiredInternetPort int
	Options             *Options

	Deleted      bool
	Relay        *relay.RelayClient
	InternetPort int
	RelayPort    int
	URI          string
	Secret       string
	CreatedAt    string
	ExpiresIn    int
	ServerBanner string
}

type ClientOptions struct {
	Host         string
	Server       string
	Token        string
	InternetPort int
}

type Options struct {
	SSL bool
}

func NewClient(port int, opts *ClientOptions, options *Options) (*Client, error) {

	if opts == nil {
		opts = &ClientOptions{}
	}
	if options == nil {
		options = &Options{}
	}

	serverParts, err := url.Parse(opts.Server)
	if err != nil {
		panic(err)
	}

	client := &Client{
		Port:                port,
		Host:                opts.Host,
		Server:              opts.Server,
		ServerParts:         serverParts,
		Token:               opts.Token,
		DesiredInternetPort: opts.InternetPort,
		Options:             options,
	}

	if client.Host == "" {
		client.Host = "localhost"
	}
	if client.Server == "" {
		client.Server = "https://hypertunnel.ga"
	}
	if client.Token == "" {
		client.Token = "free-server-please-be-nice"
	}
	return client, nil
}

func (c *Client) Create() error {
	payload := map[string]interface{}{
		"serverToken":  c.Token,
		"internetPort": c.DesiredInternetPort,
		"ssl":          c.Options.SSL,
	}
	payload_json, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println("create payload", payload)
	}

	res, err := http.Post(c.Server+"/create", "application/json", bytes.NewBuffer(payload_json))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	if debug {
		fmt.Println("create", body)
	}
	if !body["success"].(bool) {
		return errors.New(body["message"].(string))
	}
	c.CreatedAt = body["createdAt"].(string)
	c.InternetPort = int(body["internetPort"].(float64))
	c.RelayPort = int(body["relayPort"].(float64))
	c.Secret = body["secret"].(string)
	c.URI = body["uri"].(string)
	c.ExpiresIn = int(body["expiresIn"].(float64))
	c.ServerBanner = body["serverBanner"].(string)

	c.Relay = relay.NewClient(&relay.ClientOptions{
		Host:      c.Host,
		Port:      c.Port,
		RelayHost: c.ServerParts.Hostname(),
		RelayPort: c.RelayPort,
	}, &relay.RelayOptions{
		Secret: c.Secret,
	})
	return nil
}

func (c *Client) Delete() error {
	if c.Deleted {
		return nil
	}
	payload := map[string]interface{}{
		"serverToken":  c.Token,
		"internetPort": c.InternetPort,
		"secret":       c.Secret,
	}
	payload_json, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(c.Server+"/delete", "application/json", bytes.NewBuffer(payload_json))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	if !body["success"].(bool) {
		return errors.New(body["message"].(string))
	}
	c.Deleted = true
	return nil
}

func (c *Client) Close() error {
	return c.Delete()
}
