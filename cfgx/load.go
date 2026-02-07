package cfgx

// Load loads configuration into dst from the environment, applying tags such as
// `cfg`, `default`, and `required`.
//
// dst must be a pointer to a struct.
//
// In v0, sources are:
//   1) Environment variables
//   2) Declarative defaults (tag `default`)
//   3) Zero values (if neither env nor default is present)
func Load(dst any, opts ...Option) error {
	// Implementation will be added in upcoming commits.
	return nil
}

// MustLoad is like Load but panics on error.
// It is intended for fail-fast applications where invalid config must stop startup.
func MustLoad(dst any, opts ...Option) {
	if err := Load(dst, opts...); err != nil {
		panic(err)
	}
}
