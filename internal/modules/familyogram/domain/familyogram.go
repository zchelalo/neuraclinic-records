package domain

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

type Familyogram struct {
	ID        uuid.UUID
	Data      *structpb.Struct
	PatientID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
