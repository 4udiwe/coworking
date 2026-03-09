package app

import scheduler_service "github.com/4udiwe/cowoking/scheduler-service/internal/service/scheduler"

func (app *App) SchedulerService() *scheduler_service.SchedulerService {
	if app.schedulerService != nil {
		return app.schedulerService
	}
	app.schedulerService = scheduler_service.New(
		app.TimerRepo(),
		app.Postgres(),
		app.cfg.Scheduler.RemindBefore,
	)
	return app.schedulerService
}
