package restic

import (
	"net/http"

	"github.com/hongyuxuan/lizardrestic/server/internal/logic/restic"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeletepolicyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BackupReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := restic.NewDeletepolicyLogic(r.Context(), svcCtx)
		resp, err := l.Deletepolicy(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
