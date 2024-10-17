# reverssh

Reversing SSH Tarpit.

This tool makes SSH bots brute-force themselves. 
By using key `remoteports` in the config you can specify port numbers to instruct 
`reverssh` attempt connections on each of the specified ports of the SSH client. 
If any connection is successful, it forwards all incoming traffic back to the open port, 
causing bots to interact with their own servers.

If no ports are specified, or if all provided ports are closed, `reverssh` behaves as a standard SSH tarpit, 
sending one random byte per second.

Other features:
- linux packages(apk, deb) include services(openrc, systemd)
- JSON structured logs
- ability to bind to few addresses
- ability to watch active connections

## Installation

See [Releases](https://github.com/jackcvr/reverssh/releases)

## Usage

```shell
Usage of reverssh:
  -active
    	Show active connections info
  -c string
    	Path to TOML config file (default "/etc/reverssh/reverssh.toml")
```

### reverssh.toml sample
```toml
tz = "Europe/Vilnius"
verbose = false
quiet = false
bind = ["0.0.0.0:22"]
remoteports = [22]
```

## Examples

Start reversing tarpit on 2222 port (redirecting clients back to 22 port):

```shell
$ sudo reverssh -c reverssh.toml
{"time":"2024-09-18T15:17:08.854929365+03:00","level":"INFO","msg":"listening","addr":"0.0.0.0:2222"}
{"time":"2024-09-18T15:17:08.854953224+03:00","level":"INFO","msg":"listening","addr":"/var/run/reverssh.sock"}
{"time":"2024-09-18T15:17:13.053926647+03:00","level":"INFO","msg":"accepted","laddr":{"IP":"127.0.0.1","Port":2222,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":60988,"Zone":""}}
{"time":"2024-09-18T15:17:13.054203917+03:00","level":"INFO","msg":"connected","laddr":{"IP":"127.0.0.1","Port":44896,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
{"time":"2024-09-18T15:17:15.618838555+03:00","level":"INFO","msg":"accepted","laddr":{"IP":"127.0.0.1","Port":2223,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":60370,"Zone":""}}
{"time":"2024-09-18T15:17:15.618962245+03:00","level":"INFO","msg":"connected","laddr":{"IP":"127.0.0.1","Port":44908,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
{"time":"2024-09-18T15:17:18.844756922+03:00","level":"INFO","msg":"closed","laddr":{"IP":"127.0.0.1","Port":44896,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
{"time":"2024-09-18T15:17:18.844777336+03:00","level":"INFO","msg":"closed","laddr":{"IP":"127.0.0.1","Port":2222,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":60988,"Zone":""},"lifetime":4,"reversed":true}
{"time":"2024-09-18T15:17:19.238986575+03:00","level":"INFO","msg":"closed","laddr":{"IP":"127.0.0.1","Port":44908,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
{"time":"2024-09-18T15:17:19.239013755+03:00","level":"INFO","msg":"closed","laddr":{"IP":"127.0.0.1","Port":2223,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":60370,"Zone":""},"lifetime":2,"reversed":true}
```

Show current activity:

```shell
$ sudo reverssh -active
active connections:
127.0.0.1:41924 lifetime=15 reversed=true
127.0.0.1:56068 lifetime=14 reversed=true
```

## License

[MIT](https://spdx.org/licenses/MIT.html) 