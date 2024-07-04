package common

import (
	"image/color"
	"math/rand"
)

// 生成指定数区间的随机数
func GetRandomInt(min int, max int) int {
	return min + int(rand.Int31n(int32(max-min)))
}

func GetRandomHanZi() string {
	start := 0x4e00
	end := 0x9fa5
	// 生成随机汉字
	hanzi := rand.Intn(end-start+1) + start
	return string(rune(hanzi))
}

func GetRandomColor() color.Color {
	return color.RGBA{
		R: uint8(GetRandomInt(0, 255)),
		G: uint8(GetRandomInt(0, 255)),
		B: uint8(GetRandomInt(0, 255)),
		A: 255,
	}
}
