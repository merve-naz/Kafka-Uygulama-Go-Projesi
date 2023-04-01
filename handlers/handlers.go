package handlers

import (
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"

	"github.com/segmentio/kafka-go"
)

const (
	API_PREFIX = "/api-certificates"
	RN_PREFIX  = "bbrn:::certificateservice:::"
)

// CertificateService is a struct for certificate service
type CertificateService struct {
	inAppCache *persistence.InMemoryStore

	krCert   *kafka.Reader
	kwCertDL *kafka.Writer
	kwCert   *kafka.Writer
}

func NewCertificateService(

	inAppCache *persistence.InMemoryStore,
	krCert *kafka.Reader,
	kwCertDL *kafka.Writer,
	kwCert *kafka.Writer,

) *CertificateService {
	return &CertificateService{

		inAppCache: inAppCache,
		krCert:     krCert,
		kwCertDL:   kwCertDL,
		kwCert:     kwCert,
	}
}

type RespondJson struct {
	Status  bool        `json:"status" example:"true"`
	Intent  string      `json:"intent" example:"bbrn:::certificateservice:::/upload"`
	Message interface{} `json:"message" example:nil`
}

func respondJson(ctx *gin.Context, code int, intent string, message interface{}, err error) {
	if err == nil {
		ctx.JSON(code, RespondJson{
			Status:  true,
			Intent:  intent,
			Message: message,
		})
	} else {
		ctx.JSON(code, RespondJson{
			Status:  false,
			Intent:  intent,
			Message: err.Error(),
		})
	}
}

func (mss *CertificateService) InitRouter(r *gin.Engine) {
	// -- certificate service routes (group)
	v1 := r.Group(API_PREFIX)

	// manual trigger certificate topic for test purposes.
	v1.POST("/trigger-kafka", func(ctx *gin.Context) {
		code, data, err := mss.HandleManualTriggerKafka(ctx)
		respondJson(ctx, code, RN_PREFIX+"/trigger-kafka", data, err)
	})

}
