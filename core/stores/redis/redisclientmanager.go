package redis

import (
	"crypto/tls"
	"fmt"
	"io"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(r *Redis) (*red.Client, error) {
	val, err := clientManager.GetResource(fmt.Sprintf("%s:%d", r.Addr, r.Db), func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           r.Db,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.AddHook(durationHook)

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
