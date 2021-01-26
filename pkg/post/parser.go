package post

import (
	"errors"
)

func ParseEmailVars(vars []string) (subject, body, asName string, err error) {
	if len(vars) < 2 {
		err = errors.New("invalid len")
		return
	}
	subject = vars[0]
	body = vars[1]
	if len(vars) >= 3 {
		asName = vars[2]
	}
	return
}
