package errs

// Error GooError
type Error interface {
	error
	GetStack() string
}