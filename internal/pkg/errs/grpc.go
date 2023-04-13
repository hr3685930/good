package errs

import (
	"good/pkg/rpc"
	"google.golang.org/grpc/codes"
)

// Canceled 超时
func Canceled(msg string) error {
	return rpc.Err(codes.Canceled, msg)
}

// Unknown 系统错误
func Unknown(msg string) error {
	return rpc.Err(codes.Unknown, msg)
}

// InvalidArgument 参数失败
func InvalidArgument(msg string) error {
	return rpc.Err(codes.InvalidArgument, msg)
}

// OutOfRange OutOfRange
func OutOfRange(msg string) error {
	return rpc.Err(codes.OutOfRange, msg)
}

// NotFound not fount
func NotFound(msg string) error {
	return rpc.Err(codes.NotFound, msg)
}

// AlreadyExists AlreadyExists
func AlreadyExists(msg string) error {
	return rpc.Err(codes.AlreadyExists, msg)
}

// PermissionDenied authorization failed
func PermissionDenied(msg string) error {
	return rpc.Err(codes.PermissionDenied, msg)
}

// Unauthenticated authentication failed
func Unauthenticated(msg string) error {
	return rpc.Err(codes.Unauthenticated, msg)
}

// ResourceExhausted ResourceExhausted
func ResourceExhausted(msg string) error {
	return rpc.Err(codes.ResourceExhausted, msg)
}

// Unavailable Unavailable
func Unavailable(msg string) error {
	return rpc.Err(codes.Unavailable, msg)
}

// Aborted Aborted
func Aborted(msg string) error {
	return rpc.Err(codes.Aborted, msg)
}

// FailedPrecondition FailedPrecondition
func FailedPrecondition(msg string) error {
	return rpc.Err(codes.FailedPrecondition, msg)
}
