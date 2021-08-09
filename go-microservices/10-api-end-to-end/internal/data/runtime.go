package data

import (
	"encoding/json"
	"errors"
	"strconv"
)

type CanFly bool

func (c CanFly) MarshalJSON() ([]byte, error) {
	var jsonValue string
	if c {
		jsonValue = "yes"
	} else {
		jsonValue = "no"
	}

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

var ErrInvalidCanFlyFormat = errors.New("invalid format for 'can fly'")

func (c *CanFly) UnmarshalJSON(b []byte) error {
	canFlyStr, err := strconv.Unquote(string(b))
	if err != nil {
		return ErrInvalidCanFlyFormat
	}

	switch canFlyStr {
	case "yes":
		*c = true
	case "no":
		*c = false
	default:
		return &json.UnmarshalTypeError{}
	}

	return nil
}
