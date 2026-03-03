package pb

import "errors"

var (
	ErrEmptyData      = errors.New("empty data")
	ErrInvalidMsgCode = errors.New("invalid message code")
	ErrEmptyMsgName   = errors.New("message name is empty")
)