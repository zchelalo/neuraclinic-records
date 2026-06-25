package updatefamiliogram

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/ports"
	recordapp "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestExecuteUpdatesFamiliogramData(t *testing.T) {
	ctx := context.Background()
	psychologistID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	familiogramID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	now := time.Date(2026, 6, 25, 18, 0, 0, 0, time.UTC)
	data, err := structpb.NewStruct(map[string]any{
		"objeto_1": map[string]any{
			"x": 500,
			"y": 1200,
		},
	})
	if err != nil {
		t.Fatalf("NewStruct returned error: %v", err)
	}

	repo := &repoStub{}
	uc := New(repo, recordapp.Runtime{Now: func() time.Time { return now }})

	result, err := uc.Execute(ctx, Command{
		PsychologistID: psychologistID,
		ID:             familiogramID,
		Data:           data,
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if repo.psychologistID != psychologistID {
		t.Fatalf("expected psychologist id %s, got %s", psychologistID, repo.psychologistID)
	}
	if repo.id != familiogramID {
		t.Fatalf("expected familiogram id %s, got %s", familiogramID, repo.id)
	}
	if repo.data == nil || repo.data.AsMap()["objeto_1"] == nil {
		t.Fatalf("expected data to be passed to repository, got %#v", repo.data)
	}
	if !repo.now.Equal(now) {
		t.Fatalf("expected now %s, got %s", now, repo.now)
	}
	if result.ID != familiogramID {
		t.Fatalf("expected result id %s, got %s", familiogramID, result.ID)
	}
}

func TestExecuteRejectsNilData(t *testing.T) {
	uc := New(&repoStub{}, recordapp.Runtime{})

	_, err := uc.Execute(context.Background(), Command{
		PsychologistID: uuid.New(),
		ID:             uuid.New(),
		Data:           nil,
	})
	if err != recorderrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

type repoStub struct {
	psychologistID uuid.UUID
	id             uuid.UUID
	data           *structpb.Struct
	now            time.Time
}

func (r *repoStub) FamiliogramByPatientID(context.Context, uuid.UUID, uuid.UUID) (domain.Familiogram, error) {
	panic("unexpected call")
}

func (r *repoStub) UpdateFamiliogram(_ context.Context, psychologistID, id uuid.UUID, data *structpb.Struct, now time.Time) (domain.Familiogram, error) {
	r.psychologistID = psychologistID
	r.id = id
	r.data = data
	r.now = now
	return domain.Familiogram{
		ID:        id,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

var _ ports.Repository = (*repoStub)(nil)
