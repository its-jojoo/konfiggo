package parse

import (
	"reflect"
	"time"
)

var durationType = reflect.TypeOf(time.Duration(0))

func IsDurationType(t reflect.Type) bool {
	return t == durationType
}

func ParseDuration(raw string) (reflect.Value, error) {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(d), nil
}
