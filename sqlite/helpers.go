package sqlite

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func CheckConstraint(err error, constraintName string) bool {
	sqliteErr, ok := errors.AsType[*sqlite.Error](err)
	if !ok {
		return false
	}

	if sqliteErr.Code() != sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return false
	}

	return strings.Contains(sqliteErr.Error(), constraintName)
}

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Time struct{ T time.Time }

func (t Time) Value() (driver.Value, error) {
	return t.T.UTC().Format(rfc3339Milli), nil
}

func (t *Time) Scan(v any) error {
	if v == nil {
		return nil
	}

	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("unsupported scan: storing driver type %T into string", v)
	}

	parsedT, err := time.Parse(rfc3339Milli, s)
	if err != nil {
		return fmt.Errorf("unsupported scan: parsing %s into time.Time: %v", s, err)
	}

	t.T = parsedT.UTC()
	return nil
}
