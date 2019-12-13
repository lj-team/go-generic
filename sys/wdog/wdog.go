package wdog

import (
	"context"
	"os"
	"time"
)

type KillFunc func()

type WatchDog struct {
	aliveBefore int64
	done        chan bool
	ttl         int64
	stop        context.CancelFunc
	killFunc    KillFunc
}

func New(ttl int64, bf KillFunc) *WatchDog {

	wd := &WatchDog{
		aliveBefore: time.Now().Unix() + ttl,
		done:        make(chan bool),
		ttl:         ttl,
		killFunc:    bf,
	}

	ctx, fn := context.WithCancel(context.Background())

	wd.stop = fn

	go func() {

		for {

			select {
			case <-time.After(time.Nanosecond * 500000000):
			case <-ctx.Done():
				close(wd.done)
				return
			}

			if wd.aliveBefore < time.Now().Unix() {
				if wd.killFunc != nil {
					wd.killFunc()
					close(wd.done)
					break
				} else {
					os.Exit(1)
				}
			}

		}

	}()

	return wd
}

func (wd *WatchDog) Alive() {
	wd.aliveBefore = time.Now().Unix() + wd.ttl
}

func (wd *WatchDog) Close() {
	wd.stop()
	<-wd.done
}
