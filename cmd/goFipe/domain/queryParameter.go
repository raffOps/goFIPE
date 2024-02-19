package domain

type Condition struct {
	Column   string
	Operator string
	Value    interface{}
}

type Pagination struct {
	Offset int
	Limit  int
}

type OrderBy struct {
	Column string
	IsDesc bool
}
