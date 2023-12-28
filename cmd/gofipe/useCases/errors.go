package useCases

import "fmt"

type InvalidReferenceMonthError struct {
	reference ReferenceMonth
}

func (e *InvalidReferenceMonthError) Error() string {
	return fmt.Sprintf("Invalid reference month: %d/%d. The year or month is invalid.",
		int(e.reference.Year),
		int(e.reference.Month),
	)
}

type ReferenceMonthInTheFutureError struct {
	reference ReferenceMonth
}

func (e *ReferenceMonthInTheFutureError) Error() string {
	return fmt.Sprintf("The reference month is in the future: %d/%d",
		int(e.reference.Year),
		int(e.reference.Month),
	)
}

type InvalidLimitError struct {
	limit uint
}

func (e *InvalidLimitError) Error() string {
	return fmt.Sprintf("Invalid limit: %d. The limit must be between 1 and %d",
		int(e.limit),
		MaxLimit,
	)
}
