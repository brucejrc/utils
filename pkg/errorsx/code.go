package errorsx

import "net/http"

var (
	OK = &ErrorsX{Code: http.StatusOK, Message: ""}

	ErrInternal = &ErrorsX{Code: http.StatusInternalServerError, Reason: "InternalError", Message: "Internal server error."}

	ErrNotFound = &ErrorsX{Code: http.StatusNotFound, Reason: "NotFound", Message: "ReSource not found."}

	ErrBind = &ErrorsX{Code: http.StatusBadRequest, Reason: "BindError", Message: "Error occurred while binding the request body to the struct."}

	ErrInvalidArgument = &ErrorsX{Code: http.StatusBadRequest, Reason: "InvalidArgument", Message: "Argument verification failed."}

	ErrUnauthenticated = &ErrorsX{Code: http.StatusUnauthorized, Reason: "Unauthenticated", Message: "Unauthenticated."}

	ErrPermissionDenied = &ErrorsX{Code: http.StatusForbidden, Reason: "PermissionDenied", Message: "Permission denied. Access to the requested resource is forbidden."}

	ErrOperationFailed = &ErrorsX{Code: http.StatusConflict, Reason: "OperationFailed", Message: "The requested operation has failed. Please try again later."}
)
