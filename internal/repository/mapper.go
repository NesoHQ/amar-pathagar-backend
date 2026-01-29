package repository

import (
	"database/sql"
	"time"
)

// Helper functions to convert between sql.Null* and pointers

func stringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func float64Ptr(f sql.NullFloat64) *float64 {
	if f.Valid {
		return &f.Float64
	}
	return nil
}

func intPtr(i sql.NullInt64) *int {
	if i.Valid {
		val := int(i.Int64)
		return &val
	}
	return nil
}

func timePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func nullFloat64(f *float64) sql.NullFloat64 {
	if f != nil {
		return sql.NullFloat64{Float64: *f, Valid: true}
	}
	return sql.NullFloat64{}
}

func nullInt64(i *int) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{}
}

func nullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}
