// Package ptr provides pointer helper functions.
package ptr

// String returns a pointer to a string.
func String(s string) *string { return &s }

// Uint64 returns a pointer to a uint64.
func Uint64(n uint64) *uint64 { return &n }

// Int returns a pointer to an int.
func Int(n int) *int { return &n }

// Bool returns a pointer to a bool.
func Bool(b bool) *bool { return &b }
