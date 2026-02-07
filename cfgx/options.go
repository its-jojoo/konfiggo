package cfgx

// Option configures Load behavior.
// The v0 API is intentionally small; options will be added gradually.
type Option func(*options)

type options struct{}
