package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	InitRouter(r)
	Initdb()
	r.Run(":8001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
