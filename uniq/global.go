package uniq

var (
	global *Uniq
)

func init() {
	global = New()
}

func Next() string {
	return global.Next()
}

func Check(val string, fullCheck bool) bool {
	return global.Check(val, fullCheck)
}
