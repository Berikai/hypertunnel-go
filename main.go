package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	debug = false
	bugs  = "https://github.com/berikai/hypertunnel-go/issues"
)

func main() {
	port := flag.Int("port", 0, "local TCP/IP service port to tunnel")
	localhost := flag.String("localhost", "localhost", "local server")
	server := flag.String("server", "https://hypertunnel.ga", "hypertunnel server to use")
	token := flag.String("token", "free-server-please-be-nice", "token required by the server")
	internetPort := flag.Int("internet-port", 0, "the desired internet port on the public server")
	ssl := flag.Bool("ssl", false, "enable SSL termination (https://) on the public server")
	flag.Parse()

	if *port == 0 {
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := NewClient(*port, &ClientOptions{
		Host:         *localhost,
		Server:       *server,
		Token:        *token,
		InternetPort: *internetPort,
	}, &Options{
		SSL: *ssl,
	})

	if err != nil {
		fmt.Printf("An unexpected error occurred. Please report this issue at: %s\n", bugs)
		fmt.Println(err)
		os.Exit(1)
	}

	if err := client.Create(); err != nil {
		fmt.Printf("An unexpected error occurred. Please report this issue at: %s\n", bugs)
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("  ✨  Hypertunnel created.\n  Tunneling %s > %s:%d\n", client.URI, client.Host, client.Port)

	if client.ServerBanner != "" {
		fmt.Println(client.ServerBanner)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)

	go func() {
		select {
		case <-sigint:
			if err := client.Close(); err != nil {
				fmt.Printf("An unexpected error occurred. Please report this issue at: %s\n", bugs)
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Tunnel closed.")
			os.Exit(0)
		case <-sigterm:
			if err := client.Close(); err != nil {
				fmt.Printf("An unexpected error occurred. Please report this issue at: %s\n", bugs)
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Tunnel closed.")
			os.Exit(0)
		}
	}()

	if client.ExpiresIn != 0 {
		time.AfterFunc(time.Duration(client.ExpiresIn)*time.Second, func() {
			if err := client.Close(); err != nil {
				fmt.Printf("An unexpected error occurred. Please report this issue at: %s\n", bugs)
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Tunnel expired.")
		})
	}
}
