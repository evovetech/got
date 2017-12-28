package util

import (
	"fmt"
	"strings"
)

func CompositeError(errors []error) error {
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	default:
		var err = make([]string, len(errors))
		for i, e := range errors {
			err[i] = e.Error()
		}
		return fmt.Errorf("[%s]", strings.Join(err, ", "))
	}
}
