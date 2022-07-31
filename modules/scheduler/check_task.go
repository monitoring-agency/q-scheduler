package scheduler

import (
	"fmt"
	"github.com/myOmikron/echotools/utility"
	"os/exec"
	"time"

	"github.com/myOmikron/echotools/worker"

	"github.com/monitoring-agency/q-scheduler/models"
)

func runCheck(cmdline string) (*models.Result, error) {
	start := time.Now()
	cmd := exec.Command(cmdline)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return &models.Result{
		Output:        string(stdout),
		ExecutionTime: time.Now().Sub(start),
		ReturnCode:    cmd.ProcessState.ExitCode(),
	}, nil
}

func createCheckTask(check models.Check, s *scheduler) worker.Task {
	return worker.NewTask(func() error {
		now := time.Now().UTC()
		for _, day := range check.SchedulingPeriod.Days {
			if day.Day == int(now.Weekday()) {
				for _, tp := range day.TimePeriods {
					currMinute := now.Minute() + now.Hour()*60

					if tp.Start <= currMinute && tp.Stop > currMinute {
						// Check is allowed to run
						result, err := runCheck(check.Commandline)
						if err != nil {
							fmt.Println("error: " + err.Error())
							return err
						}

						utility.PPrintln(result)
					}
				}
			}
		}

		time.Sleep(time.Duration(check.SchedulingInterval) * time.Second)
		s.pool.AddTask(createCheckTask(check, s))
		return nil
	})
}
