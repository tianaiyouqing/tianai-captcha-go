package resource

import (
	"fmt"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"math/rand"
	"sync"
)

func NewMemoryImageCaptchaResourceStore() *MemoryImageCaptchaResourceStore {
	return &MemoryImageCaptchaResourceStore{
		templateResourceTagMap: make(map[string][]*model.ResourceMap),
		resourceTagMap:         make(map[string][]*model.Resource),
	}
}

type MemoryImageCaptchaResourceStore struct {
	templateResourceTagMap map[string][]*model.ResourceMap
	resourceTagMap         map[string][]*model.Resource
	mu                     sync.Mutex
}

func (self *MemoryImageCaptchaResourceStore) RandomGetResource(captchaType string, tag *string) *model.Resource {
	resources := self.resourceTagMap[mergeTypeAndTag(captchaType, tag)]
	if len(resources) == 0 {
		return nil
	}
	if len(resources) == 1 {
		return resources[0]
	}
	randomIndex := rand.Intn(len(resources))
	return resources[randomIndex]
}

func (self *MemoryImageCaptchaResourceStore) RandomGetTemplate(captchaType string, tag *string) *model.ResourceMap {
	templates := self.templateResourceTagMap[mergeTypeAndTag(captchaType, tag)]
	if len(templates) == 0 {
		return nil
	}
	if len(templates) == 1 {
		return templates[0]
	}
	randomIndex := rand.Intn(len(templates))
	return templates[randomIndex]
}

func (self *MemoryImageCaptchaResourceStore) AddResource(captchaType string, resource *model.Resource) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if resource.Tag == "" {
		resource.Tag = common.DEFAULT_TAG
	}
	key := mergeTypeAndTag(captchaType, &resource.Tag)
	resources, ok := self.resourceTagMap[key]
	if !ok {
		resources = make([]*model.Resource, 0, 5)
	}
	resources = append(resources, resource)
	self.resourceTagMap[key] = resources
}

func (self *MemoryImageCaptchaResourceStore) AddTemplate(captchaType string, template *model.ResourceMap) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if template.Tag == "" {
		template.Tag = common.DEFAULT_TAG
	}
	key := mergeTypeAndTag(captchaType, &template.Tag)
	templates, ok := self.templateResourceTagMap[key]
	if !ok {
		templates = make([]*model.ResourceMap, 0, 2)
	}
	templates = append(templates, template)
	self.templateResourceTagMap[key] = templates
}

func (self *MemoryImageCaptchaResourceStore) ClearTemplate() {
	self.mu.Lock()
	defer self.mu.Unlock()
	if len(self.templateResourceTagMap) == 0 {
		return
	}
	self.templateResourceTagMap = make(map[string][]*model.ResourceMap)
}

func (self *MemoryImageCaptchaResourceStore) ClearResource() {

	self.mu.Lock()
	defer self.mu.Unlock()
	if len(self.resourceTagMap) == 0 {
		return
	}
	self.resourceTagMap = make(map[string][]*model.Resource)
}

func mergeTypeAndTag(captchaType string, tag *string) string {
	var t string
	if tag == nil {
		t = common.DEFAULT_TAG
	} else {
		t = *tag
	}
	return fmt.Sprintf("%s|%s", captchaType, t)
}
