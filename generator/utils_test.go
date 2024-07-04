package generator

import (
	"github.com/golang/freetype"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestDraw(t *testing.T) {
	wh := int(math.Sqrt(50*50 + 50*50))
	img := common.NewTransparentImage(wh, wh)
	// 随机生成颜色
	fontColor := color.RGBA{
		R: uint8(common.GetRandomInt(0, 255)),
		G: uint8(common.GetRandomInt(0, 255)),
		B: uint8(common.GetRandomInt(0, 255)),
		A: 255,
	}
	file, _ := os.ReadFile("C:\\Users\\Thinkpad\\Desktop\\captcha\\手写字体\\ttf\\SIMSUN.TTC")
	font, _ := freetype.ParseFont(file)
	//var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}
	// 汉字编码范围
	start := 0x4e00
	end := 0x9fa5
	// 生成随机汉字
	hanzi := rand.Intn(end-start+1) + start
	str := string(rune(hanzi))

	img = img.DrawString(font, fontColor, str, (wh-50)/2, (wh-50)/2, 50)
	img = img.Rotate(45)
	img.WriteToFile("C:\\Users\\Thinkpad\\Desktop\\temp\\xxx.png")
}
