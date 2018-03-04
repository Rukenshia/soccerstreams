package soccerstreams

import (
	"context"
	"os"
	"path/filepath"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

// DatastoreClient represents a database client for Google Cloud Platform Datastore
type DatastoreClient struct {
	ctx    context.Context
	client *datastore.Client
}

// NewDatastoreClient creates a new instance of DatastoreClient
func NewDatastoreClient(ctx context.Context) (*DatastoreClient, error) {
	client, err := datastore.NewClient(ctx, "soccerstreams-web", option.WithServiceAccountFile(filepath.Join(os.Getenv("HOME"), ".gcloud/service-accounts/soc-agent.json")))

	if err != nil {
		return nil, err
	}

	return &DatastoreClient{
		ctx,
		client,
	}, nil
}

// Key returns a valid matchthread key for datastore
func (d *DatastoreClient) Key(id string) *datastore.Key {
	return datastore.NameKey("matchthread", id, nil)
}

// Upsert inserts or update a Matchthread
func (d *DatastoreClient) Upsert(m *Matchthread) error {
	_, err := d.client.Mutate(d.ctx, datastore.NewUpsert(d.Key(m.DBKey()), m))

	return err
}

// Delete deletes a matchthread based on its id
func (d *DatastoreClient) Delete(id string) error {
	return d.client.Delete(d.ctx, d.Key(id))
}

// Get returns a Matchthread with the provided id. If it does not exist, an error is returned
func (d *DatastoreClient) Get(id string) (*Matchthread, error) {
	m := NewMatchthread(d)

	if err := d.client.Get(d.ctx, d.Key(id), m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetAll returns an array of Matchthreads for the provided query
func (d *DatastoreClient) GetAll(query *datastore.Query) ([]*Matchthread, error) {
	var threads []*Matchthread

	if _, err := d.client.GetAll(d.ctx, query, &threads); err != nil {
		return threads, err
	}

	for _, thread := range threads {
		thread.SetClient(d)
	}

	return threads, nil
}

// DeleteMulti deletes multiple Matchthreads based on their ids
func (d *DatastoreClient) DeleteMulti(ids []string) error {
	var keys []*datastore.Key

	for _, id := range ids {
		keys = append(keys, d.Key(id))
	}

	return d.client.DeleteMulti(d.ctx, keys)
}
