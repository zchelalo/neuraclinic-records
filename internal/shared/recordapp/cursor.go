package recordapp

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidCursorPagination = errors.New("invalid cursor pagination")

type CursorPagination struct {
	AfterCursor  *string
	BeforeCursor *string
	Limit        int32
}

type ResolvedCursorPagination struct {
	AfterID  *uuid.UUID
	BeforeID *uuid.UUID
	Limit    int32
}

type CursorMeta struct {
	NextCursor *string
	PrevCursor *string
	Limit      *int32
}

type CursorPage[T any] struct {
	Items []T
	Meta  CursorMeta
}

func ResolveCursorPagination(p CursorPagination, cfg Config) (ResolvedCursorPagination, error) {
	afterCursor := cleanCursor(p.AfterCursor)
	beforeCursor := cleanCursor(p.BeforeCursor)
	if afterCursor != nil && beforeCursor != nil {
		return ResolvedCursorPagination{}, ErrInvalidCursorPagination
	}

	resolved := ResolvedCursorPagination{Limit: normalizeLimit(p.Limit, cfg)}
	if afterCursor != nil {
		id, err := DecodeCursor(*afterCursor)
		if err != nil {
			return ResolvedCursorPagination{}, ErrInvalidCursorPagination
		}
		resolved.AfterID = &id
	}
	if beforeCursor != nil {
		id, err := DecodeCursor(*beforeCursor)
		if err != nil {
			return ResolvedCursorPagination{}, ErrInvalidCursorPagination
		}
		resolved.BeforeID = &id
	}
	return resolved, nil
}

func (p ResolvedCursorPagination) IsBackward() bool {
	return p.BeforeID != nil
}

func (p ResolvedCursorPagination) QueryLimit() int32 {
	return p.Limit + 1
}

func BuildCursorPage[T any](records []T, pagination ResolvedCursorPagination, idOf func(T) uuid.UUID) CursorPage[T] {
	limit := int(pagination.Limit)
	hasAdditionalPage := len(records) > limit

	items := records
	if hasAdditionalPage {
		items = records[:limit]
	}

	limitValue := pagination.Limit
	meta := CursorMeta{Limit: &limitValue}
	if len(items) == 0 {
		return CursorPage[T]{Items: items, Meta: meta}
	}

	firstID := idOf(items[0])
	lastID := idOf(items[len(items)-1])

	if pagination.IsBackward() || hasAdditionalPage {
		value := EncodeCursor(lastID)
		meta.NextCursor = &value
	}
	if pagination.IsBackward() {
		if hasAdditionalPage {
			value := EncodeCursor(firstID)
			meta.PrevCursor = &value
		}
	} else if pagination.AfterID != nil {
		value := EncodeCursor(firstID)
		meta.PrevCursor = &value
	}

	return CursorPage[T]{Items: items, Meta: meta}
}

func NormalizeBackwardCursorRows[T any](records []T, queryLimit int32) []T {
	if len(records) == 0 {
		return records
	}

	normalized := make([]T, 0, len(records))
	if len(records) == int(queryLimit) {
		for i := len(records) - 2; i >= 0; i-- {
			normalized = append(normalized, records[i])
		}
		normalized = append(normalized, records[len(records)-1])
		return normalized
	}

	for i := len(records) - 1; i >= 0; i-- {
		normalized = append(normalized, records[i])
	}
	return normalized
}

func EncodeCursor(id uuid.UUID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(id.String()))
}

func DecodeCursor(value string) (uuid.UUID, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		raw, err = base64.URLEncoding.DecodeString(value)
		if err != nil {
			return uuid.Nil, err
		}
	}
	return uuid.Parse(string(raw))
}

func cleanCursor(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeLimit(limit int32, cfg Config) int32 {
	if limit <= 0 {
		limit = cfg.PaginationLimitDefault
	}
	if limit <= 0 {
		limit = 10
	}
	if cfg.PaginationLimitMax > 0 && limit > cfg.PaginationLimitMax {
		limit = cfg.PaginationLimitMax
	}
	return limit
}
