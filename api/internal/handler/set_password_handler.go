package handler

import (
	"net/http"

	"github.com/xh-polaris/account-svc/api/internal/logic"
	"github.com/xh-polaris/account-svc/api/internal/svc"
	"github.com/xh-polaris/account-svc/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSetPasswordLogic(r.Context(), svcCtx)
		err := l.SetPassword(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
