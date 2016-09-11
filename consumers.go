package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/linkosmos/mapop"
	"github.com/streadway/amqp"
)

const ACTIONINSERT = "insert"
const ACTIONUPDATE = "update"
const ACTIONDELETE = "delete"

type PayLoad struct {
	CorrelationID string `json:"correlation_id"`
	TenantID      string `json:"tenant_id"`
	Doc           User   `json:"doc"`
}

func StartConsumers(jobs chan amqp.Delivery, routines int, action string) {
	for i := 0; i < routines; i++ {
		switch action {
		case ACTIONINSERT:
			go InsertConsumer(i, jobs)
		case ACTIONUPDATE:
			go UpdateConsumer(i, jobs)
		case ACTIONDELETE:
			go DeleteConsumer(i, jobs)
		}
	}
}

func InsertConsumer(id int, jobs <-chan amqp.Delivery) {
	for job := range jobs {
		payLoad := PayLoad{}

		logFields := log.Fields{
			"worker_id":  id,
			"message_id": job.MessageId,
			"data":       string(job.Body),
		}

		log.WithFields(logFields).Infoln("Processing message")

		if err := json.Unmarshal(job.Body, &payLoad); err != nil {
			log.WithFields(mapop.Merge(logFields, log.Fields{
				"error": err,
			})).Errorln("Fail to unmarshal payload")
		}
	}
}

func UpdateConsumer(id int, jobs <-chan amqp.Delivery) {

}

func DeleteConsumer(id int, jobs <-chan amqp.Delivery) {

}
