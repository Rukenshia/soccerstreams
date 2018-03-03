package soccerstream

type DBClient interface {
	Get(string) (*Matchthread, error)
	Delete(string) error
	Upsert(*Matchthread) error
}
