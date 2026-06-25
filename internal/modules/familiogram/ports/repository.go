package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	"google.golang.org/protobuf/types/known/structpb"
)

type Repository interface {
	FamiliogramByPatientID(ctx context.Context, psychologistID, patientID uuid.UUID) (domain.Familiogram, error)
	UpdateFamiliogram(ctx context.Context, psychologistID, id uuid.UUID, data *structpb.Struct, now time.Time) (domain.Familiogram, error)
}
