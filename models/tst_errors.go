package models

// CustomError defines a custom error type with status code and message
type CustomError struct {
	StatusCode int
	Message    string
}

// Error Interface, implements the error interface for CustomError
func (e *CustomError) Error() string {
	return e.Message
}

// ?
// This method is defined with a receiver of type *CustomError (a pointer to CustomError).
// This means the method can be called on instances of CustomError
// Usage (create an instance of CustomError)
//   err := &CustomError{
// 		StatusCode: 404,
// 		Message:    "Resource not found",
// }
