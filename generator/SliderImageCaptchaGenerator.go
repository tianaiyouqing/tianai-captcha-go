package generator

import (
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
)

type SliderImageCaptchaGenerator struct {
	ResourceStore       resource.ImageCaptchaResourceStore
	ResourceImageReader *resource.ImageCaptchaResourceReaders
	ImageTransform      ImageTransform
}

func (self *SliderImageCaptchaGenerator) GenerateCaptchaImage(captchaExchange *model.CaptchaExchange) error {
	param := *captchaExchange.Param
	templateResource, err := resource.RequiredRandomGetTemplate(self.ResourceStore, param.CaptchaName, param.TemplateImageTag)
	if err != nil {
		return errors.Wrap(err, "滑动验证码获取模板失败")
	}
	bgImageResource, err := resource.RequiredRandomGetResource(self.ResourceStore, param.CaptchaName, param.BackgroundImageTag)
	if err != nil {
		return errors.Wrap(err, "滑动验证码获取背景图片失败")
	}
	bgImage, err := self.ResourceImageReader.GetResourceImage(bgImageResource)
	if err != nil {
		return errors.Wrap(err, "滑动验证码转换背景图片失败")
	}
	activeImage, err := self.ResourceImageReader.GetResourceImageByTemplate(templateResource, common.TEMPLATE_ACTIVE_IMAGE_NAME)
	if err != nil {
		return errors.Wrap(err, "滑动验证码获取模板失败")
	}
	fixedImage, err := self.ResourceImageReader.GetResourceImageByTemplate(templateResource, common.TEMPLATE_FIXED_IMAGE_NAME)
	if err != nil {
		return errors.Wrap(err, "滑动验证码获取模板失败")
	}
	randomX := common.GetRandomInt(fixedImage.Bounds().Dx()+5, bgImage.Bounds().Dx()-fixedImage.Bounds().Dx()-10)
	randomY := common.GetRandomInt(0, bgImage.Bounds().Dy()-fixedImage.Bounds().Dy())

	cutImage := bgImage.Cut(fixedImage, randomX, randomY)
	bgImage = bgImage.Overlay(fixedImage, randomX, randomY)
	if param.Obfuscate {
		// 混淆
	}
	cutImage = cutImage.Overlay(activeImage, 0, 0)
	matrixTemplate := common.NewTransparentImage(activeImage.Bounds().Dx(), bgImage.Bounds().Dy())
	matrixTemplate = matrixTemplate.Overlay(cutImage, 0, randomY)
	captchaExchange.BgImage = bgImage
	captchaExchange.TemplateImage = matrixTemplate
	captchaExchange.ResourceImage = bgImageResource
	captchaExchange.TemplateResource = templateResource
	captchaExchange.TransferData = map[string]int{
		"x": randomX,
		"y": randomY,
	}
	return nil
}

func (self *SliderImageCaptchaGenerator) WrapImageCaptchaInfo(captchaExchange *model.CaptchaExchange) (*model.ImageCaptchaInfo, error) {
	transform, err := self.ImageTransform.Transform(captchaExchange.Param, captchaExchange.BgImage,
		captchaExchange.TemplateImage,
		captchaExchange.ResourceImage,
		captchaExchange.TemplateResource, captchaExchange.CustomData)
	if err != nil {
		return nil, errors.Wrap(err, "图片转换异常")
	}
	randomX := captchaExchange.TransferData.(map[string]int)["x"]
	return &model.ImageCaptchaInfo{
		CaptchaName:           captchaExchange.Param.CaptchaName,
		CaptchaType:           model.TYPE_SLIDER,
		BackgroundImage:       transform.BgImageUrl,
		TemplateImage:         transform.TemplateImageUrl,
		TemplateImageTag:      captchaExchange.TemplateResource.Tag,
		TemplateImageWidth:    captchaExchange.TemplateImage.Bounds().Dx(),
		TemplateImageHeight:   captchaExchange.TemplateImage.Bounds().Dy(),
		BackgroundImageTag:    captchaExchange.ResourceImage.Tag,
		BackgroundImageWidth:  captchaExchange.BgImage.Bounds().Dx(),
		BackgroundImageHeight: captchaExchange.BgImage.Bounds().Dy(),
		RandomX:               &randomX,
	}, nil
}
