# reverssh

Reversing SSH tarpit.

You can specify port numbers (e.g., -p `22,2222`) for `reverssh` to attempt connections on each of the listed ports. 
If a connection is successful, `reverssh` will forward all incoming traffic back to the attacker's server, 
causing them to interact with their own system.

If no ports are specified, or if the provided ports are closed, `reverssh` operates as a standard SSH tarpit, 
sending one random byte per second.

## Installation

Download binary manually:

- [reverssh-linux-x86_64](https://raw.githubusercontent.com/jackcvr/reverssh/main/bin/x86_64/reverssh)
- [reverssh-linux-aarch64](https://raw.githubusercontent.com/jackcvr/reverssh/main/bin/aarch64/reverssh)
- [reverssh-linux-armv7l](https://raw.githubusercontent.com/jackcvr/reverssh/main/bin/armv7l/reverssh)

or via command:

```shell
`sh -c "wget -O /usr/local/bin/reverssh https://raw.githubusercontent.com/jackcvr/reverssh/main/bin/$(uname -m)/reverssh \
  && chmod +x /usr/local/bin/reverssh"`
```

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
```

## Examples

```shell
$ reverssh -b 0.0.0.0:2222 -p 22  # log sample when reversed connection succeeded
{"time":"2024-09-14T17:13:27.111626861+03:00","level":"INFO","msg":"listening","addr":"0.0.0.0:2222"}
{"time":"2024-09-14T17:13:32.080358768+03:00","level":"INFO","msg":"accepted","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:39680"}
{"time":"2024-09-14T17:13:32.08045588+03:00","level":"INFO","msg":"connected","laddr":"127.0.0.1:40136","raddr":"127.0.0.1:22"}
{"time":"2024-09-14T17:13:45.008896864+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:40136","raddr":"127.0.0.1:22"}
{"time":"2024-09-14T17:13:47.009419814+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:39680","lifetime":13}
```

```shell
$ reverssh -b 0.0.0.0:2222 -p 23  # log sample without successful connection to clients server
{"time":"2024-09-14T17:15:01.726948856+03:00","level":"INFO","msg":"listening","addr":"0.0.0.0:2222"}
{"time":"2024-09-14T17:15:04.231376092+03:00","level":"INFO","msg":"accepted","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:58262"}
{"time":"2024-09-14T17:15:11.239589332+03:00","level":"INFO","msg":"closed","laddr":"127.0.0.1:2222","raddr":"127.0.0.1:58262","lifetime":6}
```

## License

[MIT](https://spdx.org/licenses/MIT.html) 