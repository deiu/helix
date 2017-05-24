package helix

import (
	"errors"
	"github.com/albrow/zoom"
)

var (
	GraphStore *zoom.Collection
	pool       = &zoom.Pool{}
)

type VersionedGraph struct {
	ID        string
	Last      bool
	Subject   string
	Predicate string
	Object    string
	Graph     string
	zoom.RandomId
}

func NewVersionedGraph() *VersionedGraph {
	return &VersionedGraph{}
}

func initRedisPool(URL string, pool *zoom.Pool) (*zoom.Collection, error) {
	if len(URL) == 0 {
		return &zoom.Collection{}, errors.New("Cannot initialize Redis with empty URL")
	}
	pool = zoom.NewPool(URL)
	return pool.NewCollectionWithOptions(NewVersionedGraph(),
		zoom.DefaultCollectionOptions.WithIndex(true))
}
