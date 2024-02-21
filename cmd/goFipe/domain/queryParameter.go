package domain

const MaxLimit = 100

type WhereClause struct {
	Column   string
	Operator string
	Value    interface{}
}

type Pagination struct {
	Offset int
	Limit  int
}

type OrderByClause struct {
	Column string
	IsDesc bool
}
