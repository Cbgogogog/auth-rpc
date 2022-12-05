// Code generated by goctl. DO NOT EDIT.
// Source: auth.proto

package server

import (
	"context"

	"github.com/xh-polaris/auth-rpc/internal/logic"
	"github.com/xh-polaris/auth-rpc/internal/svc"
	"github.com/xh-polaris/auth-rpc/pb"
)

type AuthServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedAuthServer
}

func NewAuthServer(svcCtx *svc.ServiceContext) *AuthServer {
	return &AuthServer{
		svcCtx: svcCtx,
	}
}

func (s *AuthServer) SignIn(ctx context.Context, in *pb.SignInReq) (*pb.SignInResp, error) {
	l := logic.NewSignInLogic(ctx, s.svcCtx)
	return l.SignIn(in)
}

func (s *AuthServer) SetPassword(ctx context.Context, in *pb.SetPasswordReq) (*pb.SetPasswordResp, error) {
	l := logic.NewSetPasswordLogic(ctx, s.svcCtx)
	return l.SetPassword(in)
}

func (s *AuthServer) SendVerifyCode(ctx context.Context, in *pb.SendVerifyCodeReq) (*pb.SendVerifyCodeResp, error) {
	l := logic.NewSendVerifyCodeLogic(ctx, s.svcCtx)
	return l.SendVerifyCode(in)
}
