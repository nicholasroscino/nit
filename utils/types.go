package utils

import "errors"

type StagedObject struct {
	Hash      string
	Path      string
	Timestamp string
}
type HashAlreadyExist struct{}

func (m *HashAlreadyExist) Error() string {
	return "Hash already exists in the repository"
}

func IsHashAlreadyExist(err error) bool {
	var hashAlreadyExist *HashAlreadyExist
	ok := errors.As(err, &hashAlreadyExist)
	return ok
}
