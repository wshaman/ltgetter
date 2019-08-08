package tools

// OnErrPanic panics if err. Just a sugar
func OnErrPanic( err error) {
	if err != nil {
		panic(err)
	}
}
