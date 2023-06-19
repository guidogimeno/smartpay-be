package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guidogimeno/smartpay-be/services"
	"github.com/guidogimeno/smartpay-be/types"
)

type apiServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *apiServer {
	return &apiServer{
		listenAddr: listenAddr,
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
	paymentRequest := new(types.PaymentRequest)
	if err := json.NewDecoder(r.Body).Decode(paymentRequest); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	payment := paymentRequest.ToPayment()
	analysis := services.PaymentAnalysis(payment)
	analysisResponse := analysis.ToAnalysisResponse()

	writeJSON(w, http.StatusOK, analysisResponse)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
