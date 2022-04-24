package progorm

// Page holds information about result page
type Page struct {
	Total   uint
	Page    uint
	PerPage uint
	Pages   uint
}

// PageTyped holds information about result page
type PageTyped[T any] struct {
	Total   uint
	Page    uint
	PerPage uint
	Pages   uint
	Results []T
}
