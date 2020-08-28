package bot


import (
"fmt"
"net/http"
"os"

"github.com/goop/service-edi-purchase-orders/awshelpers"
"github.com/goop/service-edi-purchase-orders/cloudformation/config"
"github.com/goop/service-edi-purchase-orders/spree"
"go.uber.org/zap"
)

type Service struct {
	logger            *zap.Logger
}

func NewService() (*Service, error) {
	var err error

	env := loadOrPanic("DEPLOY_ENV")
	config.SetVarsWithEnvironment(env)
	fmt.Println("DEPLOY_ENV: ", env)

	// Topic ARN for Admin SNS Alert
	topicArn := loadOrPanic("FAILED_EDI_PO_TOPIC")

	// Topic ARN for Goop Engineering SNS Alert
	engTopicArn := loadOrPanic("FAILED_EDI_ENG_TOPIC")

	// Set EDI Stock Locations
	var ediStockLocations = map[string]string{}
	for k, v := range config.EDI_STOCK_LOCATION_VARIABLES {
		ediStockLocations[k] = v
	}

	// build custom logger configuration
	logger, _ := cfg.Build()
	// flushes buffer, if any
	defer logger.Sync()

	configTable := func(name string) (t awshelpers.Table) {
		if err == nil {
			if t, err = awshelpers.NewDynamoTable(name); err != nil {
				err = NewInternalError(fmt.Sprintf("Could not access DynamoDB table %s", name), err)
			}
		}
		return
	}

	// TODO: Uncomment back when Spree Client is finished
	sp, err := spree.NewClient(loadOrPanic("SPREE_ACCESS_TOKEN"), loadOrPanic("SPREE_BASE_URL"), http.Client{})
	if err != nil {
		return nil, NewInternalError("Could not create Spree client", err)
	}

	// TODO: Add spree back to service struct
	svc := Service{
		logger:            logger,
		botTable:           configTable(config.EDI_PURCHASE_ORDERS_TABLE),
	}
	return &svc, err
}

func loadOrPanic(s string) string {
	logger, _ := cfg.Build()
	// flushes buffer, if any
	defer logger.Sync()
	out := os.Getenv(s)

	logger.Debug("ENV", zap.Any("s", s), zap.Any("out", out))

	if out == "" {
		logger.DPanic("Environment not set", zap.String("Env", s))
	}

	return out
}
