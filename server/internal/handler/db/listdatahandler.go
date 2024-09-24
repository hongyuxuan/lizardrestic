package db

import (
	"net/http"

	"github.com/hongyuxuan/lizardrestic/common/errorx"
	commontypes "github.com/hongyuxuan/lizardrestic/common/types"
	"github.com/hongyuxuan/lizardrestic/server/internal/logic/db"
	"github.com/hongyuxuan/lizardrestic/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListdataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req commontypes.GetDataReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := db.NewListdataLogic(r.Context(), svcCtx)
		resp, err := l.Listdata(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
