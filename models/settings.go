package models

import (
	"fmt"
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
	FilterCompleted   bool
	FilterIncompleted bool
	CompletedFrom     time.Time
	CompletedTo       time.Time
	SortColumn        SortColumn
	SortDirection     SortDirection
}

func (t TasksQuery) String() string {
	return fmt.Sprintf("FilterCompleted: %v, CompletedFrom: %v, CompletedTo: %v, SortColumn: %v, SortDirection: %v, FilterCompleted: %v",
		t.FilterCompleted, t.CompletedFrom, t.CompletedTo, t.SortColumn, t.SortDirection, t.FilterIncompleted)
}

func (s Settings) IsSorted(c SortColumn, d SortDirection) bool {
	return s.TasksQuery.SortColumn == c && s.TasksQuery.SortDirection == d
}
