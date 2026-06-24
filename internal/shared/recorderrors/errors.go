package recorderrors

import "errors"

var ErrInvalidInput = errors.New("invalid input")
var ErrUnauthenticated = errors.New("unauthenticated")
var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrFailedPrecondition = errors.New("failed precondition")
