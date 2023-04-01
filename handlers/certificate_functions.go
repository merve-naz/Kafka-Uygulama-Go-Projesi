package handlers

import (
	"errors"
	"image"
	"log"

	"github.com/fogleman/gg"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/config"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/models"
)

// file assets
var FontBold = "font-Bold.ttf"
var FontMedium = "font-Medium.ttf"
var FontLight = "font-Light.ttf"
var FontSizeHigh = float64(50)
var FontSizeMid = float64(45)
var FontSizeLow = float64(40)

func FontLayoutPath(name string) string {
	return config.Cwdir + "/assets/fonts/" + name
}

func ImageLayoutPath(lang string) string {
	if lang == "tr" {
		return config.Cwdir + "/assets/layouts/layout_tr.jpg"
	}
	if lang == "fr" {
		return config.Cwdir + "/assets/layouts/layout_fr.jpg"
	}
	if lang == "en" {
		return config.Cwdir + "/assets/layouts/layout_en.jpg"
	}

	return config.Cwdir + "/assets/layouts/layout_tr.jpg"
}

func (mss *CertificateService) InitAssets() {
	FontBold = FontLayoutPath(FontBold)
	FontMedium = FontLayoutPath(FontMedium)
	FontLight = FontLayoutPath(FontLight)
	FontSizeHigh = float64(80)
	FontSizeMid = float64(40)
	FontSizeLow = float64(25)
}

func (mss *CertificateService) GenerateCertificate(dto models.Certificate, language string) (image.Image, error) {

	completedAt := dto.CompletedAt.Format("2006-01-02")

	absolutePath := ImageLayoutPath(language)
	layout, err := gg.LoadImage(absolutePath)
	if err != nil {
		return nil, errors.New("could not load certificate layout file")
	}

	// get layout width and height
	w := layout.Bounds().Dx()
	h := layout.Bounds().Dy()
	log.Println("******************** en: ", w, " boy ise ", h)
	// create canvas
	dc := gg.NewContext(w, h)
	dc.DrawImage(layout, 0, 0)
	dc.SetRGB(255, 255, 255)

	// 1 - write username
	x := float64(w / 2)
	y := float64((h / 2) - 100) //YAZI VEKTORLERINI DINAMIK AL.
	maxWidth := float64(w) - 60.0
	if err := dc.LoadFontFace(FontBold, FontSizeHigh); err != nil {
		return nil, err
	}
	// write text on it.
	dc.DrawStringWrapped(dto.Name, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	// 2 - write at 64(imgWidth / 2)
	y = float64((h / 2) + 40)
	if err := dc.LoadFontFace(FontMedium, FontSizeMid); err != nil {
		return nil, err
	}
	// write text
	dc.DrawStringWrapped(dto.Title, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	// 3 - write completed_at
	m := float64(w)
	z := float64(h)
	z = z - (z / 4) + 40
	m = m - (m / 4) - 40

	if err := dc.LoadFontFace(FontLight, FontSizeLow); err != nil {
		return nil, err
	}
	// write text
	dc.DrawStringWrapped(completedAt, m, z, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	// 4 - write registration
	x = float64(w/4) + 60
	dc.DrawStringWrapped(dto.RegistrationUUID, x, z, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	return dc.Image(), nil
}
