package models

import (
	"fmt"
	"slices"
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
	ColumnImpact
	ColumnWip
	ColumnPlanned
	ColumnCost
	ColumnValue
	ColumnTags
	ColumnFun
)

func (sc SortColumn) ToHumanString() string {
	return []string{"Undefined", "Completed", "Title", "Created", "Updated", "Priority", "Impact", "WIP", "Plan", "T", "Value", "Tags", "Fun"}[sc]
}

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

func (s Settings) IsSorted(c SortColumn, d SortDirection) bool {
	return s.TasksQuery.SortColumn == c && s.TasksQuery.SortDirection == d
}

type TasksQuery struct {
	FilterCompleted   bool
	FilterIncompleted bool
	CompletedFrom     time.Time
	CompletedTo       time.Time
	SortColumn        SortColumn
	SortDirection     SortDirection
	FilterWip         bool
	FilterNonWip      bool
	Planned           bool
	NonPlanned        bool
	Tags              []TaskTag
}

func (t TasksQuery) RemoveTag(target TaskTag) TasksQuery {
	var newTags []TaskTag = nil
	for i, tag := range t.Tags {
		if tag == target {
			newTags = slices.Delete(t.Tags, i, i+1)
			t.Tags = newTags
			break
		}
	}
	return t
}

func (t TasksQuery) String() string {
	return fmt.Sprintf(
		"FilterCompleted: %v, "+
			"CompletedFrom: %v, "+
			"CompletedTo: %v, "+
			"SortColumn: %v, "+
			"SortDirection: %v, "+
			"FilterCompleted: %v, "+
			"FilterWip: %v, "+
			"FilterNonWip: %v, "+
			"Planned: %v, "+
			"NonPlanned: %v, "+
			"Tags: %v",
		t.FilterCompleted,
		t.CompletedFrom,
		t.CompletedTo,
		t.SortColumn,
		t.SortDirection,
		t.FilterIncompleted,
		t.FilterWip,
		t.FilterNonWip,
		t.Planned,
		t.NonPlanned,
		t.Tags,
	)
}

func (s TasksQuery) Reset() TasksQuery {
	s.FilterCompleted = true
	s.FilterIncompleted = false
	s.CompletedFrom = time.Now().AddDate(0, 0, -14)
	s.CompletedTo = NOT_COMPLETED
	s.SortColumn = Priority
	s.SortDirection = Desc
	s.FilterWip = false
	s.FilterNonWip = false
	s.Planned = false
	s.NonPlanned = false
	s.Tags = []TaskTag{}
	return s
}
