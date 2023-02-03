package logic

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"strconv"

	"github.com/xh-polaris/auth-rpc/internal/errorx"
	"github.com/xh-polaris/auth-rpc/internal/model"
	"github.com/xh-polaris/auth-rpc/internal/svc"
	"github.com/xh-polaris/auth-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyCodeLogic {
	return &SendVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendVerifyCodeLogic) SendVerifyCode(in *pb.SendVerifyCodeReq) (*pb.SendVerifyCodeResp, error) {
	var verifyCode string
	switch in.AuthType {
	case model.PhoneAuthType:
		return nil, errorx.ErrInvalidArgument
	case model.EmailAuthType:
		c := l.svcCtx.Config.SMTP
		auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
		code, err := rand.Int(rand.Reader, big.NewInt(900000))
		code = code.Add(code, big.NewInt(100000))
		if err != nil {
			return nil, err
		}
		err = smtp.SendMail(c.Host+":"+strconv.Itoa(c.Port), auth, c.Username, []string{in.AuthId}, []byte(fmt.Sprintf(
			"To: %s\r\n"+
				"From: xh-polaris\r\n"+
				"Content-Type: text/plain"+"; charset=UTF-8\r\n"+
				"Subject: 验证码\r\n\r\n"+
				"%s\r\n", in.AuthId, code.String())))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errorx.ErrInvalidArgument
	}
	err := l.svcCtx.Redis.Hset(VerifyCodeKey, in.AuthId, verifyCode)
	if err != nil {
		return nil, err
	}
	return &pb.SendVerifyCodeResp{}, nil
}
