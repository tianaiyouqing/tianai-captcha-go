package resource

import (
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"image"
	"net/http"
	"os"
	"strings"
)

type ImageCaptchaResourceReader func(resource *model.Resource) (*common.Image, error)

func NewDefaultImageCaptchaResourceReaders() *ImageCaptchaResourceReaders {
	readers := &ImageCaptchaResourceReaders{}
	readers.AddResourceReader(GetResourceImageByFile)
	readers.AddResourceReader(GetResourceImageByUrl)
	readers.AddResourceReader(GetResourceImageByDefaultFile)
	return readers
}

type ImageCaptchaResourceReaders struct {
	readers []ImageCaptchaResourceReader
}

func (self *ImageCaptchaResourceReaders) AddResourceReader(reader ImageCaptchaResourceReader) {
	self.readers = append(self.readers, reader)
}

func (self *ImageCaptchaResourceReaders) GetResourceImage(resource *model.Resource) (*common.Image, error) {
	for _, reader := range self.readers {
		img, err := reader(resource)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		if img != nil {
			return img, nil
		}
	}
	return nil, errors.Errorf("解析资源类型 %s 失败", resource.ResourceType)
}
func (self *ImageCaptchaResourceReaders) GetResourceImageByTemplate(template *model.ResourceMap, key string) (image.Image, error) {
	if template == nil {
		return nil, errors.Errorf("获取模板资源失败，模板为空")
	}
	resource := template.Get(key)
	if resource == nil {
		return nil, errors.Errorf("获取模板资源失败，模板资源为空, key:[%s]", key)
	}
	return self.GetResourceImage(resource)
}

// ================== default reader impl =====================

func GetResourceImageByFile(resource *model.Resource) (*common.Image, error) {
	if !strings.EqualFold(resource.ResourceType, "file") {
		return nil, nil
	}
	file, err := os.Open(resource.Data)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, t, err := image.Decode(file)
	return common.NewImage(img, t), err
}

func GetResourceImageByUrl(resource *model.Resource) (*common.Image, error) {
	if !strings.EqualFold(resource.ResourceType, "url") {
		return nil, nil
	}
	resp, err := http.Get(resource.Data)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()
	img, t, err := image.Decode(body)
	return common.NewImage(img, t), err
}
