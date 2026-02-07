package reflectx

import (
	"reflect"
)

// Field represents a settable struct field discovered by walking a config struct.
type Field struct {
	// Value is the reflect.Value of the field (settable).
	Value reflect.Value

	// StructField holds tag/name/type metadata.
	StructField reflect.StructField

	// Path is the dotted path of the field inside the struct (e.g. "Server.Port").
	// Useful for debugging/tests. Optional for v0, but cheap to keep.
	Path string
}

// WalkTaggedFields walks a struct value and calls fn for every exported field
// that contains the provided tagKey (e.g. "cfg").
//
// v must be a struct value (not a pointer). It must be addressable/settable
// if you want Field.Value.CanSet() == true.
func WalkTaggedFields(v reflect.Value, tagKey string, fn func(f Field) error) error {
	if v.Kind() != reflect.Struct {
		return nil
	}
	return walk(v, v.Type(), tagKey, "", fn)
}

func walk(v reflect.Value, t reflect.Type, tagKey, prefix string, fn func(f Field) error) error {
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// Skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		fv := v.Field(i)

		// Build path
		path := sf.Name
		if prefix != "" {
			path = prefix + "." + sf.Name
		}

		// If the field is an embedded struct or a named struct, we can recurse.
		// We only recurse into structs that are not time.Time-like special cases.
		// (We'll treat time.Duration as a scalar later; it's not a struct.)
		// Here we keep recursion simple: only into struct kinds.
		if fv.Kind() == reflect.Struct && sf.Anonymous == false && sf.Tag.Get(tagKey) == "" {
			// Recurse into nested structs only if they don't declare cfg tag themselves.
			if err := walk(fv, fv.Type(), tagKey, path, fn); err != nil {
				return err
			}
			continue
		}

		// Handle embedded structs (anonymous fields). Common pattern in configs.
		if sf.Anonymous && fv.Kind() == reflect.Struct && sf.Tag.Get(tagKey) == "" {
			if err := walk(fv, fv.Type(), tagKey, prefix, fn); err != nil {
				return err
			}
			continue
		}

		// If it doesn't have the tagKey, ignore it.
		if _, ok := sf.Tag.Lookup(tagKey); !ok {
			continue
		}

		// Only yield settable fields.
		// (If the parent struct isn't addressable, these won't be settable.)
		if !fv.CanSet() {
			continue
		}

		if err := fn(Field{
			Value:       fv,
			StructField: sf,
			Path:        path,
		}); err != nil {
			return err
		}
	}
	return nil
}
