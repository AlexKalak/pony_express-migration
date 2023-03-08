package validation

import (
	"regexp"

	"github.com/alexkalak/pony_express/src/types"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func GetValidator() *validator.Validate {
	if validate != nil {
		return validate
	}
	validate = validator.New()
	validate.RegisterValidation("only-letters-and-spaces", ValidateOnlyLettersAndSpaces)
	return validate
}

func Validate(structure interface{}) []*types.ErrorResponse {
	var errors []*types.ErrorResponse

	err := GetValidator().Struct(structure)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element types.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateOnlyLettersAndSpaces(fl validator.FieldLevel) bool {
	valid, _ := regexp.MatchString(`^[A-Za-zА-Яа-яÇĞIİÖŞÜçğii̇öşü\s]+$`, fl.Field().String())

	return valid
}
