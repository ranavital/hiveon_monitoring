package scheduler

import (
	"encoding/json"
	"fmt"
	"hiveon_monitoring/config"
	"hiveon_monitoring/entities"
	"hiveon_monitoring/psql"
	"net/http"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
)

func GetWorkers() ([]string, []string, error) {
	resp, err := http.Get(config.AppConfig.HiveonWorkersPath)
	if err != nil {
		return nil, nil, err
	}

	var responseJson = map[string]map[string]map[string]interface{}{}

	if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
		return nil, nil, err
	}

	onlineWorkers := []string{}
	offlineWorkers := []string{}
	for k, v := range responseJson["workers"] {
		if v["online"] != true {
			offlineWorkers = append(offlineWorkers, k)
			continue
		}
		onlineWorkers = append(onlineWorkers, k)
	}

	return onlineWorkers, offlineWorkers, nil
}

func handleOfflineWorker(name string) error {
	worker, err := psql.GetWorker(name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	curTime := time.Now()
	if err == gorm.ErrRecordNotFound {
		worker = &entities.OfflineWorker{
			CreatedAt: &curTime,
			UpdatedAt: &curTime,
			Name:      name,
		}

		return nil
	}

	if worker.CreatedAt.Before(curTime.Add(-20 * time.Minute)) {
		if worker.LastAlertTime == nil || worker.LastAlertTime.Before(curTime.Add(-1*time.Hour)) {
			if err := handleAlert(worker, &curTime); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleOnlineWorker(name string) error {
	worker, err := psql.GetWorker(name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	curTime := time.Now()
	delta := curTime.Sub(*worker.CreatedAt)
	worker.DowntimeLength = &delta
	if err := psql.UpdateWorker(worker); err != nil {
		return err
	}

	if err := psql.DeleteWorker(worker); err != nil {
		return err
	}

	return nil
}

func handleAlert(worker *entities.OfflineWorker, curTime *time.Time) error {
	// TODO: ALERT TELEGRAM
	worker.LastAlertTime = curTime
	worker.UpdatedAt = curTime
	if err := psql.UpdateWorker(worker); err != nil {
		return err
	}

	return nil
}

func handleOnlineWorkers(workersNames []string) error {
	for _, workerName := range workersNames {
		if err := handleOnlineWorker(workerName); err != nil {
			return err
		}
	}

	return nil
}

func handleOfflineWorkers(workersNames []string) error {
	for _, workerName := range workersNames {
		if err := handleOfflineWorker(workerName); err != nil {
			return err
		}
	}

	return nil
}

func handleWorkers() {
	onlineWorkers, offlineWorkers, err := GetWorkers()
	if err != nil {
		fmt.Printf("[handleWorkers]: failed to get workers: %s\n", err)
		return
	}

	if err := handleOnlineWorkers(onlineWorkers); err != nil {
		fmt.Printf("[handleWorkers]: failed to handle online workers: %s\n", err)
		return
	}

	if err := handleOfflineWorkers(offlineWorkers); err != nil {
		fmt.Printf("[handleWorkers]: failed to handle offline workers: %s\n", err)
		return
	}
}

func RunScheduler() {
	gocron.Every(15).Minute().From(gocron.NextTick()).Do(handleWorkers)
}
