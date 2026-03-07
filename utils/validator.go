package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// 定义一个全局校验器
var validate = validator.New()

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		errs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			//得到所有的错误信息
			errs = append(errs, getFieldErrorMessage(e))
		}
		//只返回一个最前面的错误信息
		return errors.New(errs[0])
	}
	return nil
}

func getFieldErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + "是必填项"
	case "min":
		return e.Field() + "长度不能小于" + e.Param()
	case "max":
		return e.Field() + "长度不能大于" + e.Param()
	case "len":
		return e.Field() + "长度必须是" + e.Param()
	case "alphanum":
		return e.Field() + "只能包含字母和数字"
	case "containsany":
		return e.Field() + "必须包含以下字符之一: " + e.Param()
	case "gte":
		return e.Field() + "不能小于" + e.Param()
	case "lte":
		return e.Field() + "不能大于" + e.Param()
	case "email":
		return "邮箱格式不正确"
	case "oneof":
		return e.Field() + "必须是" + e.Param() + "中的一个"
	case "numeric":
		return e.Field() + "必须是数字"
	case "url":
		return e.Field() + "必须是合法的URL地址"
	default:
		return e.Field() + "验证失败"
	}
}
