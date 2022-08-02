package scheduler

import (
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/anmitsu/go-shlex"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/worker"

	"github.com/monitoring-agency/q-scheduler/models"
)

func runCheck(cmdline string, s *scheduler) (*models.Result, error) {
	start := time.Now()

	parts, err := shlex.Split(cmdline, true)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.configuration.ProcessTimeout)*time.Second)
	defer cancel()

	args := make([]string, 0)
	if len(parts) > 1 {
		args = parts[1:]
	}

	cmd := exec.CommandContext(ctx, parts[0], args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return &models.Result{
			Output:        err.Error(),
			ReturnCode:    2,
			ExecutionTime: time.Now().Sub(start),
		}, nil
	}

	if err := cmd.Wait(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return &models.Result{
				Output:        "Check took to long to execute",
				ExecutionTime: time.Now().Sub(start),
			}, nil
		}
		if e, ok := err.(*exec.ExitError); ok {
			return &models.Result{
				Output:        e.String(),
				ReturnCode:    e.ExitCode(),
				ExecutionTime: time.Now().Sub(start),
			}, nil
		} else {
			return nil, err
		}
	}

	output, _ := ioutil.ReadAll(stdout)
	return &models.Result{
		Output:        string(output),
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
						result, err := runCheck(check.Commandline, s)
						if err != nil {
							color.Println(color.RED, "error: "+err.Error())
							return err
						}

						// Fill rest of struct
						result.CheckID = check.ID

					}
				}
			}
		}

		go func() {
			time.Sleep(time.Duration(check.SchedulingInterval) * time.Second)
			s.pool.AddTask(createCheckTask(check, s))
		}()
		return nil
	})
}
