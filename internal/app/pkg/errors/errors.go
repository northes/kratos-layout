package errors

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code     ErrorCode
	Message  string
	Metadata map[string]string
}

func New(code ErrorCode, msg string, ops ...Option) *Error {
	err := &Error{
		Code:    code,
		Message: msg,
	}

	for _, option := range ops {
		option(err)
	}

	return err
}

func (e Error) Error() string {
	errStr := fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)

	if len(e.Metadata) > 0 {
		kvs := make([]string, 0, len(e.Message))
		for k, v := range e.Metadata {
			kvs = append(kvs, fmt.Sprintf("\"%s\": \"%s\"", k, v))
		}

		errStr = fmt.Sprintf("%s, metadata: {%s}", errStr, strings.Join(kvs, ", "))
	}

	return errStr
}

func (e Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Metadata: e.Metadata,
		})
	return s
}

type Option func(e *Error)

func WithMetadata(metadata map[string]string) Option {
	return func(e *Error) {
		e.Metadata = metadata
	}
}

func FromGRPCError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(ErrorCode(gs.Code()), gs.Message())
		for _, details := range gs.Details() {
			switch d := details.(type) {
			case *errdetails.ErrorInfo:
				ret.Metadata = d.Metadata
			}
		}
		return ret
	}

	return New(ServerErrorCode, err.Error())

}
