package models

// Competition represents a Football competition (e.g. Premier League)
type Competition struct {
	Name       string
	Identifier string

	Matchthreads FrontendMatchthreads
}
