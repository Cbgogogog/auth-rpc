package logic

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/xh-polaris/auth-rpc/internal/errorx"
	"github.com/xh-polaris/auth-rpc/internal/model"
	"github.com/xh-polaris/auth-rpc/internal/svc"
	"github.com/xh-polaris/auth-rpc/pb"
	"math/big"
	"net/smtp"
	"strconv"

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
				"您正在进行喵社区账号注册，本次注册验证码为：%s，5分钟内有效，请勿透露给其他人。\r\n", in.AuthId, code.String())))
		if err != nil {
			return nil, err
		}
		verifyCode = code.String()
	default:
		return nil, errorx.ErrInvalidArgument
	}
	err := l.svcCtx.Redis.SetexCtx(l.ctx, VerifyCodeKeyPrefix+in.AuthId, verifyCode, 300)
	if err != nil {
		return nil, err
	}
	return &pb.SendVerifyCodeResp{}, nil
}
