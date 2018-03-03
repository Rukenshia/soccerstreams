package soccerstream

import (
	"context"
	"os"
	"path/filepath"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

type Datastore struct {
	client *datastore.Client
}

func NewDatastoreClient() (*Datastore, error) {
	c, err := datastore.NewClient(context.Background(), "soccerstreams-web", option.WithServiceAccountFile(filepath.Join(os.Getenv("HOME"), ".gcloud/service-accounts/soc-agent.json")))

	if err != nil {
		return nil, err
	}

	return &Datastore{
		client: c,
	}, nil
}

func (d *Datastore) Key(id string) *datastore.Key {
	return datastore.NameKey("matchthread", id, nil)
}

func (d *Datastore) Upsert(m *Matchthread) error {
	_, err := d.client.Mutate(context.Background(), datastore.NewUpsert(d.Key(m.DBKey()), m))

	return err
}

func (d *Datastore) Delete(id string) error {
	return d.client.Delete(context.Background(), d.Key(id))
}

func (d *Datastore) Get(id string) (*Matchthread, error) {
	m := NewMatchthread(d)

	if err := d.client.Get(context.Background(), d.Key(id), m); err != nil {
		return nil, err
	}

	return m, nil
}
