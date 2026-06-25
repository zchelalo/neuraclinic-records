package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	filemanagementv1 "github.com/zchelalo/neuraclinic-records/gen/go/file_management/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/processfilestatus"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	"github.com/zchelalo/neuraclinic-records/internal/shared/appctx"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewHandler(service *application.Service) Handler {
	return func(ctx context.Context, routingKey string, body []byte) error {
		event, err := parseFileStatusChangedEvent(body)
		if err != nil {
			return err
		}

		fileID, err := uuid.Parse(event.GetFileId())
		if err != nil {
			return fmt.Errorf("%w: invalid file_id", ports.ErrInvalidEvent)
		}

		var occurredAt *time.Time
		if event.GetOccurredAt() != nil {
			value := event.GetOccurredAt().AsTime()
			occurredAt = &value
		}

		_, err = service.ProcessFileStatusChanged(ctx, processfilestatus.Command{
			FileID:     fileID,
			Status:     event.GetStatus(),
			OccurredAt: occurredAt,
		})
		if err == nil {
			return nil
		}
		if err == recorderrors.ErrNotFound {
			logger := appctx.Logger(ctx)
			logger.Warn("attachment not found for file status event",
				zap.String("routing_key", routingKey),
				zap.String("file_id", event.GetFileId()),
				zap.String("event_id", event.GetEventId()),
				zap.String("status", event.GetStatus().String()),
			)
			return nil
		}
		return err
	}
}

func parseFileStatusChangedEvent(body []byte) (*filemanagementv1.FileStatusChangedEvent, error) {
	event := &filemanagementv1.FileStatusChangedEvent{}
	if err := protojson.Unmarshal(body, event); err != nil {
		return nil, fmt.Errorf("%w: decode file status changed event: %v", ports.ErrInvalidEvent, err)
	}

	if _, err := uuid.Parse(event.GetEventId()); err != nil {
		return nil, fmt.Errorf("%w: invalid event_id", ports.ErrInvalidEvent)
	}
	if _, err := uuid.Parse(event.GetFileId()); err != nil {
		return nil, fmt.Errorf("%w: invalid file_id", ports.ErrInvalidEvent)
	}
	if event.GetServiceOrigin() != "record" {
		return nil, fmt.Errorf("%w: invalid service_origin", ports.ErrInvalidEvent)
	}
	if event.GetStatus() != sharedv1.FileStatus_FILE_STATUS_AVAILABLE && event.GetStatus() != sharedv1.FileStatus_FILE_STATUS_ERROR {
		return nil, fmt.Errorf("%w: invalid status", ports.ErrInvalidEvent)
	}
	if event.GetOccurredAt() == nil || event.GetOccurredAt().AsTime().IsZero() {
		return nil, fmt.Errorf("%w: missing occurred_at", ports.ErrInvalidEvent)
	}

	return event, nil
}
