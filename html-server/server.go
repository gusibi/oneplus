package main

import (
	"log"
	"net/http"

	handler "html-server/handler"

	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("start server..")
	handler.InitDB()
	router := httprouter.New()
	router.GET("/", handler.EchoHandler)
	router.POST("/code/render", handler.CodeRender)
	router.POST("/codes", handler.CreateHtml)
	router.GET("/codes/:id", handler.RenderCode)

	log.Fatal(http.ListenAndServe(":8080", router))
}
