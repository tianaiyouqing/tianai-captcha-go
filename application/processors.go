package application

import (
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
)

type ImageCaptchaPostProcessor interface {
	BeforeGenerateCaptchaImage(exchange *model.CaptchaExchange, app *TianAiCaptchaApplication) (*model.ImageCaptchaInfo, error)
	BeforeWrapImageCaptchaInfo(exchange *model.CaptchaExchange, app *TianAiCaptchaApplication) error
	AfterGenerateCaptchaImage(exchange *model.CaptchaExchange, imageCaptchaInfo *model.ImageCaptchaInfo, app *TianAiCaptchaApplication) error
}
