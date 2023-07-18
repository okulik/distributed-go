package dist

import "context"

// Circuit type represents a function that takes context and returns a
// pointer to type defined by the user and an error.
type Circuit[T any] func(context.Context) (*T, error)
