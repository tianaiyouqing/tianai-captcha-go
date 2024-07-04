package generator

import (
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"image"
	"image/jpeg"
	"image/png"
)

type ImageTransform interface {
	Transform(param *model.GenerateParam,
		bgImage image.Image,
		templateImage image.Image,
		bgResource *model.Resource,
		templateResource *model.ResourceMap,
		customData *model.CustomData,
	) (*TransFormData, error)
}

type TransFormData struct {
	BgImageUrl       string
	TemplateImageUrl string
	Data             *any
}

// base64 实现

func NewBase64ImageTransform() *Base64ImageTransform {
	return &Base64ImageTransform{}
}

type Base64ImageTransform struct{}

func (self *Base64ImageTransform) Transform(param *model.GenerateParam,
	bgImage image.Image, templateImage image.Image,
	bgResource *model.Resource, templateResource *model.ResourceMap,
	customData *model.CustomData) (*TransFormData, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, bgImage, nil)
	if err != nil {
		return nil, errors.Wrap(err, "encode bgImage to jpeg error")
	}
	bgImageUrl := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	var res = TransFormData{
		BgImageUrl: bgImageUrl,
	}
	if templateImage != nil {
		buf = bytes.Buffer{}
		err = png.Encode(&buf, templateImage)
		if err != nil {
			return nil, errors.Wrap(err, "encode templateImage to png error")
		}
		res.TemplateImageUrl = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	}
	return &res, nil

}
