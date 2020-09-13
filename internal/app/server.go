package apiserver

import (
	"database/sql"

	balanceD "github.com/Andronovdima/job-avito-assignment/internal/app/balance/delivery"
	balanceR "github.com/Andronovdima/job-avito-assignment/internal/app/balance/repository"
	balanceU "github.com/Andronovdima/job-avito-assignment/internal/app/balance/usecase"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Mux    *mux.Router
	Config *Config
	Logger *zap.SugaredLogger
}

func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	s := &Server{
		Mux:    mux.NewRouter(),
		Logger: logger,
		Config: config,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *Server) ConfigureServer(db *sql.DB) {
	balanceRep := balanceR.NewBalanceRepository(db)
	balanceUC := balanceU.NewBalanceUsecase(balanceRep)
	balanceD.NewBalanceHandler(s.Mux, *balanceUC, s.Logger)
}
