package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/domain"
	"google.golang.org/protobuf/types/known/structpb"
)

type Repository interface {
	FamilyogramByPatientID(ctx context.Context, psychologistID, patientID uuid.UUID) (domain.Familyogram, error)
	UpdateFamilyogram(ctx context.Context, psychologistID, id uuid.UUID, data *structpb.Struct, now time.Time) (domain.Familyogram, error)
}
