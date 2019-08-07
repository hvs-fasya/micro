package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/hvs-fasya/micro/internal/configure"
	"github.com/hvs-fasya/micro/internal/server/handlers"
)

//Srv package level variable
var Srv *Server

//Server http server struct
type Server struct {
	Cfg *configure.Server
}

//Run start http server
func Run(cfg *configure.Server) {
	Srv = &Server{cfg}
	connstr := Srv.Cfg.Host + ":" + Srv.Cfg.Port
	zap.L().Info("http server start at " + connstr)
	e := http.ListenAndServe(connstr, NewRouter())
	if e != nil {
		zap.L().Error(errors.Wrap(e, "http server start error").Error())
	}
}

//NewRouter create router with routes
func NewRouter() *mux.Router {
	rt := new(mux.Router)
	rt.Use(setID)
	rt.Use(recovery)
	apiRouter := rt.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(setAPIHeaders)
	apiRouter.Use(logRequest)
	apiRouter.HandleFunc("/alive", handlers.Alive).Methods("GET")
	return rt
}

func setAPIHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func setID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewV4()
		ctx := context.WithValue(r.Context(), "ID", requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		requestID := r.Context().Value("ID").(uuid.UUID)
		Srv.Cfg.Logger.Info(requestID.String(), zap.String("request", r.Method+" "+r.RequestURI), zap.String("tooktime", time.Since(start).String()))
	})
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r *http.Request) {
			rec := recover()
			if rec != nil {
				requestID := r.Context().Value("ID").(uuid.UUID)
				zap.L().Error("PANIC: " + fmt.Sprintf("%v", rec))
				Srv.Cfg.Logger.Error(fmt.Sprintf("%v", rec), zap.String("request", r.Method+" "+r.RequestURI), zap.String("ID", requestID.String()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}(r)
		next.ServeHTTP(w, r)
	})
}
