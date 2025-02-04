package models

import (
	"strconv"
	"time"
)

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
	Id         string
	TasksQuery TasksQuery
}

type TasksQuery struct {
	FilterCompleted bool
	CompletedFrom   time.Time
	CompletedTo     time.Time
	SortColumn      SortColumn
	SortDirection   SortDirection
}

func (s Settings) IsSorted(c SortColumn, d SortDirection) bool {
	return s.TasksQuery.SortColumn == c && s.TasksQuery.SortDirection == d
}
