package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// HandleManualTriggerKafka godoc
// @Summary manually trigger kafka with payload
// @Schemes
// @Description manually trigger kafka with payload
// @Tags Certificate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param certificate body models.Certificate  true "certification dto"
// @Success 200 {object} handlers.RespondJson "kafka triggered manually success"
// @Failure 400 {object} handlers.RespondJson "invalid certificate info for trigger"
// @Failure 500 {object} handlers.RespondJson "internal server error while trigger kafka"
// @Router /trigger-kafka [post]
func (cs *CertificateService) HandleManualTriggerKafka(c *gin.Context) (int, interface{}, error) {
	// get params from request body
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); !errors.Is(err, nil) {
		return http.StatusBadRequest, nil, errors.New("invalid generate certificate body " + err.Error())
	}

	// jsonize the message (it should be a map[string]interface{})
	value, err := json.Marshal(params)
	if err != nil {
		return http.StatusBadRequest, nil, errors.New("unable to jsonize certificate: " + err.Error())
	}

	// send message to kafka for trigger
	// now as format: 2023/03/12 18:41:11
	now := time.Now().Format(time.RFC3339)
	log.Println("KAFKA_TRIGGER: " + now + " - " + "SENDING...")
	if err = cs.kwCert.WriteMessages(context.Background(), kafka.Message{Value: value}); err != nil {
		return http.StatusInternalServerError, nil, errors.New("unable to trigger kafka: " + err.Error())
	}
	log.Println("KAFKA_TRIGGER: " + now + " - " + "SENT.")

	// return success
	return http.StatusOK, "kafka triggered successfully", nil
}
