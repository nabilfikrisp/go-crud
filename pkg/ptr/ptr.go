package ptr

func String(s string) *string { return &s }
func Uint64(n uint64) *uint64 { return &n }
func Int(n int) *int          { return &n }
func Bool(b bool) *bool       { return &b }
