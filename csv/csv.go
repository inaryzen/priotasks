package csv

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

type DumpScheduler struct {
	ticker *time.Ticker
	quit   chan struct{}
}

func NewDumpScheduler(intervalSec int) DumpScheduler {
	if common.IsDebug() {
		log.Printf("starting scheduler with interval: %v\n", intervalSec)
	}
	s := DumpScheduler{
		ticker: time.NewTicker(time.Duration(intervalSec) * time.Second),
		quit:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-s.ticker.C:
				Dump()
			case <-s.quit:
				return
			}
		}
	}()
	return s
}

func (s DumpScheduler) Release() {
	common.Debug("release dump scheduler...")
	close(s.quit)
}

func Load(fileName string) error {
	common.Debug("loading dump: %v", fileName)

	appDir, err := common.ResolveAppDir()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	fileName = filepath.Join(appDir, fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	var tasks []models.Task
	for _, rec := range records {
		task, err := recToTask(rec)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		tasks = append(tasks, task)
	}

	common.Debug("tasks to be uploaded: %v", len(tasks))

	err = db.DeleteAllTasks()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		err = db.SaveTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func wrapSpecialSymbols(t models.Task) models.Task {
	t.Content = strings.ReplaceAll(t.Content, "\n", "\\n")
	t.Title = strings.ReplaceAll(t.Title, "\n", "\\n")
	return t
}

func unwrapSpecialSymbols(t models.Task) models.Task {
	t.Content = strings.ReplaceAll(t.Content, "\\n", "\n")
	return t
}

const (
	Field_Id = iota
	Field_Title
	Field_Content
	Field_Created
	Field_Updated
	Field_Completed
	Field_Priority
)

func taskToRec(c models.Task) []string {
	c = wrapSpecialSymbols(c)

	return []string{
		c.Id,
		c.Title,
		c.Content,
		c.Created.Format(consts.DEFAULT_TIME_FORMAT),
		c.Updated.Format(consts.DEFAULT_TIME_FORMAT),
		c.Completed.Format(consts.DEFAULT_TIME_FORMAT),
		strconv.Itoa(int(c.Priority)),
	}
}

func recToTask(record []string) (models.Task, error) {
	created, err := time.Parse(consts.DEFAULT_TIME_FORMAT, record[Field_Created])
	if err != nil {
		return models.EMPTY_TASK, err
	}
	updated, err := time.Parse(consts.DEFAULT_TIME_FORMAT, record[Field_Updated])
	if err != nil {
		return models.EMPTY_TASK, err
	}
	completed, err := time.Parse(consts.DEFAULT_TIME_FORMAT, record[Field_Completed])
	common.Debug("%v", completed)
	if err != nil {
		return models.EMPTY_TASK, err
	}
	priority, err := strconv.Atoi(record[Field_Priority])
	if err != nil {
		return models.EMPTY_TASK, err
	}

	result := models.Task{
		Id:        record[Field_Id],
		Title:     record[Field_Title],
		Content:   record[Field_Content],
		Created:   created,
		Updated:   updated,
		Completed: completed,
		Priority:  models.TaskPriority(priority),
	}
	result = unwrapSpecialSymbols(result)

	return result, nil
}

func Dump() error {

	appDir, err := common.ResolveAppDir()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	currentTime := time.Now()
	fileName := currentTime.Format("2006_01_02_15_04_05") + ".csv"
	fileName = filepath.Join(appDir, fileName)

	if common.IsDebug() {
		log.Printf("dump DB to: %v\n", fileName)
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	cards, err := db.Tasks()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	for _, c := range cards {
		record := taskToRec(c)
		err = w.Write(record)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}

	return nil
}
