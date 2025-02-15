package main

import (
	"context"

	"github.com/skyrenx/blog-api-go/http/controller"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import PostgreSQL driver as a blank import
)

var ginLambda *ginadapter.GinLambda

func init() {

	// Create your Gin router and define routes.
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/aurora", controller.AuroraExampleHandler)

	// Wrap the router with the Lambda adapter.
	ginLambda = ginadapter.New(router)
}

// Assume ginLambda is declared and initialized in init()
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
