package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const pprofAddr = "localhost:6060"
const port = ":8080"

func main() {

	handler := gin.Default()
	handler.POST("/add", add)

	go func() {
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	if err := handler.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

}

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 64)
		return &b
	},
}

func add(c *gin.Context) {

	bufPtr := bufPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]
	defer bufPool.Put(bufPtr)

	n, err := c.Request.Body.Read(buf[:cap(buf)])
	if err != nil && err != io.EOF {
		log.Printf("failed to read request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	buf = buf[:n]

	var first, second int
	_, err = fmt.Sscanf(string(buf), `{"first":%d,"second":%d}`, &first, &second)
	if err != nil {
		log.Printf("failed to parse JSON from body %s: %v", string(buf), err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "application/json")
	_, err = c.Writer.Write([]byte(`{"result":` + strconv.Itoa(first+second) + `}`))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}

}
