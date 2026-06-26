package recordgrpc

import (
	"context"
	"testing"

	"github.com/zchelalo/neuraclinic-records/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-records/internal/shared/i18n"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapErrorUsesLocalizedMessages(t *testing.T) {
	t.Parallel()

	ctx := appctx.WithLanguage(context.Background(), i18n.Spanish)
	err := MapError(ctx, recorderrors.ErrUnauthenticated)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status, got %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("status code = %s, want %s", st.Code(), codes.Unauthenticated)
	}
	if st.Message() != "faltan credenciales" {
		t.Fatalf("status message = %q", st.Message())
	}
}
