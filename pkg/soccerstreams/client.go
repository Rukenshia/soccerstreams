package soccerstreams

import "cloud.google.com/go/datastore"

// DBClient represents an interface for persisting Matchthreads
type DBClient interface {
	// TODO: Move the actual implementation (i.e. Datastore) into another package
	Get(string) (*Matchthread, error)

	// TODO: Change to more generic Query
	GetAll(*datastore.Query) ([]*Matchthread, error)

	Delete(string) error
	DeleteMulti([]string) error
	Upsert(*Matchthread) error
}
