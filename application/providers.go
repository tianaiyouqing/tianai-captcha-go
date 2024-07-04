package application

import (
	"github.com/golang/freetype/truetype"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/generator"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
)

type CaptchaGeneratorProvider func(app *TianAiCaptchaApplication) (generator.ImageCaptchaGenerator, error)

// =============== 一些验证码的生成器实现 =================

func CreateWordClickProvider(fonts []*truetype.Font) (string, CaptchaGeneratorProvider) {
	if fonts == nil {
		font, _ := resource.GetDefaultFont()
		fonts = []*truetype.Font{font}
	}
	return common.CAPTCHA_NAME_WORD_CLICK, func(app *TianAiCaptchaApplication) (generator.ImageCaptchaGenerator, error) {
		return &generator.WordClickCaptchaGenerator{
			ResourceStore:       app.ResourceStore,
			ResourceImageReader: app.ResourceImageReader,
			ImageTransform:      app.ImageTransform,
			Fonts:               fonts,
		}, nil
	}
}

func CreateRotateProvider() (string, CaptchaGeneratorProvider) {
	return common.CAPTCHA_NAME_ROTATE, func(app *TianAiCaptchaApplication) (generator.ImageCaptchaGenerator, error) {
		return &generator.RotateImageCaptchaGenerator{
			ResourceStore:       app.ResourceStore,
			ResourceImageReader: app.ResourceImageReader,
			ImageTransform:      app.ImageTransform,
		}, nil
	}
}

func CreateSliderProvider() (string, CaptchaGeneratorProvider) {
	return common.CAPTCHA_NAME_SLIDER, func(app *TianAiCaptchaApplication) (generator.ImageCaptchaGenerator, error) {
		return &generator.SliderImageCaptchaGenerator{
			ResourceStore:       app.ResourceStore,
			ResourceImageReader: app.ResourceImageReader,
			ImageTransform:      app.ImageTransform,
		}, nil
	}
}
