package filemanagement

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestForwardMetadataIncludesLanguage(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		headerUserID, "user-1",
		headerAcceptLanguage, "es",
		headerRequestID, "req-1",
		headerTraceID, "trace-1",
	))

	outgoing := forwardMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(outgoing)
	if !ok {
		t.Fatal("expected outgoing metadata")
	}
	if got := md.Get(headerAcceptLanguage); len(got) != 1 || got[0] != "es" {
		t.Fatalf("accept-language = %#v", got)
	}
}
