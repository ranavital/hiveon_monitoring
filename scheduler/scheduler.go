package scheduler

import (
	"hiveon_monitoring/services"

	"github.com/jasonlvhit/gocron"
)

const MonitoringGaps = 1 // 1 Minutes

type Scheduler struct {
	s *gocron.Scheduler
}

func Init() *Scheduler {
	scheduler := &Scheduler{}
	scheduler.s = gocron.NewScheduler()
	scheduler.s.Every(uint64(MonitoringGaps)).Minute().From(gocron.NextTick()).Do(services.HandleWorkers)
	return scheduler
}

func (sc *Scheduler) Run() {
	<-sc.s.Start()
}
