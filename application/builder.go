package application

import (
	"github.com/tianaiyouqing/tianai-captcha-go/generator"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"github.com/tianaiyouqing/tianai-captcha-go/validator"
	"time"
)

// ===================== 构建器实现 ===================

type TianAiCaptchaApplicationBuilder struct {
	ResourceStore         resource.ImageCaptchaResourceStore
	ResourceImageReader   *resource.ImageCaptchaResourceReaders
	ImageTransform        generator.ImageTransform
	providers             map[string]CaptchaGeneratorProvider
	processors            []ImageCaptchaPostProcessor
	CacheStore            CacheStore
	ImageCaptchaValidator validator.ImageCaptchaValidator
}

func NewBuilder() *TianAiCaptchaApplicationBuilder {
	return &TianAiCaptchaApplicationBuilder{
		providers:  make(map[string]CaptchaGeneratorProvider),
		processors: make([]ImageCaptchaPostProcessor, 0),
	}
}

func (self *TianAiCaptchaApplicationBuilder) AddProvider(name string, provider CaptchaGeneratorProvider) {
	self.providers[name] = provider
}

func (self *TianAiCaptchaApplicationBuilder) AddProcessor(processor ImageCaptchaPostProcessor) {
	self.processors = append(self.processors, processor)
}

func (self *TianAiCaptchaApplicationBuilder) SetResourceStore(store resource.ImageCaptchaResourceStore) {
	self.ResourceStore = store
}
func (self *TianAiCaptchaApplicationBuilder) SetImageTransform(transform generator.ImageTransform) {
	self.ImageTransform = transform
}
func (self *TianAiCaptchaApplicationBuilder) SetCacheStore(store CacheStore) {
	self.CacheStore = store
}
func (self *TianAiCaptchaApplicationBuilder) SetImageCaptchaValidator(validator validator.ImageCaptchaValidator) {
	self.ImageCaptchaValidator = validator
}
func (self *TianAiCaptchaApplicationBuilder) SetResourceImageReader(reader *resource.ImageCaptchaResourceReaders) {
	self.ResourceImageReader = reader
}

func (self *TianAiCaptchaApplicationBuilder) Build() *TianAiCaptchaApplication {
	if self.ResourceStore == nil {
		// 设置默认ResourceStore
		self.ResourceStore = resource.NewMemoryImageCaptchaResourceStore()
		// 添加内置资源
		resource.AddDefaultResources(self.ResourceStore)
	}
	if self.ResourceImageReader == nil {
		// 设置默认ResourceImageReader
		self.ResourceImageReader = resource.NewDefaultImageCaptchaResourceReaders()
	}
	if self.ImageTransform == nil {
		// 设置默认ImageTransform
		self.ImageTransform = generator.NewBase64ImageTransform()
	}
	if self.CacheStore == nil {
		// 设置默认CacheStore
		self.CacheStore = NewMemoryCacheStore(5*time.Minute, 5*time.Minute)
	}
	if self.ImageCaptchaValidator == nil {
		// 设置默认ImageCaptchaValidator
		self.ImageCaptchaValidator = validator.NewSimpleImageCaptchaValidator(nil)
	}
	return &TianAiCaptchaApplication{
		CacheStore:            self.CacheStore,
		ImageCaptchaValidator: self.ImageCaptchaValidator,
		ResourceStore:         self.ResourceStore,
		ResourceImageReader:   self.ResourceImageReader,
		ImageTransform:        self.ImageTransform,
		providers:             self.providers,
		processors:            self.processors,
	}
}
