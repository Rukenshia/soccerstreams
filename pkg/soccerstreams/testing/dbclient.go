package testing

import (
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

var (
	// ErrNotExists represents an error when trying to retrieve a Matchthread from the database that does not exist
	ErrNotExists = errors.New("Matchthread does not exist")
)

// MockDBClient represents a mock database client (in-memory) that can be used for testing.
type MockDBClient struct {
	Threads map[string]*soccerstreams.Matchthread
}

// NewMockDBClient creates a new MockDBClient
func NewMockDBClient() *MockDBClient {
	return &MockDBClient{
		make(map[string]*soccerstreams.Matchthread),
	}
}

// Add adds a Matchthread to the mock database
func (m *MockDBClient) Add(mt *soccerstreams.Matchthread) {
	m.Threads[mt.DBKey()] = mt
}

// Get returns a Matchthread with the given key
func (m *MockDBClient) Get(key string) (*soccerstreams.Matchthread, error) {
	mt, ok := m.Threads[key]
	if !ok {
		return nil, ErrNotExists
	}

	return mt, nil
}

// GetAll returns all Matchthreads. The provided query is not respected
func (m *MockDBClient) GetAll(*datastore.Query) ([]*soccerstreams.Matchthread, error) {
	var mts []*soccerstreams.Matchthread

	for _, mt := range m.Threads {
		mts = append(mts, mt)
	}

	return mts, nil
}

// Delete deletes a matchthread with the given key
func (m *MockDBClient) Delete(key string) error {
	delete(m.Threads, key)
	return nil
}

// DeleteMulti deletes all given matchthreads
func (m *MockDBClient) DeleteMulti(keys []string) error {
	for _, k := range keys {
		m.Delete(k)
	}
	return nil
}

// Upsert inserts or updates a Matchthread
func (m *MockDBClient) Upsert(mt *soccerstreams.Matchthread) error {
	m.Add(mt)
	return nil
}
