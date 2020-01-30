package api

import "errors"

var (
	ErrMissingID    = errors.New("Missing identifier in URI")          // 400
	ErrMismatchedID = errors.New("URI doesn't match provided data")    // 400
	ErrBodyRequired = errors.New("Body is required for this endpoint") // 400

	ErrNotFound = errors.New("Resource Not Found") // 404
)
