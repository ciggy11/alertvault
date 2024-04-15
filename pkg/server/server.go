package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/ciggy11/alertvault/pkg/alertstore"
	"github.com/ciggy11/alertvault/pkg/alert"
	"github.com/ciggy11/alertvault/pkg/metrics"
	"github.com/ciggy11/alertvault/pkg/config"
)

// SupportedWebhookVersion is the alert webhook data version that is supported
// by this app
const SupportedWebhookVersion = "4"

// Server represents a web server that processes webhooks
type Server struct {
	config config.Config
	r      *mux.Router
	store  alertstore.AlertStore
}

// New returns a new web server
func New(c config.Config) (Server, error) {
	r := mux.NewRouter()
	// intialize the store
	store, err := alertstore.NewAlertStore(context.Background(), c)
	if err != nil {
		return Server{}, err
	}
	s := Server{
		config: c,
		r:      r,
		store: store,
	}
	r.HandleFunc("/webhook", s.Store).Methods("POST")
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/alert/{id}/history", s.HistoryAlert).Methods("GET")
	return s, nil
}

// Start starts a new server on the given address
func (s Server) Start(address string) {
	log.Println("Starting listener on", address, "using", s.config)
	log.Fatal(http.ListenAndServe(address, s.r))
}

func (s Server) Store(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	metrics.WebhooksReceivedTotal.Inc()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		metrics.InvalidWebhooksTotal.Inc()
		log.Errorf("Failed to read payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to read payload: %s", err), http.StatusBadRequest)
		return
	}
	log.Debug("Received webhook payload: ", string(body))
	data, err := alert.Parse(body)
	if err != nil {
		metrics.InvalidWebhooksTotal.Inc()
		log.Errorf("Failed to parse payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to parse payload: %s", err), http.StatusBadRequest)
		return
	}
	for _, a := range data.Alerts{
		var tenantID string
		if s.config.Tenant.InLabel {
			tenantID = a.Labels[s.config.Tenant.Label]
		} else if s.config.Tenant.InAnnotation {
			tenantID = a.Annotations[s.config.Tenant.Annotation]
		}
		uniqueKey := a.Labels[s.config.Tenant.UniqueName]
		alertKey := tenantID + "|" + uniqueKey
		err := s.store.SetTenantAlert(r.Context(), alertKey, &a)
		if err != nil {
			log.Errorf("Failed to store alert: %s\n", err)
		}
	}
	groupTenantID := data.CommonLabels["tenantID"]
	groupErr := s.store.SetAlertGroup(r.Context(), groupTenantID, data)
	if groupErr != nil {
		log.Errorf("Failed to store alert group: %s\n", groupErr)
		http.Error(w, fmt.Sprintf("Failed to store alert group: %s", err), http.StatusInternalServerError)
		return
	}
	metrics.AlertsReceivedTotal.WithLabelValues(data.Receiver, data.Status).Add(float64(len(data.Alerts)))
}

func (s Server) HistoryAlert(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	metrics.AlertsGetTotal.Inc()
	params := mux.Vars(r)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	score,_ := strconv.ParseFloat(r.URL.Query().Get("score"), 64)

	tenantID := r.Header.Get(s.config.Tenant.Header)
	uniqueID := params["id"]
	key := tenantID + "|" + uniqueID
	alertDesc := alert.NewAlertDesc(key, score, int64(offset), int64(limit))
	alerts, err := s.store.GetTenantAlerts(r.Context(), alertDesc)
	if err != nil {
		metrics.AlertsGetFailuresTotal.Inc()
		log.Errorf("Failed to get alerts: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to get alerts: %s", err), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(alerts)
	if err != nil {
		metrics.InvalidAlertGetTotal.Inc()
        log.Errorf("Failed to encode alerts: %s\n", err)
        http.Error(w, fmt.Sprintf("Failed to encode alerts: %s", err), http.StatusInternalServerError)
    }
}


