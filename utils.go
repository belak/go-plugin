package plugin

// Ensure is a simple wrapper which will panic if the given error is non-nil.
func Ensure(err error) {
	if err != nil {
		panic(err.Error())
	}
}
