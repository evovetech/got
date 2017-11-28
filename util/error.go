package util

import (
	"fmt"
	"strings"
)

func CompositeError(errors []error) error {
	if l := len(errors); l == 1 {
		return errors[0]
	} else if l > 0 {
		var errStr []string
		for _, err := range errors {
			errStr = append(errStr, err.Error())
		}
		return fmt.Errorf("[%s]", strings.Join(errStr, ", "))
	}
	return nil
}
