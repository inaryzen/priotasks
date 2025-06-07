package consts

const (
	FILTER_NAME_HIDE_COMPLETED   = "completed-filter"
	FILTER_NAME_HIDE_INCOMPLETED = "filter-incompleted"
	FILTER_COMPLETED_FROM        = "filter-completed-from"
	FILTER_COMPLETED_TO          = "filter-completed-to"
	FILTER_WIP                   = "filter-wip"
	FILTER_NON_WIP               = "filter-non-wip"
	FILTER_PLANNED               = "filter-planned"
	FILTER_NON_PLANNED           = "filter-non-planned"
	FILTER_TAGS                  = "filter-tags"
	FILTER_SEARCH                = "filter-search"

	PREPARED_QUERY_RESET                    = "prepared-query-clear"
	PREPARED_QUERY_COMPLETED_YESTERDAY      = "prepared-query-completed-yesterday"
	PREPARED_QUERY_COMPLETED_TODAY          = "prepared-query-completed-today"
	PREPARED_QUERY_COMPLETED_THIS_WEEK      = "prepared-query-completed-this-week"
	PREPARED_QUERY_COMPLETED_LAST_TWO_WEEKS = "prepared-query-completed-last-two-weeks"
	PREPARED_QUERY_COMPLETED_LAST_WEEK      = "prepared-query-completed-last-week"

	COMPLETED_SORT_NAME  = "completed-sort"
	SORT_COLUMN_NAME     = "sort-column"
	SORT_DIRECTION_NAME  = "sort-direction"
	MODAL_TASK_COST_NAME = "modal-task-cost"

	URL_TOGGLE_SORT_TABLE = "/toggle-sort-table"
	URL_TASKS             = "/tasks"
	URL_TASKS_ID          = "/tasks/{id}"
	URL_TASKS_EXPORT_YAML = "/tasks/export/yaml"

	DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"
	DEFAULT_DATE_FORMAT = "2006-01-02"

	INPUT_NAME_NEW_TAG = "input-name-new-tag"
)
