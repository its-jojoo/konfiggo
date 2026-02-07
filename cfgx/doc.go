// Package cfgx provides a tiny, zero-dependency way to load application configuration
// from environment variables into typed Go structs.
//
// Konfiggo follows a simple philosophy:
//
//   - Your struct is the configuration contract
//   - Explicit > magic
//   - Fail fast: the app shouldn't start with invalid config
//   - Human-friendly errors
//
// Basic usage:
//
//	type Config struct {
//		AppName string        `cfg:"APP_NAME" default:"my-app"`
//		Port    int           `cfg:"PORT" default:"8080"`
//		Debug   bool          `cfg:"DEBUG" default:"false"`
//		Timeout time.Duration `cfg:"TIMEOUT" default:"5s"`
//	}
//
//	func main() {
//		var cfg Config
//		cfgx.MustLoad(&cfg)
//	}
//
// Note: The API surface is intentionally small. Use Load for error handling or
// MustLoad for fail-fast apps.
package cfgx
