package services

import (
	"encoding/json"
	"fmt"
	"hiveon_monitoring/config"
	"hiveon_monitoring/entities"
	"hiveon_monitoring/logger"
	"hiveon_monitoring/psql"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TimeBetweenAlerts = time.Duration(-1) * time.Hour     // 1 hour
	TimeAfterOffline = time.Duration(-1*20) * time.Minute // 20 Minutes
)

func getWorkers() ([]string, []string, error) {
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

	logger.Logging.Info("[getWorkers]: Online workers: %+v", onlineWorkers)
	logger.Logging.Info("[getWorkers]: Offline workers: %+v", offlineWorkers)

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
		psql.Session.Create(&worker)
		logger.Logging.Info("[handleOfflineWorker]: Successfully inserted offline worker: %s", worker.Name)
		return nil
	}

	if worker.CreatedAt.Before(curTime.Add(TimeAfterOffline)) {
		if worker.LastAlertTime == nil || worker.LastAlertTime.Before(curTime.Add(time.Duration(TimeBetweenAlerts))) {
			if err := handleAlert(worker, &curTime); err != nil {
				return err
			}
			logger.Logging.Info("[handleOfflineWorker]: Alerted offline worker: %s", worker.Name)
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
	worker.DowntimeLength = delta.String()
	if err := psql.UpdateWorker(worker); err != nil {
		return err
	}

	if err := psql.DeleteWorker(worker); err != nil {
		return err
	}

	logger.Logging.Info("[handleOnlineWorker]: Worker %s is online", worker.Name)
	if worker.LastAlertTime == nil {
		return nil
	}

	if err := SendTelegramAlert(fmt.Sprintf("Worker %s is online, thanks!", worker.Name)); err != nil {
		return err
	}

	return nil
}

func handleAlert(worker *entities.OfflineWorker, curTime *time.Time) error {
	customMsg := ""
	switch worker.Name {
	case "MiriRegev":
		customMsg = "Ran stop playing RL, you are always losing!!!"
	case "THEOERIGISBACK2", "ARGAZ":
		customMsg = "Tal stop playing Factorio, it's a shitty game!!!"
	case "BoratSagdiyev":
		customMsg = "Matan call mama... NOW!"
	case "MainOERig":
		customMsg = "ARGAZIM ALERT, CALL mama ASAP!!!!!!!!"
	}
	if err := SendTelegramAlert(fmt.Sprintf("Worker %s is offline, %s", worker.Name, customMsg)); err != nil {
		return err
	}

	worker.LastAlertTime = curTime
	worker.UpdatedAt = curTime
	if err := psql.UpdateWorker(worker); err != nil {
		return err
	}

	logger.Logging.Info("[handleAlert]: Worker %s is offline", worker.Name)
	return nil
}

func handleOnlineWorkers(workersNames []string) error {
	for _, workerName := range workersNames {
		logger.Logging.Info("[handleOnlineWorkers]: Handling online worker %s", workerName)
		if err := handleOnlineWorker(workerName); err != nil {
			return err
		}
	}

	return nil
}

func handleOfflineWorkers(workersNames []string) error {
	for _, workerName := range workersNames {
		logger.Logging.Info("[handleOfflineWorkers]: Handling offline worker %s", workerName)
		if err := handleOfflineWorker(workerName); err != nil {
			return err
		}
	}

	return nil
}

func HandleWorkers() {
	onlineWorkers, offlineWorkers, err := getWorkers()
	if err != nil {
		logger.Logging.Error("[handleWorkers]: failed to get workers: %s\n", err)
		return
	}

	if err := handleOnlineWorkers(onlineWorkers); err != nil {
		logger.Logging.Error("[handleWorkers]: failed to handle online workers: %s\n", err)
		return
	}

	if err := handleOfflineWorkers(offlineWorkers); err != nil {
		logger.Logging.Error("[handleWorkers]: failed to handle offline workers: %s\n", err)
		return
	}

	logger.Logging.Info("[handleWorkers]: successfuly handled workers")
}
