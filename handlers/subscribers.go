package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-playground/validator"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/models"
)

func (cs *CertificateService) InitSubcriber() {
	go func() {
		for {
			// listen message
			msg, err := cs.krCert.ReadMessage(context.Background())
			if err != nil {
				log.Println("CONSUMER: unable to read message from kafka:", err)
				continue
			}

			var dto models.Certificate
			err = json.Unmarshal(msg.Value, &dto)
			if err != nil {
				continue
			}

			validate := validator.New()
			err = validate.Struct(dto)
			if err != nil {
				continue
			}

			languages := []string{"tr", "en", "fr"}
			lenLang := len(languages)
			i := 0
			for i = 0; i < lenLang; i++ {
				lang := languages[i]
				// generate certificate
				cert, err := cs.GenerateCertificate(dto, lang)
				if err != nil {

					break
				}
				//certificate has been generated.

				certBuff, err := cs.TransformCertificateImage(cert, lang)
				if err != nil {

					break
				}
				err = cs.UploadCertificateFile(dto, certBuff, lang)
				if err != nil {

					break
				}
				//upload was successful

			}

			// if all languages are processed success emit success event
			if i == lenLang {
				log.Println("Sertifika başarıyla üretilip minio'ya kaydedildi")
			}

		}

	}()
}
