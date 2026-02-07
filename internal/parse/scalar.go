package parse

import (
	"fmt"
	"reflect"
	"strconv"
)

// ParseScalar parses a single (non-slice) value from raw string into the given type.
// Supported kinds: string, bool, int*, uint*, float*.
// time.Duration is handled in duration.go via ParseDuration.
//
// It returns the parsed reflect.Value with the requested type.
func ParseScalar(typ reflect.Type, raw string) (reflect.Value, error) {
	switch typ.Kind() {
	case reflect.String:
		return reflect.ValueOf(raw).Convert(typ), nil

	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(typ), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// NOTE: time.Duration is int64 but should be parsed by time.ParseDuration.
		// Callers must special-case Duration before calling ParseScalar.
		bitSize := typ.Bits()
		v, err := strconv.ParseInt(raw, 10, bitSize)
		if err != nil {
			return reflect.Value{}, err
		}
		out := reflect.New(typ).Elem()
		out.SetInt(v)
		return out, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bitSize := typ.Bits()
		v, err := strconv.ParseUint(raw, 10, bitSize)
		if err != nil {
			return reflect.Value{}, err
		}
		out := reflect.New(typ).Elem()
		out.SetUint(v)
		return out, nil

	case reflect.Float32, reflect.Float64:
		bitSize := typ.Bits()
		v, err := strconv.ParseFloat(raw, bitSize)
		if err != nil {
			return reflect.Value{}, err
		}
		out := reflect.New(typ).Elem()
		out.SetFloat(v)
		return out, nil

	default:
		return reflect.Value{}, fmt.Errorf("unsupported scalar type: %s", typ.String())
	}
}
