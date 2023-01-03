package network

import "errors"

var (
	ErrIncorrectAddress = errors.New("IncorrectAddress")
	ErrNotOpt           = errors.New("NotOpts")
	ErrNotPack          = errors.New("NotPack")
	ErrTimeWait         = errors.New("TimeWait")
)
