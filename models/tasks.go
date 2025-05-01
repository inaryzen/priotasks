package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inaryzen/priotasks/common"
)

var EMPTY_TASK = Task{
	Priority: PriorityMedium,
	Impact:   ImpactModerate,
	Cost:     CostM,
}

var NOT_COMPLETED time.Time = time.Time{}

const TITLE_MAX_SIZE = 64

type TaskTag string

type TaskPriority int

const (
	PriorityLow TaskPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

func StrToTaskPriority(a string) (TaskPriority, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return PriorityMedium, err
	} else {
		return TaskPriority(val), nil
	}
}

func (p TaskPriority) ToStr() string {
	switch p {
	case PriorityUrgent:
		return "🚨 Urgent"
	case PriorityHigh:
		return "🔥 High"
	case PriorityMedium:
		return "🌀 Medium"
	case PriorityLow:
		return "⌛ Low"
	default:
		return "Unknown"
	}
}

type TaskImpact int

const (
	ImpactSlight TaskImpact = iota
	ImpactLow
	ImpactModerate
	ImpactConsiderable
	ImpactHigh
)

func StrToImpact(a string) (TaskImpact, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return ImpactModerate, err
	} else {
		return TaskImpact(val), nil
	}
}

func (i TaskImpact) ToHumanString() string {
	switch i {
	case ImpactHigh:
		return "XL - 🌟🌟🌟🌟"
	case ImpactConsiderable:
		return "L - 🌟🌟🌟"
	case ImpactModerate:
		return "M - 🌟🌟"
	case ImpactLow:
		return "S - 🌱"
	case ImpactSlight:
		return "XS"
	default:
		return "Unknown"
	}
}

func StrToEnum[T ~int](a string) (T, error) {
	val, err := strconv.Atoi(a)
	if err != nil {
		return T(0), err
	} else {
		return T(val), nil
	}
}

type TaskCost int

const (
	CostXS TaskCost = iota
	CostS
	CostM
	CostL
	CostXL
	CostXXL
)

func (i TaskCost) ToHumanString() string {
	switch i {
	case CostXS:
		return "~10m"
	case CostS:
		return "~30m"
	case CostM:
		return "~1h"
	case CostL:
		return "~2h"
	case CostXL:
		return "~4h"
	case CostXXL:
		return "~8h"
	default:
		return "Unknown"
	}
}

type Task struct {
	Id        string
	Title     string
	Content   string
	Created   time.Time
	Updated   time.Time
	Completed time.Time
	Priority  TaskPriority
	Wip       bool
	Planned   bool
	Impact    TaskImpact
	Cost      TaskCost
	Value     float32
	Tags      []TaskTag
}

func titleFromContent(content string) string {
	titleIdx := len(content)
	if len(content) > TITLE_MAX_SIZE {
		titleIdx = TITLE_MAX_SIZE
	}
	result := content[:titleIdx]

	// exclude line-breaks
	titleIdx = strings.Index(result, "\n")
	if titleIdx != -1 {
		result = content[:titleIdx]
	}
	titleIdx = strings.Index(result, "\r\n")
	if titleIdx != -1 {
		result = content[:titleIdx]
	}
	titleIdx = strings.Index(result, "\r")
	if titleIdx != -1 {
		result = content[:titleIdx]
	}
	return result
}

func (t Task) AsNewTask() Task {
	if t.Title == "" {
		t.Title = titleFromContent(t.Content)
	}
	t.Id = uuid.NewString()
	t.Created = time.Now()
	return t
}

func (c Task) Update(change Task) Task {
	if change.Title == "" {
		change.Title = titleFromContent(change.Content)
	}

	common.Debug("Update: c.IsCompleted()=%v; change.IsCompleted=%v", c.IsCompleted(), change.IsCompleted())

	var completed = c.Completed
	if !c.IsCompleted() || !change.IsCompleted() {
		completed = change.Completed
	}

	return Task{
		Id:        c.Id,
		Title:     change.Title,
		Content:   change.Content,
		Created:   c.Created,
		Updated:   time.Now(),
		Completed: completed,
		Priority:  change.Priority,
		Wip:       change.Wip,
		Planned:   change.Planned,
		Impact:    change.Impact,
		Cost:      change.Cost,
		Value:     change.Value,
		Tags:      change.Tags,
	}
}

func (c Task) IsCompleted() bool {
	return c.Completed != NOT_COMPLETED
}

func (c Task) Complete() Task {
	c.Completed = time.Now()
	c.Updated = time.Now()
	return c
}

func (c Task) Uncomplete() Task {
	c.Completed = NOT_COMPLETED
	c.Updated = time.Now()
	return c
}

func (c Task) CalculateValue() Task {
	var priorityMultipliers = map[TaskPriority]float32{
		PriorityUrgent: 1.7,
		PriorityHigh:   1.2,
		PriorityMedium: 1.0,
		PriorityLow:    0.8,
	}

	baseValue := (float32(c.Impact)*1.1 + 1) * (float32(CostXXL-c.Cost) + 1)
	c.Value = baseValue * priorityMultipliers[c.Priority]
	return c
}

func (c Task) ValueAsHumanStr() string {
	buckets := []struct {
		low  int
		high int
		str  string
	}{
		{low: 30, high: 60, str: "💵💵💵💵"},
		{low: 22, high: 30, str: "💵💵💵"},
		{low: 13, high: 22, str: "💵💵"},
		{low: 7, high: 13, str: "💵"},
		{low: 0, high: 7, str: ""},
	}

	for _, b := range buckets {
		if c.Value > float32(b.low) && c.Value < float32(b.high) {
			return b.str
		}
	}

	return "???"
}

func (t Task) IsEmpty() bool {
	return t.Title == EMPTY_TASK.Title &&
		t.Content == EMPTY_TASK.Content &&
		t.Priority == EMPTY_TASK.Priority &&
		t.Impact == EMPTY_TASK.Impact &&
		t.Cost == EMPTY_TASK.Cost &&
		t.Value == EMPTY_TASK.Value &&
		!t.Wip &&
		!t.Planned &&
		len(t.Tags) == 0 &&
		t.Completed == NOT_COMPLETED
}
