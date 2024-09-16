# reverssh

Reversing SSH tarpit.

This tool makes SSH bots brute-force themselves. 
If port numbers are specified (e.g., `-p 22,2222`), 
`reverssh` will attempt connections on each of the listed ports of the SSH client. 
If any connection is successful, it forwards all incoming traffic back to the open port, 
causing bots to interact with their own servers.

If no ports are specified, or if all provided ports are closed, `reverssh` behaves as a standard SSH tarpit, 
sending one random byte per second.

## Installation

See [Releases](https://github.com/jackcvr/reverssh/releases)

## Usage

```shell
Usage of reverssh:
  -b string
    	Local address to listen on (default "0.0.0.0:22")
  -f string
    	Log file (default stdout)
  -l value
    	Log level. Possible values: debug, info, warn, error (default info)
  -p value
    	Remote ports to connect to, e.g. '22,2222'
  -q	Do not print anything (default false)
  -stats
    	Show active connections info
```

## Examples

Start reversing tarpit on 22 port (with redirect clients back to the same port):

```shell
$ sudo reverssh -b 0.0.0.0:22 -p 22
{"time":"2024-09-14T17:13:27.111626861+03:00","level":"INFO","msg":"listening","addr":"0.0.0.0:2222"}
{"time":"2024-09-14T17:13:32.080358768+03:00","level":"INFO","msg":"accepted","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:39680"}
{"time":"2024-09-14T17:13:32.08045588+03:00","level":"INFO","msg":"connected","laddr":"127.0.0.1:40136","raddr":"127.0.0.1:22"}
{"time":"2024-09-14T17:13:45.008896864+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:40136","raddr":"127.0.0.1:22"}
{"time":"2024-09-14T17:13:47.009419814+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:39680","lifetime":13,"reversed":true}
```

Start normal tarpit on 2222 port:

```shell
$ sudo reverssh -b 0.0.0.0:2222
{"time":"2024-09-14T17:15:01.726948856+03:00","level":"INFO","msg":"listening","addr":"0.0.0.0:2222"}
{"time":"2024-09-14T17:15:04.231376092+03:00","level":"INFO","msg":"accepted","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:58262"}
{"time":"2024-09-14T17:15:11.239589332+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:58262","lifetime":6,"reversed":false}
```

Show current activity:

```shell
$ sudo reverssh -stats
active connections:
127.0.0.1:48800 lifetime=4 reversed=true
127.0.0.1:57968 lifetime=23 reversed=true
```

## License

[MIT](https://spdx.org/licenses/MIT.html) 