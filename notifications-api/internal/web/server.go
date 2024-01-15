package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi"

	apiv1 "github.com/mkermilska/notifications-api/api/v1"
	"github.com/mkermilska/notifications-api/pkg/service"
)

const (
	errMsgMissingBody = "missing body"
	errMsgDecoding    = "error decoding request body"
	errMsgProcessing  = "error processing notification"
)

type APIServer struct {
	port            int
	notificationSvc *service.NotificationService
	logger          *zap.Logger
	httpServer      *http.Server
}

func New(port int, notificationSvc *service.NotificationService, logger *zap.Logger) *APIServer {
	return &APIServer{
		port:            port,
		notificationSvc: notificationSvc,
		logger:          logger,
	}
}

func (a *APIServer) Start() {
	a.logger.Info("Starting API Server", zap.Int("port", a.port))
	a.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.port),
		Handler:           a.handler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	err := a.httpServer.ListenAndServe()
	if err != nil {
		a.logger.Error("Error starting API Sever")
	}
}

func (a *APIServer) handler() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Post("/notifications", a.publishNotification)
	})

	return r
}

func (a *APIServer) publishNotification(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Error(errMsgMissingBody, zap.Error(err))
		http.Error(w, errMsgMissingBody, http.StatusBadRequest)
		return
	}
	notification := &apiv1.Notification{}
	if err = json.Unmarshal(bodyBytes, &notification); err != nil {
		a.logger.Error(errMsgDecoding, zap.ByteString("request_body", bodyBytes), zap.Error(err))
		http.Error(w, errMsgDecoding, http.StatusBadRequest)
		return
	}

	err = a.notificationSvc.ValidateNotification(notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = a.notificationSvc.PublishNotification(notification)
	if err != nil {
		a.logger.Error(errMsgProcessing, zap.Error(err))
		http.Error(w, errMsgProcessing, http.StatusInternalServerError)
		return
	}
}
