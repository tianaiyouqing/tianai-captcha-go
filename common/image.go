package common

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/tianaiyouqing/tianai-captcha-go/common/imaging"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type Image struct {
	draw.Image
	Type string
}

func NewTransparentImage(width int, height int) *Image {
	return NewImage(image.NewRGBA(image.Rect(0, 0, width, height)), "png")
}
func NewImage(img image.Image, imgType string) *Image {
	return &Image{
		Image: toDrawImage(img),
		Type:  imgType,
	}
}
func toDrawImage(img image.Image) draw.Image {
	if d, ok := img.(draw.Image); ok {
		return d
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}
func (self *Image) Cut(templateImage image.Image, xPos int, yPos int) *Image {
	bw := templateImage.Bounds().Dx()
	bh := templateImage.Bounds().Dy()
	resImage := image.NewRGBA(image.Rect(0, 0, bw, bh))
	for y := 0; y < bh; y++ {
		for x := 0; x < bw; x++ {
			_, _, _, a := templateImage.At(x, y).RGBA()
			if a > 100 {
				bgRgb := self.At(x+xPos, y+yPos)
				resImage.Set(x, y, bgRgb)
			}
		}
	}
	return NewImage(resImage, "png")
}
func (self *Image) Overlay(overlayImage image.Image, xPos int, yPos int) *Image {
	dstRect := image.Rect(xPos, yPos, overlayImage.Bounds().Dx()+xPos, overlayImage.Bounds().Dy()+yPos)
	draw.Draw(self, dstRect, overlayImage, overlayImage.Bounds().Min, draw.Over)
	return self
}

func (self *Image) Rotate(angle float64) *Image {
	angle = -angle
	rotate := imaging.Rotate(self, angle, color.Transparent)
	return NewImage(rotate, "png")
}

func (self *Image) DrawString(font *truetype.Font, c color.Color, str string, x int, y int, fontsize float64) *Image {
	ctx := freetype.NewContext()
	// default 72dpi
	ctx.SetDst(self)

	ctx.SetClip(self.Bounds())
	ctx.SetSrc(image.NewUniform(c))
	ctx.SetFontSize(fontsize)
	ctx.SetFont(font)
	// 写入文字的位置
	pt := freetype.Pt(x, y+int(-fontsize/6)+ctx.PointToFixed(fontsize).Ceil())
	_, _ = ctx.DrawString(str, pt)
	return self
}

func (self *Image) WriteToFile(filePath string) {
	file, _ := os.Create(filePath)
	defer file.Close()
	if strings.EqualFold("jpg", self.Type) || strings.EqualFold("jpeg", self.Type) {
		err := jpeg.Encode(file, self, nil)
		if err != nil {
			return
		}
	} else if strings.EqualFold("png", self.Type) {
		_ = png.Encode(file, self)
	}
}
