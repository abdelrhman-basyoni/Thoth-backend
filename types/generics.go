package typ

type PaginatedEntities[T any] struct {
	Total    int64
	Entities []T
}
