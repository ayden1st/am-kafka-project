package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"am-kafka-project/internal/kafka/producer"
	"am-kafka-project/internal/model/alert"
)

// NewAlerts - function that creates a new alert handler.
// It takes a producer service as an argument and returns a gin.HandlerFunc.
func NewAlerts(producer producer.ProducerService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		parent := context.Background()
		defer parent.Done()

		var alerts alert.Alerts

		// Bind the request body to the alerts variable
		err := ctx.Bind(&alerts)
		if err != nil {
			// If there is an error while binding the JSON, return a bad request response
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("error while binding json: %s", err.Error()),
				},
			})
			ctx.Abort()
			return
		}

		// Loop through each alert in the alerts
		for _, o_alert := range alerts.Alerts {

			// Serialize the alert
			payload, err := producer.Serializer.Serialize(producer.Topic, o_alert.ToKafkaAlert())
			if err != nil {
				// If there is an error while pushing the message to Kafka
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": map[string]interface{}{
						"message": fmt.Sprintf("failed to serialize payload: %s", err.Error()),
					},
				})
				ctx.Abort()
				return
			}

			// Push the serialized alert to Kafka
			err = producer.Push(producer.Topic, nil, payload)
			if err != nil {
				// If there is an error while pushing the message to Kafka
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": map[string]interface{}{
						"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
					},
				})

				ctx.Abort()
				return
			}
		}

		// If all alerts are successfully pushed to Kafka, return a success response
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "success push data into kafka",
		})
	}
}
