package cfgx

import (
	"os"
	"reflect"

	cfgerr "github.com/its-jojoo/konfiggo/internal/errors"
	"github.com/its-jojoo/konfiggo/internal/parse"
	"github.com/its-jojoo/konfiggo/internal/reflectx"
)

func Load(dst any, opts ...Option) error {
	var o options
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return cfgerr.InvalidTarget("dst must be a non-nil pointer to a struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return cfgerr.InvalidTarget("dst must point to a struct")
	}

	// Walk all fields that declare the `cfg` tag.
	return reflectx.WalkTaggedFields(rv, tagKeyCfg, func(f reflectx.Field) error {
		tags := parseFieldTags(f.StructField)
		if tags.Var == "" {
			return nil
		}

		raw, ok := lookupEnv(o, tags.Var)
		if !ok {
			if tags.HasDef {
				raw = tags.Default
				ok = true
			} else if tags.Required {
				return cfgerr.Required(tags.Var, f.StructField.Name)
			} else {
				// No env, no default, not required: leave zero/current value.
				return nil
			}
		}

		fieldType := f.Value.Type()

		// Parse into the target type.
		var parsed reflect.Value
		var err error

		switch fieldType.Kind() {
		case reflect.Slice:
			parsed, err = parse.ParseSlice(fieldType, raw)
			if err != nil {
				return cfgerr.Parse(tags.Var, f.StructField.Name, fieldType.String(), raw, err)
			}

		default:
			// Special-case time.Duration (int64 under the hood).
			if parse.IsDurationType(fieldType) {
				parsed, err = parse.ParseDuration(raw)
				if err != nil {
					return cfgerr.Parse(tags.Var, f.StructField.Name, fieldType.String(), raw, err)
				}
			} else {
				parsed, err = parse.ParseScalar(fieldType, raw)
				if err != nil {
					// If the scalar parser says unsupported, return a nicer error kind.
					// (We treat it as parse error only when parsing fails for supported types.)
					return cfgerr.Parse(tags.Var, f.StructField.Name, fieldType.String(), raw, err)
				}
			}
		}

		if !parsed.IsValid() {
			return cfgerr.Unsupported(tags.Var, f.StructField.Name, fieldType.String())
		}

		// Set value (parsed already matches target type).
		f.Value.Set(parsed.Convert(fieldType))
		return nil
	})
}

func MustLoad(dst any, opts ...Option) {
	if err := Load(dst, opts...); err != nil {
		panic(err)
	}
}

func lookupEnv(o options, key string) (string, bool) {
	if o.env != nil {
		v, ok := o.env[key]
		return v, ok
	}
	return os.LookupEnv(key)
}
