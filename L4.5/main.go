package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

const pprofAddr = "localhost:6060"
const port = ":8080"

type request struct {
	First  int `json:"first"`
	Second int `json:"second"`
}

type response struct {
	Result int `json:"result"`
}

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

func add(c *gin.Context) {

	var request request

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := request.First + request.Second

	c.JSON(http.StatusOK, response{Result: result})

}
