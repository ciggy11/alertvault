package redisclient

import (
	"context"
	"encoding/json"

	"github.com/ciggy11/alertvault/pkg/alert"
	"github.com/ciggy11/alertvault/pkg/db/redis"
	log "github.com/sirupsen/logrus"
)

const (
	NAME = "redis"
)

type Store struct {
	alertClient      redis.RedisClient
	alertGroupClient redis.RedisClient
}

func NewStore(ac redis.RedisClient, agc redis.RedisClient) *Store {
	return &Store{
		alertClient:      ac,
		alertGroupClient: agc,
	}
}

func (c *Store) GetTenantAlerts(ctx context.Context, ad *alert.AlertsDesc) (*alert.AlertsResp, error) {
	alertJson, err := c.alertClient.ZGetByScore(ctx, ad)
	if err != nil {
		log.Errorf("Error getting alerts from redis: %v", err)
		return nil, err
	}
	total, errTotal := c.alertClient.Count(ctx, ad.Key)
	if errTotal != nil {
		log.Errorf("Error getting total alerts from redis: %v", errTotal)
		return nil, errTotal
	}
	alertsRsp := &alert.AlertsResp{	
		Alerts: alertJson,
		Total:  total,
		Offset: ad.Offset,
		Limit:  total,
	}
	return alertsRsp, nil
}

func (c *Store) SetTenantAlert(ctx context.Context, key string, a *alert.Alert) error {
	alertsJson, err := json.Marshal(a)
	if err != nil {
		log.Errorf("Error marshalling alert: %v", err)
		panic(err)
	}
	c.alertClient.ZSet(ctx, alertsJson, key, float64(a.StartsAt.Unix()))
	return nil
}

func (c *Store) SetAlertGroup(ctx context.Context, tenantID string, a *alert.AlertGroup) error {
	alertGroupJson, err := json.Marshal(a)
	if err != nil {
		log.Errorf("Error marshalling alert group: %v", err)
		panic(err)
	}
	c.alertGroupClient.Set(ctx, tenantID, alertGroupJson)
	return nil
}

func (c *Store) TotalByKey(ctx context.Context, key string) (int64, error) {
	total, err := c.alertClient.Count(ctx, key)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// Currently not implemented
func (c *Store) DeleteTenantAlerts(ctx context.Context, tenantID string) error {
	return nil
}
