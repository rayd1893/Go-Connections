package middlewares

import (
	"net/http"

	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/gorilla/mux"
)

func Recovery() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger := logging.FromContext(ctx).Named("middleware.Recover")

			defer func() {
				if p := recover(); p != nil {
					logger.Errorw("http handler panic", "panic", p)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(w, r)
			return
		})
	}
}
