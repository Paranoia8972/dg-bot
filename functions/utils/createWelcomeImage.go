package utils

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func CreateWelcomeImage(s *discordgo.Session, user *discordgo.User, guildID string) {

	// load local bg img
	bgFile, err := os.Open("assets/bg.png")
	if err != nil {
		return
	}
	defer bgFile.Close()
	bgImg, err := png.Decode(bgFile)
	if err != nil {
		return
	}

	// download users pfp
	pfpResp, err := http.Get(user.AvatarURL("512"))
	if err != nil {
		return
	}
	defer pfpResp.Body.Close()
	pfpImg, err := png.Decode(pfpResp.Body)
	if err != nil {
		return
	}

	// resize pfp
	pfpImg = resize.Resize(512, 512, pfpImg, resize.Lanczos3)

	// create new image
	outputImg := image.NewRGBA(bgImg.Bounds())
	draw.Draw(outputImg, bgImg.Bounds(), bgImg, image.Point{}, draw.Over)

	dc := gg.NewContextForRGBA(outputImg)

	// calculate position center pfp
	pfpX := float64((bgImg.Bounds().Dx() - pfpImg.Bounds().Dx()) / 2)
	pfpY := float64((bgImg.Bounds().Dy() - pfpImg.Bounds().Dy()) / 2)

	// move pfp up
	moveUp := 50.0 // up by 50 pixels
	pfpY -= moveUp

	// circular mask for pfp
	dc.DrawCircle(pfpX+256, pfpY+256, 256)
	dc.Clip()

	dc.DrawImage(pfpImg, int(pfpX), int(pfpY))

	dc.ResetClip()

	// white border for pfp
	dc.SetLineWidth(10)
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(pfpX+256, pfpY+256, 256)
	dc.Stroke()

	// add username
	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace("assets/Geist-Bold.ttf", 100)
	if err != nil {
		return
	}

	// adjust y to add spacing and move up
	spacing := 50.0 // spacing between pfp and username
	dc.DrawStringAnchored(user.GlobalName, float64(bgImg.Bounds().Dx()/2), float64(pfpY+float64(pfpImg.Bounds().Dy())+30+spacing), 0.5, 0.5)
	dc.Fill()

	var buf bytes.Buffer
	err = png.Encode(&buf, outputImg)
	if err != nil {
		return
	}

}
