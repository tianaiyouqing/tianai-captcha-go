package resource

import (
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
)

type ImageCaptchaResourceStore interface {
	RandomGetResource(captchaType string, tag *string) *model.Resource

	RandomGetTemplate(captchaType string, tag *string) *model.ResourceMap

	AddResource(captchaType string, resource *model.Resource)

	AddTemplate(captchaType string, template *model.ResourceMap)
}

func RequiredRandomGetTemplate(store ImageCaptchaResourceStore, captchaType string, tag *string) (*model.ResourceMap, error) {
	if store == nil {
		return nil, errors.Errorf("验证码类型[%s]获取模板错误， store为空,请手动配置resource.SetStore()", captchaType)
	}
	template := store.RandomGetTemplate(captchaType, tag)
	if template == nil {
		return nil, errors.Errorf("验证码类型[%s]获取模板错误， 模板为空", captchaType)
	}
	return template, nil
}
func RequiredRandomGetResource(store ImageCaptchaResourceStore, captchaType string, tag *string) (*model.Resource, error) {
	if store == nil {
		return nil, errors.Errorf("验证码类型[%s]获取资源错误， store为空,请手动配置resource.SetStore()", captchaType)
	}
	resource := store.RandomGetResource(captchaType, tag)
	if resource == nil {
		return nil, errors.Errorf("验证码类型[%s]获取资源错误， 资源为空", captchaType)
	}
	return resource, nil
}
