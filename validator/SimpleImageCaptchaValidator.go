package validator

import (
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"strconv"
	"strings"
)

const (
	default_tolerant float64 = 0.2
	tolerant_key             = "tolerant"
	name_key                 = "name"
	captcha_type_key         = "captcha_type"
	percentage_key           = "percentage"
)

func NewSimpleImageCaptchaValidator(defaultTolerant *float64) *SimpleImageCaptchaValidator {
	tolerant := default_tolerant
	if defaultTolerant != nil {
		tolerant = *defaultTolerant
	}
	return &SimpleImageCaptchaValidator{
		DefaultTolerant: tolerant,
	}
}

type SimpleImageCaptchaValidator struct {
	DefaultTolerant float64
}

func (self *SimpleImageCaptchaValidator) GenerateImageCaptchaValidData(imageCaptchaInfo *model.ImageCaptchaInfo) (ValidData, error) {
	var result = make(ValidData)
	check, err := self.BeforeGenerateImageCaptchaValidData(imageCaptchaInfo, result)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if check {
		err := self.DoGenerateImageCaptchaValidData(imageCaptchaInfo, result)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	err = self.AfterGenerateImageCaptchaValidData(imageCaptchaInfo, result)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	//panic("implement me")
	return result, nil
}

func (SimpleImageCaptchaValidator) AfterGenerateImageCaptchaValidData(imageCaptchaInfo *model.ImageCaptchaInfo, result ValidData) error {

	return nil
}

func (SimpleImageCaptchaValidator) DoGenerateImageCaptchaValidData(imageCaptchaInfo *model.ImageCaptchaInfo, result ValidData) error {
	if model.TYPE_SLIDER == imageCaptchaInfo.CaptchaType {
		// 滑动类
		percentage := float64(*imageCaptchaInfo.RandomX) / float64(imageCaptchaInfo.BackgroundImageWidth)
		result[percentage_key] = percentage
	} else if model.TYPE_CLICK == imageCaptchaInfo.CaptchaType {
		clickImages := imageCaptchaInfo.Data.Expand.([]*model.ClickImageCheckDefinition)
		var sb string
		for i, img := range clickImages {
			vx := img.X / imageCaptchaInfo.BackgroundImageWidth
			vy := img.Y / imageCaptchaInfo.BackgroundImageHeight
			sb += strconv.Itoa(vx) + "," + strconv.Itoa(vy) + ";"
			if i == 0 && result.ConstantKey(tolerant_key) {
				minLeft := (img.X - img.Width/2) / imageCaptchaInfo.BackgroundImageWidth
				tolerant := vx - minLeft
				result[tolerant_key] = tolerant
			}
		}
		result[percentage_key] = imageCaptchaInfo.CaptchaName
	}
	return nil
}

func (self *SimpleImageCaptchaValidator) BeforeGenerateImageCaptchaValidData(imageCaptchaInfo *model.ImageCaptchaInfo, result ValidData) (bool, error) {
	tolerant := imageCaptchaInfo.Tolerant
	if tolerant != nil && *tolerant > 0 {
		result[tolerant_key] = *tolerant
	} else {
		result[tolerant_key] = self.DefaultTolerant
	}
	result[name_key] = imageCaptchaInfo.CaptchaName
	result["type"] = imageCaptchaInfo.CaptchaName // 兼容java版本
	result[captcha_type_key] = imageCaptchaInfo.CaptchaType
	return true, nil
}

func (self *SimpleImageCaptchaValidator) Valid(imageCaptchaTrack *model.ImageCaptchaTrack, imageCaptchaValidData ValidData) (*model.ApiResponse, error) {

	tolerant, err := imageCaptchaValidData.getFloat(tolerant_key, &self.DefaultTolerant)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// 读验证码名称
	captchaType, err := imageCaptchaValidData.getInt(captcha_type_key, &model.TYPE_SLIDER)

	bgImageWidth := imageCaptchaTrack.BgImageWidth
	if bgImageWidth == nil {
		return nil, errors.New("bgImageWidth is nil")
	}
	trackList := imageCaptchaTrack.TrackList
	if trackList == nil || len(trackList) == 0 {
		return nil, errors.New("trackList is nil")
	}
	var response *model.ApiResponse
	valid, err := self.DoValid(imageCaptchaTrack, imageCaptchaValidData, *tolerant, *captchaType)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if valid {
		// 验证后
		response, err = self.AfterValid(imageCaptchaTrack, imageCaptchaValidData, *tolerant, *captchaType)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	} else {
		apiResponse := model.ApiResponse{}
		response = apiResponse.BasicCheckFail()
	}
	return response, nil
}

func (self *SimpleImageCaptchaValidator) DoValid(imageCaptchaTrack *model.ImageCaptchaTrack,
	imageCaptchaValidData ValidData,
	tolerant float64,
	captchaType int) (bool, error) {
	if model.TYPE_SLIDER == captchaType {
		return self.DoValidSliderCaptcha(imageCaptchaTrack, imageCaptchaValidData, tolerant)
	}
	if model.TYPE_CLICK == captchaType {
		return self.DoValidClickCaptcha(imageCaptchaTrack, imageCaptchaValidData, tolerant)

	}
	return false, nil
}
func (self *SimpleImageCaptchaValidator) DoValidSliderCaptcha(imageCaptchaTrack *model.ImageCaptchaTrack,
	imageCaptchaValidData ValidData,
	tolerant float64) (bool, error) {
	oriPercentage, err := imageCaptchaValidData.getFloat(percentage_key, nil)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	lastTrack := imageCaptchaTrack.TrackList[len(imageCaptchaTrack.TrackList)-1]
	calcPercentage := float64(*lastTrack.X) / float64(*imageCaptchaTrack.BgImageWidth)

	maxTolerant := *oriPercentage + tolerant
	minTolerant := *oriPercentage - tolerant
	check := calcPercentage >= maxTolerant && calcPercentage <= minTolerant
	return check, nil
}
func (self *SimpleImageCaptchaValidator) DoValidClickCaptcha(imageCaptchaTrack *model.ImageCaptchaTrack, imageCaptchaValidData ValidData, tolerant float64) (bool, error) {
	validStr, err := imageCaptchaValidData.getStr(percentage_key, nil)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	splitArr := strings.Split(*validStr, ";")
	trackList := imageCaptchaTrack.TrackList
	// 取出点击事件的轨迹数据
	var clickTrackList []*model.Track
	for _, track := range trackList {
		if strings.EqualFold(common.TRACK_TYPE_CLICK, *track.Type) {
			clickTrackList = append(clickTrackList, &track)
		}
	}
	if len(clickTrackList) != len(splitArr) {
		return false, nil
	}
	for i, pos := range splitArr {
		posArr := strings.Split(pos, ",")
		xPercentage, err := strconv.ParseFloat(posArr[0], 32)
		if err != nil {
			return false, errors.Wrap(err, "")
		}
		yPercentage, err := strconv.ParseFloat(posArr[1], 32)
		if err != nil {
			return false, errors.Wrap(err, "")
		}
		track := clickTrackList[i]
		calcXPercentage := *track.X / float32(*imageCaptchaTrack.BgImageWidth)
		calcYPercentage := *track.Y / float32(*imageCaptchaTrack.BgImageHeight)
		if !checkPercentage(calcXPercentage, float32(xPercentage), float32(tolerant)) || !checkPercentage(calcYPercentage, float32(yPercentage), float32(tolerant)) {
			return false, nil
		}
	}
	return true, nil
}

func (self *SimpleImageCaptchaValidator) AfterValid(track *model.ImageCaptchaTrack, data ValidData, f float64, i int) (*model.ApiResponse, error) {
	response := model.ApiResponse{}
	return response.Success(), nil
}

func checkPercentage(newPercentage float32, oriPercentage float32, tolerant float32) bool {
	maxTolerant := oriPercentage + tolerant
	minTolerant := oriPercentage - tolerant
	if newPercentage >= minTolerant && newPercentage <= maxTolerant {
		return true
	}
	return false

}
