<div align="center">

![][image-logo]

### tianaiCAPTCHA - 天爱验证码
#### 基于 Golang 实现的滑动/旋转/点选 验证码  **[在线体验][online-demo-link]**
<sup>可能是开源界最好用的行为验证码工具</sup><br>

[![][github-release-shield]][github-release-link] [![][github-license-shield]][github-license-link]<br>

</div>

---

- **tianai CAPTCHA Golang 版** 目前含有的的行为验证码类型
<table>
  <tr>
    <th>滑块验证码</th>
    <th>旋转验证码</th>
    <th>文字点选验证码</th>
  </tr>
  <tr>
    <td><img src="https://minio.tianai.cloud/public/demo-view/go-slider-1.png" alt="滑块验证码"></td>
    <td><img src="https://minio.tianai.cloud/public/demo-view/go-rotate-1.png" alt="旋转验证码"></td>
    <td><img src="https://minio.tianai.cloud/public/demo-view/go-click-1.png" alt="文字点选验证码"></td>
  </tr>
  <tr>
    <td colspan="3" align="center"><a href="http://captcha.tianai.cloud">在线体验 🚀</a></td>
  </tr>
</table>

> 更多验证码类型待开发，敬请期待! ✨

> [!NOTE]
> 功能更丰富的 Java 版：[github.com/dromara/tianai-captcha][tianai-captcha-java-link]

# 快速开始

> [!TIP]
> Demo 参考 [gitee.com/tianai/captcha-go-demo][captcha-go-demo-link]

## 1. mod 导入
```shell
go get github.com/tianaiyouqing/tianai-captcha-go@v1.0.1
```
> #### 需要本地开发，可拉取至项目根目录，并在 `go.mod` 添加
> ```mod
> replace github.com/tianaiyouqing/tianai-captcha-go => ./tianai-captcha-go
> ```

## 2. 初始化验证码

<details>
  <summary>展开</summary>

```go
import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
)

var Captcha *application.TianAiCaptchaApplication

func init() {
    builder := application.NewBuilder()
    // 添加滑块验证码
    builder.AddProvider(application.CreateSliderProvider())
    // 添加旋转验证码
    builder.AddProvider(application.CreateRotateProvider())
    // 添加文字点选验证码， 参数为nil时会读取默认的字体，可以替换成自定义字体， 传入多个字体会随机选择
    builder.AddProvider(application.CreateWordClickProvider(nil))
	// 注意 构建出来的 CaptchaApplication 是单例的，所以可以全局使用
    Captcha = builder.Build()
}
```
</details>

## 3. 在项目中使用

<details>
  <summary>展开</summary>

```go
// 这里以gin框架为例，其它框架自行修改即可

// 生成验证码
func GenCaptcha(c *gin.Context) {
    // 这里生成类型为 SLIDER的验证码， 目前支持 SLIDER、ROTATE、WORD_IMAGE_CLICK
    captcha, err := Captcha.GenerateCaptcha(&model.GenerateParam{
        CaptchaName: "SLIDER",
    })
    if err != nil {
        c.JSON(500, gin.H{
            "code": 500,
            "msg":  err.Error(),
        })
        return
    }
    // 这边返回的结构是为了适配 tianai-captcha-web-sdk 前端项目
    c.JSON(200, gin.H{
        "code":    200,
        "msg":     "success",
        "id":      captcha.Id,
        "captcha": captcha,
    })

}

// 校验验证码
func Valid(c *gin.Context) {
	
    // 该接收的参数结构是前端项目 tianai-captcha-web-sdk 的 ValidParam
    validParam := new(ValidParam)
    if err := c.ShouldBindJSON(validParam); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    valid, err := Captcha.Valid(validParam.Id, &validParam.Data)
    if err != nil {
        c.JSON(500, gin.H{
        "code": 500,
        "msg":  err.Error(),
        })
        return
    }
    c.JSON(200, valid)
}

type ValidParam struct {
    Id   string                  `json:"id" binding:"required"`
    Data model.ImageCaptchaTrack `json:"data" binding:"required"`
}

```
</details>

---

