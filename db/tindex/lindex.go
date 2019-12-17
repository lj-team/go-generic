package tindex

import (
	"github.com/lj-team/go-generic/db/tindex/engine"
)

type Engine = engine.Engine

func Open(driver string, dsn string) (Engine, error) {

	return engine.Open(driver, dsn)
}
