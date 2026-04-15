package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksComplited             int
	TasksComplitedRate         *float64
	TasksAvarageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksComplited int,
	tasksComplitedRate *float64,
	tasksAvarageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TasksComplited:             tasksComplited,
		TasksComplitedRate:         tasksComplitedRate,
		TasksAvarageCompletionTime: tasksAvarageCompletionTime,
	}
}
