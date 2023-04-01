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
				_, err = cs.GenerateCertificate(dto, lang)
				if err != nil {

					break
				}

			}

			// if all languages are processed success emit success event
			if i == lenLang {
				log.Println("Sertifika başarıyla üretildi: ")
			}

		}

	}()
}
