package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func InitTrans(locale string) (err error) {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT)

		var translator ut.Translator
		var found bool

		switch locale {
		case "zh":
			translator, found = uni.GetTranslator("zh")
		case "en":
			translator, found = uni.GetTranslator("en")
		default:
			translator, found = uni.GetTranslator("en")
		}

		if !found {
			return fmt.Errorf("failed to get translator for locale: %s", locale)
		}

		trans = translator

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			label := fld.Tag.Get("label")
			if label != "" {
				return label
			}
			jsonTag := fld.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				return strings.Split(jsonTag, ",")[0]
			}
			return fld.Name
		})

		// 注册翻译
		switch locale {
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(validate, trans)
		case "en":
			err = en_translations.RegisterDefaultTranslations(validate, trans)
		default:
			err = en_translations.RegisterDefaultTranslations(validate, trans)
		}
		return err
	}
	return fmt.Errorf("validator engine is not of type *validator.Validate")
}

func TranslateValidationError(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var messages []string
		for _, e := range errs {
			message := e.Translate(trans)
			messages = append(messages, message)
		}
		return strings.Join(messages, ",")
	}
	return err.Error()
}
