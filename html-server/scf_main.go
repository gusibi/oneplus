package main

// https://github.com/go-swagger/go-swagger/issues/962

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/julienschmidt/httprouter"
	scf "github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/cloudfunction"

	handler "html-server/handler"
	"html-server/httpadapter"
)

var httpAdapter *httpadapter.HandlerAdapter

func init() {
	log.Println("start server...")
	handler.InitDB()
	router := httprouter.New()
	router.GET("/", handler.EchoHandler)
	router.POST("/code/render", handler.CodeRender)
	router.POST("/codes", handler.CreateHtml)
	router.GET("/codes/:id", handler.RenderCode)

	httpAdapter = httpadapter.New(router)
	log.Println("adapter: ", httpAdapter)
}

// Handler go swagger aws lambda handler
func Handler(req events.APIGatewayProxyRequest) (scf.APIGatewayProxyResponse, error) {

	return httpAdapter.Proxy(req)
}

func main() {
	cloudfunction.Start(Handler)
}
