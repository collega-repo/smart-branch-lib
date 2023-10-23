package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"reflect"
	"repo.collega.co.id/olibs724/smart-branch-service/commons/errs"
	"strconv"
	"strings"
)

var validate *validator.Validate

func NewValidation() error {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get(`json`), `,`, 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})

	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.InexactFloat64()
		}
		return 0
	}, decimal.Decimal{})

	if err := validate.RegisterValidation(`min_digit`, rangeDigitValidator); err != nil {
		return err
	}
	if err := validate.RegisterValidation(`max_digit`, rangeDigitValidator); err != nil {
		return err
	}
	return nil
}

func rangeDigitValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	digit, _ := strconv.Atoi(param)

	field := fl.Field()
	kind := field.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		absValue := field.Int()
		if absValue < 0 {
			absValue = -absValue
		}
		digitValue := len(strconv.FormatInt(absValue, 10))
		switch fl.GetTag() {
		case "min_digit":
			return digitValue >= digit
		case "max_digit":
			return digitValue <= digit
		}
	case reflect.Float32, reflect.Float64:
		absValue := field.Float()
		if absValue < 0 {
			absValue = -absValue
		}
		digitValue := len(strings.Split(strconv.FormatFloat(absValue, []byte(`f`)[0], 2, 64), `.`)[0])
		switch fl.GetTag() {
		case "min_digit":
			return digitValue >= digit
		case "max_digit":
			return digitValue <= digit
		}
	}
	return false
}

func ValidateStruct(i any) errs.ErrMap {
	err := validate.Struct(i)
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		errMap := loopMessage(valErrs)
		if len(errMap) > 0 {
			return errMap
		}
	}
	return nil
}

func loopMessage(valErrs validator.ValidationErrors) errs.ErrMap {
	errMap := make(errs.ErrMap)
	for _, valErr := range valErrs {
		errMap[valErr.Field()] = toMessage(valErr.Tag(), valErr.Param())
	}
	return errMap
}

func toMessage(tag, param string) string {
	tag = strings.ToLower(tag)
	param = strings.ToLower(param)
	switch tag {
	case "required":
		return "Harus Diisi"
	case "min":
		return fmt.Sprintf("Minimum %s karakter", param)
	case "max", "len":
		return fmt.Sprintf("Maksimal %s karakter", param)
	case "min_digit":
		return fmt.Sprintf("Minimum %s digit", param)
	case "max_digit":
		return fmt.Sprintf("Maksimal %s digit", param)
	case "decimal", "number", "numeric":
		return "Harus digit angka"
	case "email", "datetime":
		return "tidak valid"
	default:
		return "tag tidak terdifinisi"
	}
}
