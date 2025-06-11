```markdown
[![Build Status][gha-workflow-badge]][gha-workflow-url] [![Latest release][gha-latest-release]][gha-releases-url]

# Kong Plugin Development Kit - Go edition

Docs: https://pkg.go.dev/github.com/Kong/go-pdk

The Kong Plugin Development Kit (PDK) for Go lets you write custom Kong Gateway plugins in Go.

## Getting Started

To write a Kong Gateway plugin in Go:

1. Define a `struct` type to hold configuration.
2. Write a `New()` function to create instances of your struct.
3. Add methods to that struct to handle request phases.
4. Include the `go-pdk/server` sub-library.
5. Add a `main()` function that calls `server.StartServer(New, Version, Priority)`.
6. Compile as an executable with `go build`.

## Plugin Structure

### Configuration struct

The plugin uses a Go `struct` to receive configuration from Kong:

```go
type MyConfig struct {
    Path   string `json:"my_file_path"`
    Reopen bool   `json:"reopen"`
}
```

### `New()` constructor

Define a function to return a new instance of your config struct:

```go
func New() interface{} {
    return &MyConfig{}
}
```

### `main()` function

Include `github.com/Kong/go-pdk/server` and start the plugin:

```go
func main () {
  server.StartServer(New, Version, Priority)
}
```

When run, the plugin creates a socket file within the Kong prefix directory.

### Phase handlers

You can implement logic in the following phases using the same signature:

```go
func (conf *MyConfig) Access(kong *pdk.PDK) {
  ...
}
```

Supported phases:

- `Certificate`
- `Rewrite`
- `Access`
- `Response` *(enables buffered proxy mode)*
- `Preread`
- `Log`

### Version and priority

Set execution order with constants:

```go
const Version = "1.0.0"
const Priority = 1
```

Higher priority runs earlier.

## Configuring with `kong.conf`

To register plugins in Kong, define plugin server settings in `kong.conf`:

```conf
pluginserver_names = my-plugin

pluginserver_my_plugin_socket = /usr/local/kong/my-plugin.socket
pluginserver_my_plugin_start_cmd = /usr/local/bin/my-plugin
pluginserver_my_plugin_query_cmd = /usr/local/bin/my-plugin -dump
```

The socket and start command lines can be omitted if using defaults.

## Example Plugins

Explore [example plugins](https://github.com/Kong/go-pdk/tree/master/examples).

[gha-workflow-badge]: https://github.com/Kong/go-pdk/actions/workflows/test.yml/badge.svg
[gha-workflow-url]: https://github.com/Kong/go-pdk/actions/workflows/test.yml
[gha-latest-release]: https://img.shields.io/github/v/release/Kong/go-pdk.svg
[gha-releases-url]: https://github.com/Kong/go-pdk/releases
```
