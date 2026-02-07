package cfgx

// Option configures Load behavior.
// The v1 API is intentionally small; options will be added gradually.
type Option func(*options)

type options struct{}
