package models

import "strconv"

type SortColumn int

const (
	ColumnUndefined SortColumn = iota
	Completed
	Title
	Created
	Updated
	Priority
)

func ColumnFromString(str string) (result SortColumn) {
	num, err := strconv.Atoi(str)
	if err != nil {
		result = ColumnUndefined
	} else {
		result = SortColumn(num)
	}
	return
}

type SortDirection int

const (
	DirectionUndefined SortDirection = iota
	Desc
	Asc
)

func (d SortDirection) Flip() SortDirection {
	if d == Desc {
		return Asc
	} else {
		return Desc
	}
}

func DirectionFromString(str string) (result SortDirection) {
	num, err := strconv.Atoi(str)
	if err != nil {
		result = DirectionUndefined
	} else {
		result = SortDirection(num)
	}
	return
}

type Settings struct {
	Id                  string
	FilterCompleted     bool
	ActiveSortColumn    SortColumn
	ActiveSortDirection SortDirection
}

func (s Settings) IsSorted(c SortColumn, d SortDirection) bool {
	return s.ActiveSortColumn == c && s.ActiveSortDirection == d
}