> [!NOTE]
> 默认前端 SDK<br>
> [github.com/tianaiyouqing/captcha-web-sdk][tianai-captcha-web-sdk-link]


## 扩展

<details>
  <summary>自定义背景图片和模板图片</summary>

```go
package controller

import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
)

var Captcha *application.TianAiCaptchaApplication

func init() {

	// 创建一个resourceStore存储
	store := resource.NewMemoryImageCaptchaResourceStore()

	// ==================== 设置背景图片 ====================

	// 注意: 背景图片的宽高为 600*360， 后台支持自定义调整宽高， 但是 tianai-captcha-web-sdk 项目设置的样式必须为 600*360的比例， 可以缩放大小，但是小白的话直接就设置600*360

	// 第一个参数为验证码名称，第二个参数为资源, 验证码名称目前支持 "SLIDER" "ROTATE" "WORD_IMAGE_CLICK"
	store.AddResource(common.CAPTCHA_NAME_SLIDER, &model.Resource{
		ResourceType: "file",              // 这个参数指定资源类型，目前支持file、url， 可以自己扩展相对应的资源读取器，这里演示file类型
		Data:         "./rsources/1.jpeg", // 如果类型为file， 则data为文件路径， 如果类型为url， 则data为url地址， 以此类推
	})
	store.AddResource(common.CAPTCHA_NAME_ROTATE, &model.Resource{
		ResourceType: "file",
		Data:         "./rsources/2.jpeg",
	})
	store.AddResource(common.CAPTCHA_NAME_SLIDER, &model.Resource{
		ResourceType: "file",
		Data:         "./rsources/3.jpeg",
	})
	store.AddResource(common.CAPTCHA_NAME_SLIDER, &model.Resource{
		ResourceType: "file",
		Data:         "./rsources/4.jpegg",
	})

	// ============================ 设置模板图片 ============================
	// 注意模板图片要按照指定的格式设置样式和尺寸， 下面给出默认大小，该宽高是按照背景图为 600*360计算的， 如果模板图片不符合要求， 则需要自己调整样式和尺寸
	// slider 底图样式大小为 110*110, 滑块样式大小为 110*110
	// rotate 底图样式大小为 200*200, 旋转图片样式大小为 200*200

	// ============= 设置自定义模板 ====================
	resourceMap := model.NewResourceMap()
	resourceMap.PutValue("active.png", &model.Resource{
		ResourceType: "file",
		Data:         "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-滑块.png",
	})
	resourceMap.PutValue("fixed.png", &model.Resource{
		ResourceType: "file",
		Data:         "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-底图.png",
	})
	store.AddTemplate(common.CAPTCHA_NAME_SLIDER, resourceMap)

	// ============= 设置默认模板 ====================
	// 如果小白只想替换背景图，不想设置模板，也可以使用系统自带的模板，代码如下
	defaultSliderTemplates := resource.GetDefaultSliderTemplates()
	for _, template := range defaultSliderTemplates {
		store.AddTemplate(common.CAPTCHA_NAME_SLIDER, template)
	}
	defaultRotateTemplates := resource.GetDefaultRotateTemplate()
	for _, template := range defaultRotateTemplates {
		store.AddTemplate(common.CAPTCHA_NAME_ROTATE, template)
	}

	// ========= 构建验证码应用 ==========

	builder := application.NewBuilder()
	// 设置存放资源的存储器
	builder.SetResourceStore(store)
	// 添加验证码生成器
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	Captcha = builder.Build()
}
```
</details>

<details>
  <summary>验证码校验信息的存储位置</summary>

#### 验证码的校验信息默认存储在内存中，若要换成 Redis 之类的，自定义扩展即可，例子：

```go
package controller

import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
)

var Captcha *application.TianAiCaptchaApplication

func init() {
	// 自定义缓存存储器
	customCacheStore := &CustomCacheStore{}

	builder := application.NewBuilder()
	// 设置自定义缓存存储器
	builder.SetCacheStore(customCacheStore)
	// 添加验证码生成器
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	Captcha = builder.Build()
}

type CustomCacheStore struct{}

func (CustomCacheStore) GetCache(key string) (value map[string]any, ok bool) {
	//TODO 通过key获取缓存
	panic("implement me")
}

func (CustomCacheStore) GetAndRemoveCache(key string) (value map[string]any, ok bool) {
	//TODO 通过key获取缓存并删除
	panic("implement me")
}

func (CustomCacheStore) SetCache(key string, data map[string]any, captchaInfo *model.ImageCaptchaInfo) error {
	//TODO 设置缓存
	panic("implement me")
}
```
</details>

