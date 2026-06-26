package filemanagement

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	filemanagementv1 "github.com/zchelalo/neuraclinic-records/gen/go/file_management/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	headerUserID         = "x-user-id"
	headerPsychologistID = "x-psychologist-id"
	headerAdminID        = "x-admin-id"
	headerRequestID      = "x-request-id"
	headerTraceID        = "x-trace-id"
	headerAcceptLanguage = "accept-language"
)

type Config struct {
	Addr               string
	TLSEnabled         bool
	CACertPath         string
	InsecureSkipVerify bool
}

type Client struct {
	conn   *grpc.ClientConn
	client filemanagementv1.FileManagementServiceClient
}

func New(cfg Config) (*Client, error) {
	creds, err := transportCredentials(cfg)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("create file-management grpc client: %w", err)
	}

	return &Client{
		conn:   conn,
		client: filemanagementv1.NewFileManagementServiceClient(conn),
	}, nil
}

func (c *Client) RequestUpload(ctx context.Context, originalName, mimeType string, sizeBytes int64, isPublic bool, serviceOrigin string) (uuid.UUID, string, time.Time, error) {
	resp, err := c.client.RequestUpload(forwardMetadata(ctx), &filemanagementv1.FileManagementServiceRequestUploadRequest{
		OriginalName:  originalName,
		MimeType:      mimeType,
		SizeBytes:     sizeBytes,
		IsPublic:      isPublic,
		ServiceOrigin: serviceOrigin,
	})
	if err != nil {
		return uuid.Nil, "", time.Time{}, mapGRPCError(err)
	}
	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, "", time.Time{}, fmt.Errorf("parse file id: %w", err)
	}
	expiresAt := time.Time{}
	if resp.GetExpiresAt() != nil {
		expiresAt = resp.GetExpiresAt().AsTime()
	}
	return id, resp.GetUploadUrl(), expiresAt, nil
}

func (c *Client) GenerateDownloadURL(ctx context.Context, id uuid.UUID) (string, time.Time, error) {
	resp, err := c.client.GenerateDownloadUrl(forwardMetadata(ctx), &filemanagementv1.FileManagementServiceGenerateDownloadUrlRequest{
		Id: id.String(),
	})
	if err != nil {
		return "", time.Time{}, mapGRPCError(err)
	}
	expiresAt := time.Time{}
	if resp.GetExpiresAt() != nil {
		expiresAt = resp.GetExpiresAt().AsTime()
	}
	return resp.GetDownloadUrl(), expiresAt, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func forwardMetadata(ctx context.Context) context.Context {
	incoming, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	out := metadata.MD{}
	for _, key := range []string{headerUserID, headerPsychologistID, headerAdminID, headerRequestID, headerTraceID, headerAcceptLanguage} {
		if values := incoming.Get(key); len(values) > 0 {
			out.Set(key, values...)
		}
	}
	if len(out) == 0 {
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, out)
}

func transportCredentials(cfg Config) (credentials.TransportCredentials, error) {
	if !cfg.TLSEnabled {
		return insecure.NewCredentials(), nil
	}

	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
	if cfg.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true
		return credentials.NewTLS(tlsCfg), nil
	}
	if cfg.CACertPath != "" {
		ca, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("read file-management ca cert: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("append file-management ca cert")
		}
		tlsCfg.RootCAs = pool
	}

	return credentials.NewTLS(tlsCfg), nil
}

func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return recorderrors.ErrUnauthenticated
	case codes.PermissionDenied:
		return recorderrors.ErrForbidden
	case codes.NotFound:
		return recorderrors.ErrNotFound
	case codes.InvalidArgument:
		return recorderrors.ErrInvalidInput
	case codes.AlreadyExists:
		return recorderrors.ErrConflict
	case codes.FailedPrecondition:
		return recorderrors.ErrFailedPrecondition
	default:
		return err
	}
}

var _ ports.FileManagementClient = (*Client)(nil)
