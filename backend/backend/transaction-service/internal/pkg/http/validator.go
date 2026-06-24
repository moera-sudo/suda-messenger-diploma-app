package http

import (
	"github.com/go-playground/validator/v10"
)

// Validator — реализация echo.Validator поверх go-playground/validator.
// Регистрируется в echo через e.Validator = NewValidator().
type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

func (cv *Validator) Validate(i any) error {
	return cv.v.Struct(i)
}