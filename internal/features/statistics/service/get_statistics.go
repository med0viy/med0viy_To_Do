package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/med0viy/practika/internal/core/domain"
	core_errors "github.com/med0viy/practika/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	listID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, listID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}

	tasksCreated := len(tasks)

	tasksComplited := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Complited == true {
			tasksComplited++
		}

		complitionDuration := task.ComplitionDuration()

		if complitionDuration != nil {
			totalCompletionDuration += *complitionDuration
		}
	}

	tasksComplitedRate := float64(tasksComplited) / float64(tasksCreated) * 100

	var tasksAvarageCompletionTime *time.Duration
	if tasksComplited > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksComplited)

		tasksAvarageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksComplited,
		&tasksComplitedRate,
		tasksAvarageCompletionTime,
	)
}
