package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FileManagementClient interface {
	RequestUpload(ctx context.Context, originalName, mimeType string, sizeBytes int64, isPublic bool, serviceOrigin string) (uuid.UUID, string, time.Time, error)
	GenerateDownloadURL(ctx context.Context, id uuid.UUID) (string, time.Time, error)
	Close() error
}
