
###
```
go get -u github.com/jonnywang/redcon2
```

### example
```
package main

import (
	"github.com/tidwall/redcon"
	"log"
	"github.com/jonnywang/redcon2"
)

func main()  {
	mux := redcon2.NewRedconServeMux()
	mux.Handle("version", func (conn redcon.Conn, cmd redcon.Command) {
		conn.WriteBulkString("1.0.0")
	})

	err := mux.Run(":9191")
	if err != nil {
		log.Fatal(err)
	}
}

```