package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WrapError wraps an error with a specific gRPC code and message.
func WrapError(code codes.Code, err error) error {
	return status.Errorf(code, "%s: %v", code.String(), err)
}

// New wraps an error with a specific gRPC code and message.
// It's a shorthand for WrapError.
func New(code codes.Code, format string, args ...interface{}) error {
	return WrapError(code, fmt.Errorf(format, args...))
}

// ExtractCode extracts the gRPC code from an error.
func ExtractCode(err error) codes.Code {
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}
	return codes.Unknown
}
