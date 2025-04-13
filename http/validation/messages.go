package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/sjdaws/pkg/common"
)

// processFailures and convert them to sensible error messages.
func (t *Tester) processFailures(failures []validator.FieldError) []string {
	errs := make([]string, 0)

	for _, failure := range failures {
		field := common.FriendlyName(failure.Field())
		tag := failure.Tag()

		switch tag {
		case "email":
			errs = append(errs, fmt.Sprintf("%s '%s' is not a valid email address", field, failure.Value()))
		case "endswith":
			errs = append(errs, fmt.Sprintf("%s must end with '%s'", field, failure.Param()))
		case "oneof":
			errs = append(errs, fmt.Sprintf("%s must be one of: %s", field, common.Options(failure.Param(), "or")))
		case "required":
			errs = append(errs, field+" is required")
		case "uuid":
			errs = append(errs, field+" must be a valid uuid")
		default:
			// Check if tag has been registered with a custom message handler
			if t.messages[tag] != nil {
				errs = append(errs, t.messages[tag](field, failure))

				continue
			}

			errs = append(errs, failure.Error())
		}
	}

	return errs
}
