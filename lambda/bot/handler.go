package bot

import (
	"log"
)

func getService() *Service {
	svc, err := NewService()
	if err != nil {
		log.Panicf("ERROR: unable to create service - %v", err)
	}
	return svc
}

func BotHandler() (e error) {
	svc := getService()
	svc.logger.Info("scheduled job initiated")

	// test log
	svc.logger.Info("lambda handler triggered")

	return nil
}
