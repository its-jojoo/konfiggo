package parse

import (
	"fmt"
	"reflect"
	"strings"
)

// ParseSlice parses a CSV-like string into a slice of the given type.
// Example: "a,b,c" -> []string{"a","b","c"}
// Rules:
//   - Split by comma (,)
//   - Trim spaces around each element
//   - Empty raw string => empty slice
func ParseSlice(sliceType reflect.Type, raw string) (reflect.Value, error) {
	if sliceType.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("expected slice type, got: %s", sliceType.String())
	}

	elemType := sliceType.Elem()

	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		// empty slice
		return reflect.MakeSlice(sliceType, 0, 0), nil
	}

	parts := strings.Split(trimmed, ",")
	out := reflect.MakeSlice(sliceType, 0, len(parts))

	for i, p := range parts {
		itemRaw := strings.TrimSpace(p)
		var item reflect.Value
		var err error

		// Special-case Duration (it's int64 underneath).
		if IsDurationType(elemType) {
			item, err = ParseDuration(itemRaw)
		} else {
			item, err = ParseScalar(elemType, itemRaw)
		}

		if err != nil {
			return reflect.Value{}, fmt.Errorf("element %d: %w", i+1, err)
		}

		out = reflect.Append(out, item.Convert(elemType))
	}

	return out, nil
}
