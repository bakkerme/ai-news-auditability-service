package customerrors

// ErrorString represents a string value and an associated error
type ErrorString struct {
	Value string
	Err   error
}
