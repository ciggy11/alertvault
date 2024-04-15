# AlertVault
AlertVault receives [HTTP Webhook](https://prometheus.io/docs/alerting/configuration/#webhook-receiver-%3Cwebhook_config%3E) notifications from [Cortex/Mimir AlertManager](https://github.com/cortexproject) and inserted into selected database for storing and analysis.
It stores both [Alert Group](https://github.com/prometheus/alertmanager) and [Alert](https://github.com/prometheus/alertmanager) per tenant.

Having this data can used for:
- Tune alerting rules
- Understand incident
- Understand alert's behavior during incident

# Limitation
AlertVault can not capture silenced or inhibited alerts.  

# Building
You will need:  
- Make
- Go 1.22 or above
- a working GOPATH

```
# This will create amd64 binary 
make build
```
# Configuration
```
# Where to listen for webhook payload
http_listen_address: 127.0.0.1:8080 
log_level: info

# Database that used for storing alerts. Currently supported S3 and Redis
backend: redis

# Where to find tenant infomations
tenant:
  in_label: true
  in_annotation: false
  label: tenantID
  annotation: tenantID
  unique_name: fingerprint
  header: X-Scope-OrgID

vaultdb:
  redis:
    addrs: 0.0.0.0:6379
    timeout: 3s
    expiration: 604800s
    alerts_db: 0
    alert_group_db: 1
  s3:
    bucket: vaultdb
    region: us-west-2
    endpoint: http://localhost:9000
    access_key: minio
    secret_key: minio123
```
# Alertmanager configuration example

```
route:
  routes:
  - receiver: alertvault
    continue: true

- name: alertvault
  webhook_configs:
    - url: https://alertvault.example.com/webhook
```
