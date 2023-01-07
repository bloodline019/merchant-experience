package main

import (
	"github.com/gin-gonic/gin"
	"merchant-experience/handlers"
)

var r = gin.Default()

func main() {

	r.POST("/upload", handlers.HandleXlsxProcessing)
	r.POST("/getGoods", handlers.HandleGetGoods)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
