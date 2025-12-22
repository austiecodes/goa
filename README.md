# README

install `go` first, 1.24 or higher is recommended.

```shell
go install github.com/austiecodes/goa/cmd/goa@latest
```

the goa binary will be installed to `~/go/bin` directory.
add it to your `PATH` environment variable, like this:

```shell
export PATH=$PATH:~/go/bin
```

then add following config to your `mcp-server` config file:

```json
{
  "mcpServers": {
    "goa": {
      "command": "goa",
      "args": ["mcp"]
    }
  }
}
```

now you are ok to goa!