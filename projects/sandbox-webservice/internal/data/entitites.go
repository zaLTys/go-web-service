package data

import (
	"time"
)

type Entity struct {
	ID        int64
	Name      string
	Labels    []string
	Version   int32
	CreatedAt time.Time
}
