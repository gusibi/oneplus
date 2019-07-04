package main

// https://github.com/go-swagger/go-swagger/issues/962

import (
	"context"
	handler "mobile-attribution/handler"

	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

func Handler(ctx context.Context, req scf.APIGatewayProxyRequest) (scf.APIGatewayProxyResponse, error) {

	return handler.GetMobileAttributionHandler(req)
}

func main() {
	cloudfunction.Start(Handler)
}
