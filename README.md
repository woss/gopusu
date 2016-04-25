# PuSu Engine client for Golang

This is a client for the PuSu Engine written in Go. PuSu Engine is a Pub-Sub engine.

More information on the server repository at [https://github.com/lietu/pusud](https://github.com/lietu/pusud).

Get it by running

```bash
go get github.com/lietu/gopusu
```


## Usage

Examples are under `tests/`. Running the tests will require screen or 
multiple terminal windows.

One terminal window, running a `listener`;

```
cd tests
go run listener.go
```

Running basic test:
```
cd tests
go run basic.go
```

Running `throughput`:

```
cd tests
go run throughput.go
```

Short example how to connect to PuSu server:

```go
package main

import "github.com/lietu/gopusu"

func main() {
    pc, _ := gopusu.NewPuSuClient("127.0.0.1", 55000)
    defer pc.Close()
    pc.Authorize("foo")
    pc.Subscribe("channel.1", listener)
    pc.Publish("channel.2", "message")
}

func listener(msg *gopusu.Publish) {
	// ...
}
```


## License

Short version: MIT + New BSD.

Long version: Read the LICENSE.md -file.
