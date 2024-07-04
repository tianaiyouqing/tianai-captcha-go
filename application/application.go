package application

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/generator"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"github.com/tianaiyouqing/tianai-captcha-go/validator"
)

type TianAiCaptchaApplication struct {
	ResourceStore         resource.ImageCaptchaResourceStore
	ResourceImageReader   *resource.ImageCaptchaResourceReaders
	ImageTransform        generator.ImageTransform
	providers             map[string]CaptchaGeneratorProvider
	processors            []ImageCaptchaPostProcessor
	CacheStore            CacheStore
	ImageCaptchaValidator validator.ImageCaptchaValidator
}

func (self *TianAiCaptchaApplication) GenerateCaptcha(param *model.GenerateParam) (*model.ImageCaptchaVO, error) {
	exchange := &model.CaptchaExchange{
		Param: param,
		CustomData: &model.CustomData{
			ViewData: make(map[string]any),
			Data:     make(map[string]any),
		},
	}
	captchaInfo, err := self.applyPostProcessorBeforeGenerate(exchange)
	if err != nil {
		return nil, errors.Wrap(err, "GenerateCaptcha error")
	}
	if captchaInfo != nil {
		return self.wrapImageCaptchaVO(captchaInfo, exchange)
	}

	provider, ok := self.providers[param.CaptchaName]
	if !ok {
		return nil, errors.Errorf("GenerateCaptcha error 未找到对应的验证码类型生成器 [%s]", param.CaptchaName)
	}
	captchaGenerator, err := provider(self)
	if err != nil {
		return nil, err
	}
	err = captchaGenerator.GenerateCaptchaImage(exchange)
	if err != nil {
		return nil, err
	}
	err = self.applyPostProcessorBeforeWrapImageCaptchaInfo(exchange)
	if err != nil {
		return nil, err
	}
	captchaInfo, err = captchaGenerator.WrapImageCaptchaInfo(exchange)
	if err != nil {
		return nil, err
	}
	captchaInfo.Data = exchange.CustomData

	err = self.applyPostProcessorAfterGenerateCaptchaImage(exchange, captchaInfo)
	if err != nil {
		return nil, err
	}
	return self.wrapImageCaptchaVO(captchaInfo, exchange)
}

func (self *TianAiCaptchaApplication) Valid(id string, track *model.ImageCaptchaTrack) (*model.ApiResponse, error) {
	validData, ok := self.CacheStore.getCache(id)
	if !ok {
		response := model.ApiResponse{}
		return response.Expire(), nil
	}
	return self.ImageCaptchaValidator.Valid(track, validData)
}

// ======================私有方法=====================

func (self *TianAiCaptchaApplication) wrapImageCaptchaVO(captchaInfo *model.ImageCaptchaInfo, exchange *model.CaptchaExchange) (*model.ImageCaptchaVO, error) {
	//生成uuid
	id := fmt.Sprintf("%s_%s", captchaInfo.CaptchaName, uuid.New().String())
	validData, err := self.ImageCaptchaValidator.GenerateImageCaptchaValidData(captchaInfo)
	if err != nil {
		return nil, errors.Wrap(err, "wrapImageCaptchaVO error")
	}
	if validData != nil && len(validData) != 0 {
		err = self.CacheStore.SetCache(id, validData, captchaInfo)
		if err != nil {
			return nil, errors.Wrap(err, "缓存数据出错")
		}
	}
	return &model.ImageCaptchaVO{
		Id:                    id,
		CaptchaName:           captchaInfo.CaptchaName,
		BackgroundImage:       captchaInfo.BackgroundImage,
		BackgroundImageTag:    captchaInfo.BackgroundImageTag,
		BackgroundImageHeight: captchaInfo.BackgroundImageHeight,
		BackgroundImageWidth:  captchaInfo.BackgroundImageWidth,
		TemplateImage:         captchaInfo.TemplateImage,
		TemplateImageTag:      captchaInfo.TemplateImageTag,
		TemplateImageHeight:   captchaInfo.TemplateImageHeight,
		TemplateImageWidth:    captchaInfo.TemplateImageWidth,
		Data:                  captchaInfo.Data.ViewData,
	}, nil
}

func (self *TianAiCaptchaApplication) applyPostProcessorAfterGenerateCaptchaImage(exchange *model.CaptchaExchange, imageCaptchaInfo *model.ImageCaptchaInfo) error {
	for _, processor := range self.processors {
		err := processor.AfterGenerateCaptchaImage(exchange, imageCaptchaInfo, self)
		if err != nil {
			return errors.Wrap(err, "applyPostProcessorBeforeWrapImageCaptchaInfo error")
		}
	}
	return nil
}

func (self *TianAiCaptchaApplication) applyPostProcessorBeforeWrapImageCaptchaInfo(exchange *model.CaptchaExchange) error {
	for _, processor := range self.processors {
		err := processor.BeforeWrapImageCaptchaInfo(exchange, self)
		if err != nil {
			return errors.Wrap(err, "applyPostProcessorBeforeWrapImageCaptchaInfo error")
		}
	}
	return nil
}

func (self *TianAiCaptchaApplication) applyPostProcessorBeforeGenerate(exchange *model.CaptchaExchange) (*model.ImageCaptchaInfo, error) {
	for _, processor := range self.processors {
		image, err := processor.BeforeGenerateCaptchaImage(exchange, self)
		if err != nil {
			return image, errors.Wrap(err, "applyPostProcessorBeforeGenerate error")
		}
		if image != nil {
			return image, nil
		}
	}
	return nil, nil
}
