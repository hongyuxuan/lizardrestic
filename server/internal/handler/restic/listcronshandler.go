package restic

import (
	"net/http"

	"github.com/hongyuxuan/lizardrestic/server/internal/logic/restic"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListcronsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := restic.NewListcronsLogic(r.Context(), svcCtx)
		resp, err := l.Listcrons()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}