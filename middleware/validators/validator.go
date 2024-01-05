package validators

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Validations struct {
	Tag string
	Fn  validator.Func
	Msg string
}

var (
	DefaultLanguage = "zh"
	uni             *ut.UniversalTranslator
	trans           ut.Translator
	valid           = []*Validations{
		{"mobile", ValidateMobile, "{0} 手机号码不正确"},
	}

	_Validator = &Validators{}
)

type Validators struct {
	Validate *validator.Validate
}

func (v *Validators) Register(val *Validations) {
	_ = v.Validate.RegisterValidation(val.Tag, val.Fn)
	_ = v.Validate.RegisterTranslation(val.Tag, trans, func(ut ut.Translator) error {
		return ut.Add(val.Tag, val.Msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(val.Tag, fe.Field())
		return t
	})

	// 反馈的信息使用form标签内容
	v.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("form")
	})

	//注册翻译器
	_ = zhtranslations.RegisterDefaultTranslations(v.Validate, trans)
}

func Init() {
	//注册翻译器
	translator := zh.New()
	uni = ut.New(translator, translator)
	trans, _ = uni.GetTranslator(DefaultLanguage)

	// 注册自定义校验器
	_Validator = &Validators{Validate: validator.New()}
	for _, val := range valid {
		_Validator.Register(val)
	}

}

func (v *Validators) Check(params any) map[string][]string {
	err := v.Validate.Struct(params)
	var result = make(map[string][]string)
	if _, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err.(validator.ValidationErrors) {
			result[e.Field()] = append(result[e.StructField()], e.Translate(trans))
		}
		return result
	}
	return nil
}

func Check(params any) map[string][]string {
	return _Validator.Check(params)
}
