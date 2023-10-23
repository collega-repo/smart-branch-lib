package errs

import (
	"errors"
	"github.com/goccy/go-json"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrRecordFound      = errors.New("record already exists")
	ErrBindingRole      = errors.New("binding to role")
	ErrAuthFailed       = errors.New("authentication failed")
	ErrTrxNotFound      = errors.New("transaction is not found")
	ErrDuplicate        = errors.New("duplicate key value violates unique constraint")
	ErrApiNotFound      = errors.New("external api is not found")
	ErrpApiRequest      = errors.New("failed request external api")
	ErrAccountNotActive = errors.New("account is not active")
	ErrAccountClosed    = errors.New("account is closed")
	ErrBranchNotFound   = errors.New("branch is not found")
)

type ErrMap map[string]any

func (e ErrMap) Error() string {
	jsonByte, _ := json.Marshal(e)
	return string(jsonByte)
}
