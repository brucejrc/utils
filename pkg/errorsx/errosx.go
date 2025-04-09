package errorsx

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrorsX struct {
	Code     int               `json:"code,omitempty"`
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func New(code int, reason string, format string, args ...any) *ErrorsX {
	return &ErrorsX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

func (err *ErrorsX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", err.Code, err.Reason, err.Message, err.Metadata)
}

func (err *ErrorsX) WithMessage(format string, args ...any) *ErrorsX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

func (err *ErrorsX) WithMetadata(md map[string]string) *ErrorsX {
	err.Metadata = md
	return err
}

func (err *ErrorsX) KV(kvs ...string) *ErrorsX {
	if err.Metadata == nil {
		err.Metadata = make(map[string]string)
	}

	for i := 0; i < len(kvs); i = i + 2 {
		if i < len(kvs)-1 {
			err.Metadata[kvs[i]] = kvs[i+1]
		}
	}
	return err
}

func (err *ErrorsX) GrpcStatus() *status.Status {
	details := errdetails.ErrorInfo{Reason: err.Reason, Metadata: err.Metadata}
	s, _ := status.New(httpstatus.ToGRPCCode(err.Code), err.Message).WithDetails(&details)
	return s
}

func (err *ErrorsX) WithRequestID(requestID string) *ErrorsX {
	return err.KV("x-request-id", requestID)
}

func (err *ErrorsX) Is(target error) bool {
	if errx := new(ErrorsX); errors.As(target, &errx) {
		return errx.Code == err.Code && errx.Reason == err.Reason
	}
	return false
}

func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return FromError(err).Code
}

func Reason(err error) string {
	if err == nil {
		return ""
	}
	return FromError(err).Reason
}

func FromError(err error) *ErrorsX {
	if err == nil {
		return nil
	}

	if errx := new(ErrorsX); errors.As(err, &errx) {
		return errx
	}

	gs, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, gs.Message())
	}

	ret := New(httpstatus.FromGRPCCode(gs.Code()), ErrInternal.Reason, gs.Message())

	for _, detail := range gs.Details() {
		if typed, ok := detail.(*errdetails.ErrorInfo); ok {
			ret.Reason = typed.Reason
			return ret.WithMetadata(typed.Metadata)
		}
	}
	return ret
}
