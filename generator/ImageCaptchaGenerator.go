package generator

import (
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
)

type ImageCaptchaGenerator interface {
	GenerateCaptchaImage(captchaExchange *model.CaptchaExchange) error

	WrapImageCaptchaInfo(captchaExchange *model.CaptchaExchange) (*model.ImageCaptchaInfo, error)
}
