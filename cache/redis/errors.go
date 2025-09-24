package redis

import "errors"

var (
	ErrURLParseFailed   = errors.New("url parse failed")
	ErrConnectionFailed = errors.New("url connection failed")
	ErrSetDataFailed    = errors.New("set data failed")
	ErrGetDataFailed    = errors.New("get data failed")
)
