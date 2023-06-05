package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/silenceper/wechat/v2/util"

	"github.com/xh-polaris/auth-rpc/internal/errorx"
	"github.com/xh-polaris/auth-rpc/internal/model"
	"github.com/xh-polaris/auth-rpc/internal/svc"
	"github.com/xh-polaris/auth-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/crypto/bcrypt"
)

type SignInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	OAuthUrl            = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	VerifyCodeKeyPrefix = "verify:"
)

func NewSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInLogic {
	return &SignInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignInLogic) SignIn(in *pb.SignInReq) (resp *pb.SignInResp, err error) {
	resp = &pb.SignInResp{}
	resp.User = &pb.User{}
	switch in.AuthType {
	case model.EmailAuthType:
		fallthrough
	case model.PhoneAuthType:
		resp.User.UserId, err = l.signInByPassword(in)
	case model.WechatAuthType:
		resp.User.UserId, resp.User.UnionId, resp.User.OpenId, resp.User.AppId, err = l.signInByWechat(in)
	default:
		return nil, errorx.ErrInvalidArgument
	}
	if err != nil {
		return nil, err
	}
	return
}

func (l *SignInLogic) signInByPassword(in *pb.SignInReq) (string, error) {
	userModel := l.svcCtx.UserModel

	// 检查是否设置了验证码，若设置了检查验证码是否合法
	ok, err := l.checkVerifyCode(in.Params, in.AuthId)
	if err != nil {
		return "", err
	}

	auth := model.Auth{
		Type:  in.AuthType,
		Value: in.AuthId,
	}
	user, err := userModel.FindOneByAuth(l.ctx, auth)

	switch err {
	case nil:
	case model.ErrNotFound:
		if !ok {
			return "", errorx.ErrNoSuchUser
		}

		user = &model.User{Auth: []model.Auth{auth}}
		err := userModel.Insert(l.ctx, user)
		if err != nil {
			return "", err
		}
		return user.ID.Hex(), nil
	default:
		return "", err
	}

	if ok {
		return user.ID.Hex(), nil
	}

	// 验证码未通过，尝试密码登录
	if user.Password == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		return "", errorx.ErrWrongPassword
	}

	return user.ID.Hex(), nil
}

func (l *SignInLogic) checkVerifyCode(opts []string, authValue string) (bool, error) {
	verifyCode, err := l.svcCtx.Redis.GetCtx(l.ctx, VerifyCodeKeyPrefix+authValue)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	} else if len(opts) < 1 || verifyCode != opts[0] {
		return false, nil
	} else {
		return true, nil
	}
}

// return userId unionId openId appid
func (l *SignInLogic) signInByWechat(in *pb.SignInReq) (string, string, string, string, error) {
	opts := in.Params
	if len(opts) < 1 {
		return "", "", "", "", errorx.ErrInvalidArgument
	}
	jsCode := opts[0]

	var unionId string
	var openId string
	var appId string
	if len(opts) < 2 {
		// 向微信开放接口提交临时code
		res, err := l.svcCtx.Meowchat.GetAuth().Code2SessionContext(l.ctx, jsCode)
		if err != nil {
			return "", "", "", "", err
		} else if res.ErrCode != 0 {
			return "", "", "", "", errors.New(res.ErrMsg)
		}
		unionId = res.UnionID
		openId = res.OpenID
		appId = l.svcCtx.Meowchat.GetContext().AppID
	} else if opts[1] == "old" {
		res, err := l.svcCtx.MeowchatOld.GetAuth().Code2SessionContext(l.ctx, jsCode)
		if err != nil {
			return "", "", "", "", err
		} else if res.ErrCode != 0 {
			return "", "", "", "", errors.New(res.ErrMsg)
		}
		unionId = res.UnionID
		openId = res.OpenID
		appId = l.svcCtx.MeowchatOld.GetContext().AppID
	} else if opts[1] == "manager" {
		c := l.svcCtx.Config.MeowchatManager
		res, err := util.HTTPGetContext(l.ctx, fmt.Sprintf(OAuthUrl, c.AppID, c.AppSecret, jsCode))
		if err != nil {
			return "", "", "", "", err
		}
		var j map[string]any
		if err = json.Unmarshal(res, &j); err != nil {
			return "", "", "", "", err
		}
		if _, ok := j["unionid"]; !ok {
			return "", "", "", "", errorx.ErrWrongWechatCode
		}
		unionId = j["unionid"].(string)
		if _, ok := j["openid"]; !ok {
			return "", "", "", "", errorx.ErrWrongWechatCode
		}
		openId = j["openid"].(string)
		appId = c.AppID
	}

	userModel := l.svcCtx.UserModel
	auth := model.Auth{
		Type:  in.AuthType,
		Value: unionId,
	}
	user, err := userModel.FindOneByAuth(l.ctx, auth)
	switch err {
	case nil:
	case model.ErrNotFound:
		user = &model.User{Auth: []model.Auth{auth}}
		err := userModel.Insert(l.ctx, user)
		if err != nil {
			return "", "", "", "", err
		}
		return user.ID.Hex(), unionId, openId, appId, nil
	default:
		return "", "", "", "", err
	}

	return user.ID.Hex(), unionId, openId, appId, nil
}
