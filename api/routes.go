package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guidogimeno/smartpay-be/db"
	"github.com/guidogimeno/smartpay-be/services"
	"github.com/guidogimeno/smartpay-be/types"
)

type apiServer struct {
	listenAddr string
	db         db.Storer
}

func NewAPIServer(listenAddr string, db db.Storer) *apiServer {
	return &apiServer{
		listenAddr: listenAddr,
		db:         db,
	}
}

func (s *apiServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/ping", handlePing).Methods(http.MethodGet)
	router.HandleFunc("/payment/analysis", handlePaymentAnalysis).Methods(http.MethodPost)

	return http.ListenAndServe(":"+s.listenAddr, router)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("pong")
}

func handlePaymentAnalysis(w http.ResponseWriter, r *http.Request) {
	var payment *types.Payment
	if err := json.NewDecoder(r.Body).Decode(payment); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err := payment.IsValid()
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	analysis, err := services.PaymentAnalysis(payment)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	analysisResponse := analysis.ToAnalysisResponse()

	writeJSON(w, http.StatusOK, analysisResponse)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
