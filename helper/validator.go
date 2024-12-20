package helper

import (
	"project/domain"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NotEmptySlice(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().([]domain.OrderItem)
	return ok && len(val) > 0
}

// NewValidator membuat instance Validator dengan registrasi validator kustom
func NewValidator() *Validator {
	v := validator.New()
	_ = v.RegisterValidation("notempty", NotEmptySlice)
	return &Validator{validate: v}
}

// ValidateStruct melakukan validasi terhadap struct yang diberikan
func (v *Validator) ValidateStruct(data interface{}) error {
	return v.validate.Struct(data) // Validasi menggunakan rules yang telah didaftarkan
}

// FormatValidationError mengubah error validasi menjadi pesan yang lebih ramah
func FormatValidationError(err error) string {
	errorMessages := map[string]string{
		"Name_required":        "Name is required",
		"Name_min":             "Name must have at least 3 characters",
		"Description_required": "Description is required",
		"Description_min":      "Description must have at least 20 characters",
		"Quantity_gt":          "Quantity must be greater than 0",
		"ProductID_required":   "ProductID is required",
		"OrderItems_required":  "Order item cannot nil",
		"OrderItems_notempty":  "Order item cannot empty",
		"StatusPayment_oneof":  "Status payment must be 'In Process', 'Completed' or 'Cancelled'",
		"StatusKitchen_oneof":  "Status kitchen must be 'In The Kitchen', 'Cooking Now', or 'Ready To Serve'",
	}

	var errMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {

			key := ve.Field() + "_" + ve.Tag()

			if message, found := errorMessages[key]; found {
				errMessages = append(errMessages, message)
			} else {

				errMessages = append(errMessages, ve.Field()+" is invalid: "+ve.Tag())
			}
		}
	}

	return strings.Join(errMessages, ", ")
}
