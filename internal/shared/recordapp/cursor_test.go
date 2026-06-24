package recordapp

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

type cursorTestRecord struct {
	ID uuid.UUID
}

func TestResolveCursorPagination(t *testing.T) {
	id := uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb01")
	cursor := EncodeCursor(id)

	pagination, err := ResolveCursorPagination(CursorPagination{
		AfterCursor: &cursor,
		Limit:       100,
	}, Config{PaginationLimitDefault: 10, PaginationLimitMax: 50})
	if err != nil {
		t.Fatalf("ResolveCursorPagination() error = %v", err)
	}
	if pagination.AfterID == nil || *pagination.AfterID != id {
		t.Fatalf("AfterID = %v, want %v", pagination.AfterID, id)
	}
	if pagination.Limit != 50 {
		t.Fatalf("Limit = %d, want 50", pagination.Limit)
	}
}

func TestResolveCursorPaginationRejectsBothDirections(t *testing.T) {
	after := EncodeCursor(uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb01"))
	before := EncodeCursor(uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb02"))

	_, err := ResolveCursorPagination(CursorPagination{
		AfterCursor:  &after,
		BeforeCursor: &before,
	}, Config{})
	if !errors.Is(err, ErrInvalidCursorPagination) {
		t.Fatalf("error = %v, want ErrInvalidCursorPagination", err)
	}
}

func TestBuildCursorPageForward(t *testing.T) {
	afterID := uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb00")
	records := []cursorTestRecord{
		{ID: uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb03")},
		{ID: uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb02")},
		{ID: uuid.MustParse("018ff8f2-8f0a-7b14-a98b-8f21fbe2fb01")},
	}

	page := BuildCursorPage(records, ResolvedCursorPagination{
		AfterID: &afterID,
		Limit:   2,
	}, func(record cursorTestRecord) uuid.UUID {
		return record.ID
	})

	if len(page.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(page.Items))
	}
	if page.Meta.NextCursor == nil || *page.Meta.NextCursor != EncodeCursor(records[1].ID) {
		t.Fatalf("NextCursor = %v, want cursor for %s", page.Meta.NextCursor, records[1].ID)
	}
	if page.Meta.PrevCursor == nil || *page.Meta.PrevCursor != EncodeCursor(records[0].ID) {
		t.Fatalf("PrevCursor = %v, want cursor for %s", page.Meta.PrevCursor, records[0].ID)
	}
}

func TestNormalizeBackwardCursorRows(t *testing.T) {
	records := []int{1, 2, 3, 4}
	normalized := NormalizeBackwardCursorRows(records, 4)

	want := []int{3, 2, 1, 4}
	for i := range want {
		if normalized[i] != want[i] {
			t.Fatalf("normalized[%d] = %d, want %d", i, normalized[i], want[i])
		}
	}
}
