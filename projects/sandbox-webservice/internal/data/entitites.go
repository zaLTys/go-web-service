package data

import (
	"time"
)

type Entity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Labels    []string  `json:"labels,omitempty"`
	Version   int32     `json:"-"`
	CreatedAt time.Time `json:"-"` //wont be displayed after marshaling
}
