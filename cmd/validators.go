package main

import (
	"slices"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/ishowsagar/go-pizza-shop/internal/models"
)

// acts like a validator for client data
func RegisterCustomValidator() {
	v,ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	// & validation checks
	v.RegisterValidation("pizza_valid_type",createSliceValidator(models.PizzaTypes))
	v.RegisterValidation("pizza_valid_size",createSliceValidator(models.PizzaSizes))
}

// validaion fnc
func createSliceValidator(allowedVals []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return slices.Contains(allowedVals,fl.Field().String())
	}
}