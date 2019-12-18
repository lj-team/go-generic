package redis

import (
	"time"

	r "github.com/gomodule/redigo/redis"
	"github.com/lj-team/go-generic/db/redis/sentinel"
)

var sntnl *sentinel.Sentinel

type Pool = r.Pool

func RedisPool(cfg *Config) *Pool {

	if cfg.Failover {

		sntnl = &sentinel.Sentinel{
			Addrs:      cfg.Sentinel,
			MasterName: cfg.Cluster,
			Dial: func(addr string) (r.Conn, error) {
				timeout := 500 * time.Millisecond
				c, err := r.DialTimeout("tcp", addr, timeout, timeout, timeout)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}

	}

	return &r.Pool{
		MaxIdle:   100,
		MaxActive: 1200,
		Wait:      true,
		Dial: func() (r.Conn, error) {

			addr := cfg.Addr

			if cfg.Failover {
				addr, _ = sntnl.MasterAddr()
			}

			c, err := r.Dial("tcp", addr)

			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}
