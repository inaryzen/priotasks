package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

var EMPTY_CARD = Card{
	Priority: PriorityMedium,
}

var NOT_COMPLETED time.Time = time.Time{}

const TITLE_MAX_SIZE = 64

type TaskPriority int

const (
	PriorityLow TaskPriority = iota
	PriorityMedium
	PriorityHigh
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

type Card struct {
	Id        string
	Title     string
	Content   string
	Created   time.Time
	Updated   time.Time
	Completed time.Time
	Priority  TaskPriority
}

func titleFromContent(content string) string {
	titleIdx := len(content)
	if len(content) > TITLE_MAX_SIZE {
		titleIdx = TITLE_MAX_SIZE
	}
	return content[:titleIdx]
}

func Create(prototype Card) Card {
	if prototype.Title == "" {
		prototype.Title = titleFromContent(prototype.Content)
	}

	return Card{
		Id:        uuid.NewString(),
		Title:     prototype.Title,
		Content:   prototype.Content,
		Created:   time.Now(),
		Updated:   time.Now(),
		Completed: NOT_COMPLETED,
		Priority:  prototype.Priority,
	}
}

func (c Card) Update(change Card) Card {
	if change.Title == "" {
		change.Title = titleFromContent(change.Content)
	}

	return Card{
		Id:        c.Id,
		Title:     change.Title,
		Content:   change.Content,
		Created:   c.Created,
		Updated:   time.Now(),
		Completed: c.Completed,
		Priority:  change.Priority,
	}
}

func (c Card) IsCompleted() bool {
	return c.Completed != NOT_COMPLETED
}

func (c Card) Complete() Card {
	c.Completed = time.Now()
	c.Updated = time.Now()
	return c
}

func (c Card) Uncomplete() Card {
	c.Completed = NOT_COMPLETED
	c.Updated = time.Now()
	return c
}