<details>
  <summary>自定义图片转换器</summary>

```go
package controller

import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/generator"
)

var Captcha *application.TianAiCaptchaApplication

func init() {
	// 自定义缓存存储器


	builder := application.NewBuilder()
	// 设置自定义图片转换器， 默认是base64格式的转换前， 背景图为 jpg， 模板图为png， 如有需要可自定义实现 `generator.ImageTransform` 接口进行转换
	builder.SetImageTransform(generator.NewBase64ImageTransform())
	// 添加验证码生成器
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	Captcha = builder.Build()
}
```
</details>

<details>
  <summary>其他</summary>

```go
package controller

import (
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/generator"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"github.com/tianaiyouqing/tianai-captcha-go/validator"
	"time"
)

var Captcha *application.TianAiCaptchaApplication

func init() {
	builder := application.NewBuilder()

	// 设置资源存储器
	builder.SetResourceStore(resource.NewMemoryImageCaptchaResourceStore())
	// 设置资源读取器，负责把Resource对象转换成Image图片对象
	readers := resource.NewDefaultImageCaptchaResourceReaders()
	//readers.AddResourceReader(nil)// 自定义可以添加自定义的资源读取器
	builder.SetResourceImageReader(readers)
	// 设置自定义图片转换器， 默认是base64格式的转换前， 背景图为 jpg， 模板图为png， 如有需要可自定义实现 `generator.ImageTransform` 接口进行转换
	builder.SetImageTransform(generator.NewBase64ImageTransform())
	// 设置缓冲存储器， 默认是内存存储器， 如需要扩展redis之类， 可自定义实现 `application.CacheStore` 接口
	builder.SetCacheStore(application.NewMemoryCacheStore(5*time.Minute, 5*time.Minute))
	// 设置验证码验证器 参数为默认的容错值，传nil则容错值默认设置为 0.02
	builder.SetImageCaptchaValidator(validator.NewSimpleImageCaptchaValidator(nil))

	// 添加验证码生成器
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	Captcha = builder.Build()
}
```
</details>

---

> [!TIP]
> ### 👋 加入我们，与其他爱好者交流学习
>
> QQ群：`1021884609`<br>
>
> 微信群：群满，可扫码添加好友拉入群聊，或搜索微信号: `youseeseeyou-1ttd` 备注 `验证码`
>
> ![][qrcode-link] <br>

---

## 贡献者

<a href="https://github.com/tianaiyouqing/tianai-captcha-go/graphs/contributors" target="_blank">
  <table>
    <tr>
      <th>
          <img src="https://contrib.rocks/image?repo=tianaiyouqing/tianai-captcha-go">
      </th>
    </tr>
  </table>
</a>

---

[image-logo]: https://minio.tianai.cloud/public/captcha/logo/logo-519x100.png
[github-release-shield]: https://img.shields.io/github/v/release/tianaiyouqing/tianai-captcha-go?color=369eff&labelColor=black&logo=github&style=flat-square
[github-release-link]: https://github.com/tianaiyouqing/tianai-captcha-go/releases
[github-license-link]: https://github.com/tianaiyouqing/tianai-captcha-go/blob/master/LICENSE
[github-license-shield]: https://img.shields.io/badge/MulanPSL-2.0-white?labelColor=black&style=flat-square
[tianai-captcha-java-link]: https://github.com/dromara/tianai-captcha
[captcha-go-demo-link]: https://gitee.com/tianai/captcha-go-demo
[tianai-captcha-web-sdk-link]: https://github.com/tianaiyouqing/captcha-web-sdk
[online-demo-link]: http://captcha.tianai.cloud
[qrcode-link]: https://minio.tianai.cloud/public/qun4.png