package service

import (
	"fmt"
	"fundinvest/internal/task007/effects"
	"fundinvest/internal/task007/model"
	"fundinvest/internal/task007/worker"
)

func Retry(job model.Job, client *effects.Client, w worker.Worker) (model.Job, error) {
	for attempt := 0; attempt < 2; attempt++ {
		key := fmt.Sprintf("settlement:%s:%d", job.ID, attempt)
		if err := client.Execute(key); err != nil {
			if attempt == 1 {
				return job, err
			}
			continue
		}
		completed := job.Complete()
		w.Publish(completed)
		return completed, nil
	}
	return job, nil
}
