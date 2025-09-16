# PuSu Engine client for Golang

This is a client for the PuSu Engine written in Go. PuSu Engine is a Pub-Sub engine.

More information on the server repository at [https://github.com/PuSuEngine/pusud](https://github.com/PuSuEngine/pusud).

Get it by running

```bash
go get github.com/PuSuEngine/gopusu
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

import "github.com/PuSuEngine/gopusu"

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


# Financial support

This project has been made possible thanks to [Cocreators](https://cocreators.ee) and [Lietu](https://lietu.net). You can help us continue our open source work by supporting us on [Buy me a coffee](https://www.buymeacoffee.com/cocreators).

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cocreators)
