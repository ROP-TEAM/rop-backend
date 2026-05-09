package repository

// internal error for repo to service
import "errors"

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)
