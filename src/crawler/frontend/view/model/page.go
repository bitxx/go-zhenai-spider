package model

type SearchResult struct {
	Hits int64
	Start int
	Items []interface{}
	Query string
	PrevFrom int
	NextFrom int
}
