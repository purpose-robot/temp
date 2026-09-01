package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	driverName    = "sqlite"
	readerOptions = "?mode=ro&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	writerOptions = "?mode=rw&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_txlock=immediate"
)

type db struct {
	Reader *sql.DB
	Writer *sql.DB
}

func (h *db) ping(ctx context.Context) error {
	if err := h.Reader.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database reader: %v", err)
	}

	if err := h.Writer.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database writer: %v", err)
	}

	return nil
}

func (h *db) close() error {
	if err := h.Reader.Close(); err != nil {
		return fmt.Errorf("close database reader: %v", err)
	}

	if err := h.Writer.Close(); err != nil {
		return fmt.Errorf("close database writer: %v", err)
	}

	return nil
}

func Open(name string) (*db, error) {
	reader, err := sql.Open(driverName, name+readerOptions)
	if err != nil {
		return nil, fmt.Errorf("open database reader: %v", err)
	}

	writer, err := sql.Open(driverName, name+readerOptions)
	if err != nil {
		closeErr := reader.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return nil, fmt.Errorf("open database writer: %v", err)
	}

	writer.SetMaxOpenConns(1)
	writer.SetMaxIdleConns(1)

	writer.SetConnMaxLifetime(0)
	writer.SetConnMaxIdleTime(0)

	handler := &db{Reader: reader, Writer: writer}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = handler.ping(ctx)
	if err != nil {
		closeErr := handler.close()
		if err != nil {
			err = errors.Join(err, closeErr)
		}

		return nil, err
	}

	return handler, nil
}
