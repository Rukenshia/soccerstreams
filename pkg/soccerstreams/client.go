package soccerstreams

import "cloud.google.com/go/datastore"

type DBClient interface {
	Get(string) (*Matchthread, error)

	// TODO: Change to more generic Query
	GetAll(*datastore.Query) ([]*Matchthread, error)

	Delete(string) error
	DeleteMulti([]string) error
	Upsert(*Matchthread) error
}
