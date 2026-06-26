package recordgrpc

import (
	"context"
	"errors"

	"github.com/zchelalo/neuraclinic-records/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-records/internal/shared/i18n"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(ctx context.Context, err error) error {
	language := appctx.Language(ctx)
	switch {
	case errors.Is(err, recorderrors.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, i18n.Message(language, i18n.KeyMissingCredentials))
	case errors.Is(err, recorderrors.ErrForbidden):
		return status.Error(codes.PermissionDenied, i18n.Message(language, i18n.KeyForbidden))
	case errors.Is(err, recorderrors.ErrNotFound):
		return status.Error(codes.NotFound, i18n.Message(language, i18n.KeyNotFound))
	case errors.Is(err, recorderrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, i18n.Message(language, i18n.KeyInvalidInput))
	case errors.Is(err, recorderrors.ErrConflict):
		return status.Error(codes.AlreadyExists, i18n.Message(language, i18n.KeyConflict))
	case errors.Is(err, recorderrors.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, i18n.Message(language, i18n.KeyFailedPrecondition))
	default:
		return status.Error(codes.Internal, i18n.Message(language, i18n.KeyInternalServerError))
	}
}
