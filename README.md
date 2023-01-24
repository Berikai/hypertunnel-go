# ✨  hypertunnel-client, rewritten in go.

<a href="https://github.com/berstend/hypertunnel"><img src="https://i.stack.imgur.com/MN8RF.gif" width="280px" height="230px" align="right" /></a>

> When localtunnel/ngrok is not enough.

You can visit the main hypertunnel repo by clicking [here](https://github.com/berstend/hypertunnel). 

Keep in mind that **free public server** is not working anymore!

Hypertunnel is a free TCP relay/reverse proxy service can be used to **expose any TCP/IP service** running behind a NAT. It's using [hypertunnel-tcp-relay](/packages/hypertunnel-tcp-relay) under the hood, which itself is based on the excellent [node-tcp-relay](https://github.com/tewarid/node-tcp-relay) from [tewarid](https://github.com/tewarid), adding self-service multi-client support similar to localtunnel, a cool project name with "hyper" in it. This repository contains a hypertunnel client written in go.

## Status

**It's still in development and currently not working!**

It can connect to a hypertunnel server and get the following data.

        ✨  Hypertunnel created.
        Tunneling hypertunnel.ga:35209 > localhost:54323

        Hit ctrl+c to close the tunnel (maximum tunnel age is 1 day).

But the program ends after that. The later process is not made yet.

## Usage
```bash
❯❯❯ hypertunnel-go -help

  -internet-port int
        the desired internet port on the public server
  -localhost string
        local server (default "localhost")
  -port int
        local TCP/IP service port to tunnel
  -server string
        hypertunnel server to use (default "https://hypertunnel.ga")
  -ssl
        enable SSL termination (https://) on the public server
  -token string
        token required by the server (default "free-server-please-be-nice")
```


## Contributing

Contributions are welcome.

## Related

- [localtunnel](https://github.com/localtunnel/localtunnel)
- [ngrok](https://ngrok.com/)
- [serveo](https://serveo.net/)
- [telebit.js](https://git.coolaj86.com/coolaj86/telebit.js)

## License

MIT
