package model

import (
	"image"
)

var TYPE_SLIDER = 0
var TYPE_CLICK = 1

type ImageCaptchaInfo struct {
	CaptchaName           string
	CaptchaType           int
	BackgroundImage       string
	TemplateImage         string
	BackgroundImageTag    string
	TemplateImageTag      string
	BackgroundImageWidth  int
	BackgroundImageHeight int
	TemplateImageWidth    int
	TemplateImageHeight   int
	RandomX               *int
	Tolerant              *float64
	Data                  *CustomData
}

type CustomData struct {
	ViewData map[string]any
	Data     map[string]any
	Expand   any
}

func NewResourceMap() *ResourceMap {
	return &ResourceMap{
		ResourceMap: make(map[string]*Resource),
	}
}

type ResourceMap struct {
	Tag         string
	ResourceMap map[string]*Resource
}

func (self *ResourceMap) PutValue(key string, resource *Resource) {
	self.ResourceMap[key] = resource
}

func (self *ResourceMap) Get(key string) *Resource {
	resource := self.ResourceMap[key]
	return resource
}

type Resource struct {
	ResourceType string // 资源类型
	Data         string // 数据
	Tag          string // 标签
	Tip          string // 提示
	Extra        *any   // 扩展
}

func NewResource(resourceType string, data string) *Resource {
	return &Resource{
		ResourceType: resourceType,
		Data:         data,
	}
}

type GenerateParam struct {
	Obfuscate          bool
	CaptchaName        string
	BackgroundImageTag *string
	TemplateImageTag   *string
	Param              map[string]any
}

type CaptchaExchange struct {
	TemplateResource *ResourceMap
	ResourceImage    *Resource
	BgImage          image.Image
	TemplateImage    image.Image
	CustomData       *CustomData
	Param            *GenerateParam
	TransferData     any
}
type ClickImageCheckDefinition struct {
	X      int
	Y      int
	Width  int
	Height int
	Value  string
}

type ImageCaptchaVO struct {
	CaptchaName           string `json:"type,name"`
	BackgroundImage       string `json:"backgroundImage"`
	TemplateImage         string `json:"templateImage"`
	BackgroundImageTag    string `json:"backgroundImageTag"`
	TemplateImageTag      string `json:"templateImageTag"`
	BackgroundImageWidth  int    `json:"backgroundImageWidth"`
	BackgroundImageHeight int    `json:"backgroundImageHeight"`
	TemplateImageWidth    int    `json:"templateImageWidth"`
	TemplateImageHeight   int    `json:"templateImageHeight"`
	Data                  any    `json:"data"`
}

type ImageCaptchaTrack struct {
	BgImageWidth        *int    `json:"bgImageWidth"`
	BgImageHeight       *int    `json:"bgImageHeight"`
	TemplateImageWidth  *int    `json:"templateImageWidth"`
	TemplateImageHeight *int    `json:"templateImageHeight"`
	StartTime           *int64  `json:"startTime"`
	StopTime            *int64  `json:"stopTime"`
	StartLeft           *int    `json:"startLeft"`
	StartTop            *int    `json:"startTop"`
	TrackList           []Track `json:"trackList"`
	Data                any     `json:"data"`
}

type Track struct {
	X    *float32 `json:"x"`
	Y    *float32 `json:"y"`
	T    *float32 `json:"t"`
	Type *string  `json:"type"`
}

type ApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (self *ApiResponse) Success() *ApiResponse {
	self.Code = 200
	self.Msg = "OK"
	return self
}
func (self *ApiResponse) Expire() *ApiResponse {
	self.Code = 4000
	self.Msg = "已失效"
	return self
}
func (self *ApiResponse) BasicCheckFail() *ApiResponse {
	self.Code = 4001
	self.Msg = "基础校验失败"
	return self
}
