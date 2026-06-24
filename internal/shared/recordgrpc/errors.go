package recordgrpc

import (
	"errors"

	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	switch {
	case errors.Is(err, recorderrors.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, "missing credentials")
	case errors.Is(err, recorderrors.ErrForbidden):
		return status.Error(codes.PermissionDenied, "forbidden")
	case errors.Is(err, recorderrors.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, recorderrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	case errors.Is(err, recorderrors.ErrConflict):
		return status.Error(codes.AlreadyExists, "conflict")
	case errors.Is(err, recorderrors.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, "failed precondition")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
