package alert

import (
	"encoding/json"
	"fmt"
	"time"
)

// AlertGroup is the data read from a webhook call
type AlertGroup struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`

	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alerts is a slice of Alert
type Alerts []Alert

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// AlertsResp is the response for the GetTenantAlerts method
type AlertsResp struct {
	Alerts []*Alert `json:"alerts"`
	Total  int64    `json:"total"`
	Offset int64    `json:"offset"`
	Limit  int64    `json:"limit"`
}

// AlertsDesc is the struct to get alerts from redis
type AlertsDesc struct {
	Key    string
	Score  float64
	Offset int64
	Count  int64
}

// Convert string to Alert
func StringToAlerts(strings []string) ([]*Alert, error) {
	var alerts []*Alert
	for _, str := range strings {
		var alert Alert
		err := json.Unmarshal([]byte(str), &alert)
		if err != nil {
			return nil, fmt.Errorf("error parsing string to TenantAlert: %s", err)
		}
		alerts = append(alerts, &alert)
	}
	return alerts, nil
}

// Parse gets a webhook payload and parses it returning
// AlertGroup object if successful
func Parse(payload []byte) (*AlertGroup, error) {
	d := AlertGroup{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %s", err)
	}
	return &d, nil
}

// NewAlertDesc creates a new AlertDesc object
func NewAlertDesc(key string, score float64, offset, count int64) *AlertsDesc {
	return &AlertsDesc{
		Key:    key,
		Score:  score,
		Offset: offset,
		Count:  count,
	}
}
