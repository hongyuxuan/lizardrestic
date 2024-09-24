package restic

import (
	"net/http"

	"github.com/hongyuxuan/lizardrestic/common/errorx"
	"github.com/hongyuxuan/lizardrestic/server/internal/logic/restic"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/hongyuxuan/lizardrestic/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LssnapshotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LsSnapshotsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := restic.NewLssnapshotsLogic(r.Context(), svcCtx)
		resp, err := l.Lssnapshots(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
