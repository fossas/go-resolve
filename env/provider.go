package env

import "os"

// Provide reads an environment variable using the Provider pattern.
// Given some variable FOO:
// 1. Look up `$0_PROVIDER`.
// 2. If `$FOO_PROVIDER` is set to some `BAR`, return the value of `$BAR`.
// 3. If `$FOO_PROVIDER` is not set, return the value of `$FOO`.
func Provide(name string) string {
	provider, ok := os.LookupEnv(name + "_PROVIDER")
	if ok {
		return os.Getenv(provider)
	}
	return os.Getenv(name)
}
