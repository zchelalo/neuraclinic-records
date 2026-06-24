package recordgrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/shared/appctx"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

func PsychologistID(ctx context.Context) (uuid.UUID, error) {
	id, ok := appctx.PsychologistID(ctx)
	if !ok || id == uuid.Nil {
		return uuid.Nil, recorderrors.ErrUnauthenticated
	}
	return id, nil
}

func UserID(ctx context.Context) *uuid.UUID {
	id, ok := appctx.UserID(ctx)
	if !ok || id == uuid.Nil {
		return nil
	}
	return &id
}

func ParseID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, recorderrors.ErrInvalidInput
	}
	return id, nil
}

func ParseOptionalID(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	id, err := ParseID(*value)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
