package main


import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

func init() {
	for k, v := range config.AWS_ENVIRONMENT_VARIABLES {
		os.Setenv(k, v)
	}
}

func main() {
	lambda.Start(purchaseorders.PurchaseOrdersIngestHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
