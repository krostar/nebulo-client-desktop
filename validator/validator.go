package validator

import (
	"fmt"

	gvalidator "github.com/krostar/nebulo-golib/tools/validator"
	validator "gopkg.in/validator.v2"
)

func init() {
	var err error
	// tell to the validator lib that we have some function to use for our custom validators
	if err = validator.SetValidationFunc("file", gvalidator.File); err != nil {
		panic(fmt.Errorf("unable to set validation function %q: %v", "file", err))
	}
	if err = validator.SetValidationFunc("string", gvalidator.String); err != nil {
		panic(fmt.Errorf("unable to set validation function %q: %v", "string", err))
	}
}
