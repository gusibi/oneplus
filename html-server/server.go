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
	router.POST("/create", handler.CreateHtml)
	router.GET("/render/:id", handler.RenderCode)
	router.GET("/sleep/:seconds", handler.Sleep)

	log.Fatal(http.ListenAndServe(":8080", router))
}
