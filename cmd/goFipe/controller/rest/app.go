package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest/handler"
	"github.com/raffops/gofipe/cmd/goFipe/controller/rest/middleware"
	"github.com/raffops/gofipe/cmd/goFipe/logger"
	"github.com/raffops/gofipe/cmd/goFipe/port"
	"net/http"
	"os"
)

func Start(vehicleService port.VehicleService) {
	sanityCheck()
	router := mux.NewRouter()
	vehicleHandler := handler.NewHandler(vehicleService)
	router.HandleFunc("/health-check", healthCheck).Methods("GET")
	router.HandleFunc("/vehicles", vehicleHandler.Get).Methods("GET")

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	logger.Info("Starting server %s:%s", logger.String("host", appHost), logger.String("port", appPort))

	loggedRouter := middleware.LoggingMiddleware()(router)

	err := http.ListenAndServe(
		fmt.Sprintf("%s:%s", appHost, appPort),
		loggedRouter,
	)
	if err != nil {
		logger.Fatal("Error starting server", logger.String("error", err.Error()))
	}
}

func sanityCheck() {
	if _, ok := os.LookupEnv("APP_HOST"); !ok {
		logger.Fatal("APP_HOST environment variable is not set")
	}

	if _, ok := os.LookupEnv("APP_PORT"); !ok {
		logger.Fatal("APP_PORT environment variable is not set")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
