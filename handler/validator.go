package handler

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v

		// 注册自定义验证器
		_ = v.RegisterValidation("order_status", validateOrderStatus)

		// 初始化翻译器
		zhTrans := zh.New()
		uni = ut.New(zhTrans, zhTrans)
		trans, _ = uni.GetTranslator("zh")

		// 注册翻译器
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)

		// 注册自定义翻译
		registerCustomTranslations(validate, trans)

		// 注册一个函数，获取struct tag中的label作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			if name == "" {
				name = fld.Tag.Get("json")
			}
			return name
		})

		// 注册自定义错误信息
		registerCustomErrorMessages(validate, trans)
	}
}

func validateOrderStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}
	return validStatuses[status]
}

func registerCustomTranslations(v *validator.Validate, trans ut.Translator) {
	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "order_status",
			translation: "{0}必须是有效的订单状态",
			override:    false,
		},
		{
			tag:         "required",
			translation: "{0}不能为空",
			override:    true,
		},
		{
			tag:         "gt",
			translation: "{0}必须大于{1}",
			override:    true,
		},
		{
			tag:         "gte",
			translation: "{0}必须大于或等于{1}",
			override:    true,
		},
		{
			tag:         "email",
			translation: "{0}必须是有效的电子邮件地址",
			override:    true,
		},
		{
			tag:         "len",
			translation: "{0}长度必须等于{1}",
			override:    true,
		},
		{
			tag:         "min",
			translation: "{0}长度必须大于或等于{1}",
			override:    true,
		},
		{
			tag:         "max",
			translation: "{0}长度必须小于或等于{1}",
			override:    true,
		},
	}

	for _, t := range translations {
		if t.customTransFunc != nil && t.customRegisFunc != nil {
			_ = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
		} else {
			_ = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}
	}
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.Error()
	}
	return t
}

func registerCustomErrorMessages(v *validator.Validate, trans ut.Translator) {
	messages := []struct {
		tag     string
		message string
	}{
		{"required_with", "{0}为必填项"},
		{"numeric", "{0}必须是数字"},
		{"oneof", "{0}必须是[{1}]中的一个"},
	}

	for _, m := range messages {
		registerFn := func(ut ut.Translator) error {
			return ut.Add(m.tag, m.message, false)
		}

		transFn := func(ut ut.Translator, fe validator.FieldError) string {
			param := fe.Param()
			if m.tag == "oneof" {
				param = strings.Join(strings.Split(param, " "), ",")
			}
			t, err := ut.T(fe.Tag(), fe.Field(), param)
			if err != nil {
				return fe.Error()
			}
			return t
		}

		_ = v.RegisterTranslation(m.tag, trans, registerFn, transFn)
	}
}

// GetValidationErrors 获取验证错误的中文信息
func GetValidationErrors(err error) []string {
	var errors []string
	validationErrors := err.(validator.ValidationErrors)
	
	for _, e := range validationErrors {
		errors = append(errors, e.Translate(trans))
	}
	
	return errors
}

// 订单状态的中文描述
var OrderStatusMap = map[string]string{
	"pending":   "待支付",
	"paid":      "已支付",
	"shipped":   "已发货",
	"delivered": "已送达",
	"cancelled": "已取消",
}