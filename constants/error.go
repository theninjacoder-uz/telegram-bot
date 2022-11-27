package constants

import "errors"

var (
	// ErrTableAlreadyExists ...
	ErrTableAlreadyExists = errors.New("table already exists")
	// ErrTableNotExists ...
	ErrTableNotExists = errors.New("table not exists")
)

const (
	//PGForeignKeyViolationCode ...
	PGForeignKeyViolationCode = "23503"
	//PGUniqueKeyViolationCode ...
	PGUniqueKeyViolationCode = "23505"
)
