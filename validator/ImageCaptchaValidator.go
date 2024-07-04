package validator

import (
	"fmt"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"strconv"
)

type ImageCaptchaValidator interface {
	GenerateImageCaptchaValidData(imageCaptchaInfo *model.ImageCaptchaInfo) (ValidData, error)
	Valid(imageCaptchaTrack *model.ImageCaptchaTrack, imageCaptchaValidData ValidData) (*model.ApiResponse, error)
}

type ValidData map[string]any

func (self ValidData) getInt(key string, defaultValue *int) (*int, error) {
	value, ok := self[key]
	if !ok {
		return defaultValue, nil
	}
	switch v := value.(type) {
	case int:
		return &v, nil
	case int32:
		cov := int(v)
		return &cov, nil
	case int64:
		cov := int(v)
		return &cov, nil
	default:
		str := fmt.Sprintf("%v", v)
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		cov := int(val)
		return &cov, nil
	}

}

func (self ValidData) ConstantKey(key string) bool {
	_, ok := self[key]
	return ok
}

func (self ValidData) getStr(key string, defaultValue *string) (*string, error) {
	value, ok := self[key]
	if !ok {
		return defaultValue, nil
	}
	switch v := value.(type) {
	case string:
		return &v, nil
	default:
		str := fmt.Sprintf("%v", v)
		return &str, nil
	}
}

func (self ValidData) getFloat(key string, defaultValue *float64) (*float64, error) {
	value, ok := self[key]
	if !ok {
		return defaultValue, nil
	}
	switch v := value.(type) {
	case float64:
		return &v, nil
	case float32:
		f := float64(v)
		return &f, nil
	case int32:
		f := float64(v)
		return &f, nil
	case int64:
		f := float64(v)
		return &f, nil
	default:
		// 先转换成字符串再转换成float64
		str := fmt.Sprintf("%v", v)
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		return &float, nil
	}
}
