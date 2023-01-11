package node

import "errors"

var (
	ErrNilPackageResp    = errors.New("nil resp pack")
	ErrNotValidBlock     = errors.New("not valid resp block")
	ErrBlockAlreadyExist = errors.New("block already exist")
)
