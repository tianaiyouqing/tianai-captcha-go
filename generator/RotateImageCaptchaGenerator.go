package generator

import (
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
)

type RotateImageCaptchaGenerator struct {
	ResourceStore       resource.ImageCaptchaResourceStore
	ResourceImageReader *resource.ImageCaptchaResourceReaders
	ImageTransform      ImageTransform
}

func (self *RotateImageCaptchaGenerator) GenerateCaptchaImage(captchaExchange *model.CaptchaExchange) error {
	param := *captchaExchange.Param
	templateResource, err := resource.RequiredRandomGetTemplate(self.ResourceStore, param.CaptchaName, param.TemplateImageTag)
	if err != nil {
		return errors.Wrap(err, "旋转验证码获取模板失败")
	}
	bgImageResource, err := resource.RequiredRandomGetResource(self.ResourceStore, param.CaptchaName, param.BackgroundImageTag)
	if err != nil {
		return errors.Wrap(err, "旋转验证码获取背景图片失败")
	}
	bgImage, err := self.ResourceImageReader.GetResourceImage(bgImageResource)
	if err != nil {
		return errors.Wrap(err, "旋转验证码转换背景图片失败")
	}
	activeImage, err := self.ResourceImageReader.GetResourceImageByTemplate(templateResource, common.TEMPLATE_ACTIVE_IMAGE_NAME)
	if err != nil {
		return errors.Wrap(err, "旋转验证码获取模板失败")
	}
	fixedImage, err := self.ResourceImageReader.GetResourceImageByTemplate(templateResource, common.TEMPLATE_FIXED_IMAGE_NAME)
	if err != nil {
		return errors.Wrap(err, "旋转验证码获取模板失败")
	}
	x := bgImage.Bounds().Dx()/2 - fixedImage.Bounds().Dx()/2
	y := bgImage.Bounds().Dy()/2 - fixedImage.Bounds().Dy()/2

	cutImage := bgImage.Cut(fixedImage, x, y)
	bgImage = bgImage.Overlay(fixedImage, x, y)
	if param.Obfuscate {
		// 混淆
	}

	randomX := common.GetRandomInt(fixedImage.Bounds().Dx()+10, bgImage.Bounds().Dx()-fixedImage.Bounds().Dx()-10)
	degree := float64(360) - float64(randomX)/(float64(bgImage.Bounds().Dx())/float64(360))
	//
	matrixTemplate := common.NewTransparentImage(cutImage.Bounds().Dx(), bgImage.Bounds().Dy())
	cutImage = cutImage.Rotate(degree)
	cutImage = cutImage.Overlay(activeImage, 0, 0)

	matrixTemplate = matrixTemplate.Overlay(cutImage,
		matrixTemplate.Bounds().Dx()/2-cutImage.Bounds().Dy()/2,
		matrixTemplate.Bounds().Dy()/2-cutImage.Bounds().Dy()/2)
	captchaExchange.BgImage = bgImage
	captchaExchange.TemplateImage = matrixTemplate
	captchaExchange.ResourceImage = bgImageResource
	captchaExchange.TemplateResource = templateResource
	captchaExchange.TransferData = map[string]int{
		"x":      randomX,
		"degree": int(degree),
	}
	return nil
}

func (self *RotateImageCaptchaGenerator) WrapImageCaptchaInfo(captchaExchange *model.CaptchaExchange) (*model.ImageCaptchaInfo, error) {
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
