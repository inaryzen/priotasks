package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var EMPTY_TASK = Task{
	Priority: PriorityMedium,
	Impact:   ImpactModerate,
}

var NOT_COMPLETED time.Time = time.Time{}

const TITLE_MAX_SIZE = 64

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

func Create(prototype Task) Task {
	if prototype.Title == "" {
		prototype.Title = titleFromContent(prototype.Content)
	}

	return Task{
		Id:        uuid.NewString(),
		Title:     prototype.Title,
		Content:   prototype.Content,
		Created:   time.Now(),
		Updated:   time.Now(),
		Completed: prototype.Completed,
		Priority:  prototype.Priority,
		Wip:       prototype.Wip,
		Planned:   prototype.Planned,
		Impact:    prototype.Impact,
	}
}

func (c Task) Update(change Task) Task {
	if change.Title == "" {
		change.Title = titleFromContent(change.Content)
	}

	var completed time.Time
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
