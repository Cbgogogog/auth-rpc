package errorx

import "google.golang.org/grpc/status"

var (
	ErrNoSuchUser      = status.Error(10001, "no such user")
	ErrWrongWechatCode = status.Error(10002, "wrong wechat code")
	ErrInvalidArgument = status.Error(10003, "invalid argument")
	ErrWrongPassword   = status.Error(10004, "wrong password")
)
