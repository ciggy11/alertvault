package alertstore

import (
	"context"

	"github.com/ciggy11/alertvault/pkg/alert"
	redisClient "github.com/ciggy11/alertvault/pkg/alertstore/redisclient"
	"github.com/ciggy11/alertvault/pkg/config"
	"github.com/ciggy11/alertvault/pkg/db/redis"
)

type AlertStore interface {
	// Returns all history of alerts by tenant
	GetTenantAlerts(ctx context.Context, ad *alert.AlertsDesc) (*alert.AlertsResp, error)

	// Store history of alert by tenant
	SetTenantAlert(ctx context.Context, key string, alert *alert.Alert) error

	// Store history of alert group by tenant
	SetAlertGroup(ctx context.Context, tenantID string, alerts *alert.AlertGroup) error

	// Total Alert by tenant and alert name
	TotalByKey(ctx context.Context, key string) (int64, error)

	// DeleteAlerts delete history of alert by tenant
	DeleteTenantAlerts(ctx context.Context, tenantID string) error
}

func NewAlertStore(ctx context.Context, cfg config.Config) (AlertStore, error) {
	if cfg.Backend == redisClient.NAME {
		alertClient := redis.NewRedisAlertsClient(&cfg.VaultDB.Redis)
		alertGroupClient := redis.NewRedisGroupAlertsClient(&cfg.VaultDB.Redis)
		c := redisClient.NewStore(*alertClient, *alertGroupClient)
		return c, nil
	}
	return nil, nil
}
