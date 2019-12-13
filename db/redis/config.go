package redis

type Config struct {
	Addr     string   `json:"addr"`
	Failover bool     `json:"failover"`
	Cluster  string   `json:"cluster"`
	Sentinel []string `json:"sentinel"`
}
