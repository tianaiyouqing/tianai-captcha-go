package resource

import (
	"embed"
	"github.com/golang/freetype/truetype"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"image"
	"strings"
)

//go:embed default/*
var defaultResources embed.FS

func GetResourceImageByDefaultFile(resource *model.Resource) (*common.Image, error) {
	if !strings.EqualFold(resource.ResourceType, "default") {
		return nil, nil
	}
	file, err := defaultResources.Open(resource.Data)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, t, err := image.Decode(file)
	return common.NewImage(img, t), err
}

func AddDefaultResources(store ImageCaptchaResourceStore) {
	for _, r := range GetDefaultResources() {
		store.AddResource(common.CAPTCHA_NAME_SLIDER, r)
		store.AddResource(common.CAPTCHA_NAME_ROTATE, r)
		store.AddResource(common.CAPTCHA_NAME_WORD_CLICK, r)
	}

	sliderTemplates := GetDefaultSliderTemplates()
	for _, t := range sliderTemplates {
		store.AddTemplate(common.CAPTCHA_NAME_SLIDER, t)
	}

	rotateTemplates := GetDefaultRotateTemplate()
	for _, t := range rotateTemplates {
		store.AddTemplate(common.CAPTCHA_NAME_ROTATE, t)
	}
}

func GetDefaultFont() (*truetype.Font, error) {
	file, err := defaultResources.ReadFile("default/fonts/SIMSUN.TTC")
	if err != nil {
		return nil, err
	}
	return truetype.Parse(file)
}

func GetDefaultResources() []*model.Resource {
	return []*model.Resource{
		{
			ResourceType: "default",
			Data:         "default/image/1.jpeg",
		},
	}
}

func GetDefaultSliderTemplates() []*model.ResourceMap {
	return []*model.ResourceMap{
		{
			Tag: "default",
			ResourceMap: map[string]*model.Resource{
				"active.png": {
					ResourceType: "default",
					Data:         "default/templates/1/active.png",
				},
				"fixed.png": {
					ResourceType: "default",
					Data:         "default/templates/1/fixed.png",
				},
			},
		},
		{
			Tag: "default",
			ResourceMap: map[string]*model.Resource{
				"active.png": {
					ResourceType: "default",
					Data:         "default/templates/2/active.png",
				},
				"fixed.png": {
					ResourceType: "default",
					Data:         "default/templates/2/fixed.png",
				},
			},
		},
	}
}

func GetDefaultRotateTemplate() []*model.ResourceMap {
	return []*model.ResourceMap{
		{
			Tag: "default",
			ResourceMap: map[string]*model.Resource{
				"active.png": {
					ResourceType: "default",
					Data:         "default/templates/3/active.png",
				},
				"fixed.png": {
					ResourceType: "default",
					Data:         "default/templates/3/fixed.png",
				},
			},
		},
	}
}
