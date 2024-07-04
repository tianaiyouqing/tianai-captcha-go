package test

import (
	"fmt"
	"github.com/tianaiyouqing/tianai-captcha-go/application"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"testing"
	"time"
)

func TestGen2(t *testing.T) {
	store := resource.NewMemoryImageCaptchaResourceStore()
	store.AddResource("slider", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\codeing\\captcha-server-demo\\src\\main\\resources\\bgimages\\a.jpg"))
	resourceMap := model.NewResourceMap()
	resourceMap.PutValue("active.png", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-滑块.png"))
	resourceMap.PutValue("fixed.png", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-底图.png"))
	store.AddTemplate("slider", resourceMap)

	store.AddResource("rotate", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\codeing\\captcha-server-demo\\src\\main\\resources\\bgimages\\a.jpg"))
	resourceMap = model.NewResourceMap()
	resourceMap.PutValue("active.png", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\active.png"))
	resourceMap.PutValue("fixed.png", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\fixed.png"))
	store.AddTemplate("rotate", resourceMap)

	store.AddResource("word_click", model.NewResource("file", "C:\\Users\\Thinkpad\\Desktop\\captcha\\codeing\\captcha-server-demo\\src\\main\\resources\\bgimages\\a.jpg"))

	builder := application.NewBuilder()
	//builder.SetResourceStore(store)
	builder.AddProvider(application.CreateSliderProvider())
	builder.AddProvider(application.CreateRotateProvider())
	builder.AddProvider(application.CreateWordClickProvider(nil))
	captcha := builder.Build()

	//for i := 0; i < 100; i++ {
	// 计算耗时
	start := time.Now()
	res, err := captcha.GenerateCaptcha(&model.GenerateParam{
		CaptchaName: "word_click",
	})
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(res)
	fmt.Printf("耗时：%v\n", time.Since(start).Milliseconds())
	//}
}
func TestGen(t *testing.T) {
	store := resource.NewMemoryImageCaptchaResourceStore()

	store.AddResource("slider", &model.Resource{
		ResourceType: "file",
		Data:         "C:\\Users\\Thinkpad\\Desktop\\captcha\\codeing\\captcha-server-demo\\src\\main\\resources\\bgimages\\a.jpg",
	})
	resourceMap := model.NewResourceMap()
	resourceMap.PutValue("active.png", &model.Resource{
		ResourceType: "file",
		Data:         "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-滑块.png",
	})
	resourceMap.PutValue("fixed.png", &model.Resource{
		ResourceType: "file",
		Data:         "C:\\Users\\Thinkpad\\Desktop\\captcha\\templates\\六边形-底图.png",
	})
	store.AddTemplate("slider", resourceMap)
	builder := application.NewBuilder()
	builder.SetResourceStore(store)
	builder.AddProvider(application.CreateSliderProvider())
	captcha := builder.Build()
	for i := 0; i < 100; i++ {
		// 计算耗时
		start := time.Now()
		_, err := captcha.GenerateCaptcha(&model.GenerateParam{
			CaptchaName: "slider",
		})
		if err != nil {
			fmt.Printf("%v", err)
		}
		fmt.Printf("耗时：%v\n", time.Since(start).Milliseconds())
	}

}
