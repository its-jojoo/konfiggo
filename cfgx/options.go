package cfgx

// Option configures Load behavior.
type Option func(*options)

type options struct {
	// env allows injecting environment values (useful for tests).
	// If nil, Load reads from the process environment.
	env map[string]string
}

// WithEnv injects an environment map (mainly for tests).
// When set, Load will read from this map instead of os.LookupEnv.
func WithEnv(env map[string]string) Option {
	return func(o *options) {
		o.env = env
	}
}
