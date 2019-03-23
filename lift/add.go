package main

import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/awinterman/lifting"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func logWorkout(cmd *cobra.Command, args []string) {

	var (
		// temp vars for ui input
		sessionDateString string
		sets              string
		exercise          string
		effort            string
		volume            string
		weight            string
		duration          string
		failure           string

		// things can go wrong literally whenever
		err error

		// loaded from database
		exercises      []string
		unitsOptions   []string
		labels         []string
		recent         []lifting.Repetition
		repsByExercise = make(map[string]lifting.Repetition)

		// ui components

		enterDate = Ask{
			Label:     "Date: ",
			Default:   civil.DateOf(time.Now()).String(),
			AllowEdit: true,
		}
		selectLabel = Ask{
			Label:    "Workout Type: ",
			Items:    labels,
			AddLabel: "Add Workout Type: ",
		}
		selectExercise = Ask{
			Label:    "Exercise: ",
			Items:    exercises,
			AddLabel: "Add Exercise: ",
		}
		enterUnits = Ask{
			Label:    "Units: ",
			Items:    unitsOptions,
			AddLabel: "Add Units: ",
		}
		enterEffort = Ask{
			Label:    "Effort [0-100], 0 sleeping, 30 moderate, 50 hard, 70 very hard, 100 incredibly hard/failure: ",
			Validate: validateEffort,
		}

		enterSets = Ask{
			Label:    "Sets: ",
			Validate: validateInt,
		}

		enterVolume = Ask{
			Label: "Volume: ",
			Validate: func(s string) error {
				_, err = strconv.ParseFloat(s, 64)
				return err
			},
		}

		enterWeight = Ask{
			Label:    "Weight: ",
			Default:  "0",
			Validate: validateInt,
		}

		enterDuration = Ask{
			Label:     "Duration: ",
			Validate:  validateDuration,
			Default:   "00:00:00",
			AllowEdit: true,
		}
		enterFailure = Ask{
			IsConfirm: true,
			Label:     "Failure?",
		}
		enterConfirm = Ask{
			IsConfirm: true,
		}
		enterDone = Ask{
			Label:     "Done",
			IsConfirm: true,
		}

		rep        = lifting.Repetition{}
		previously lifting.Repetition

		// to insert
		toLoad = make([]lifting.Repetition, 0)
	)

	if len(args) == 0 {
		sessionDateString, err = enterDate.Run()
		handle(err)
	} else {
		sessionDateString = args[0]
	}
	rep.SessionDate, err = lifting.ParseSessionDateString(sessionDateString)

	handle(err)

	getLabel(storage, &rep, &selectLabel)

	recent, err = storage.GetByCategory(rep.Category, 10, 0)
	handle(err)

	exercises = make([]string, len(recent))
	unitsOptions = make([]string, len(recent))

	for i, r := range recent {
		repsByExercise[r.Exercise] = r
		exercises[i] = r.Exercise
		unitsOptions[i] = r.Units
	}

	for true {
		exercise, err = selectExercise.Run()
		previously = repsByExercise[exercise]

	INPUT_START:
		handle(err)
		if _, ok := repsByExercise[exercise]; ok {
			exercises = append(exercises, exercise)
		}
		rep.Exercise = exercise
		selectExercise.Items = exercises
		enterEffort.Default = strconv.Itoa(previously.Effort)
		enterSets.Default = strconv.Itoa(previously.Sets)
		enterVolume.Default = strconv.FormatFloat(previously.Volume, 'f', 2, 64)
		enterWeight.Default = strconv.Itoa(previously.Weight)
		enterDuration.Default = previously.Elapsed.String()
		enterFailure.Default = strconv.FormatBool(previously.Failure)

		// units
		rep.Units, err = enterUnits.Run()
		handle(err)

		// volume
		volume, err = enterVolume.Run()
		handle(err)
		rep.Volume, err = strconv.ParseFloat(volume, 64)
		handle(err)

		// sets
		sets, err = enterSets.Run()
		handle(err)
		rep.Sets, err = strconv.Atoi(sets)
		if rep.Sets < 1 {
			rep.Sets = 1
		}
		handle(err)

		// weight
		weight, err = enterWeight.Run()
		handle(err)
		rep.Weight, err = strconv.Atoi(weight)
		handle(err)

		// failure
		failure, err = enterFailure.Run()
		if err == promptui.ErrAbort {
			failure = "false"
		}
		rep.Failure, err = strconv.ParseBool(failure)
		handle(err)

		// duration
		duration, err = enterDuration.Run()
		handle(err)
		rep.Elapsed, err = civil.ParseTime(duration)
		handle(err)

		if rep.Failure {
			rep.Effort = 100
		} else {
			// effort
			effort, err = enterEffort.Run()
			handle(err)
			rep.Effort, err = strconv.Atoi(effort)
			handle(err)
		}

		// confirm
		handle(err)
		enterConfirm.Label = fmt.Sprintf(
			"Does the following look correct? %#v",
			rep,
		)

		confirmed, err := enterConfirm.Confirm()
		handle(err)
		if !confirmed {
			previously = rep

			enterDate.Default = rep.SessionDate.String()
			sessionDateString, err = enterDate.Run()
			handle(err)
			rep.SessionDate, err = lifting.ParseSessionDateString(sessionDateString)

			rep.Category, err = selectLabel.Run()
			handle(err)

			exercise, err = selectExercise.Run()
			handle(err)
			goto INPUT_START
		} else if err == nil {
			for i := 0; i < rep.Sets; i++ {
				toLoad = append(toLoad, rep)
			}
			done, err := enterDone.Confirm()
			handle(err)
			if done {
				break
			}
		} else {
			handle(err)
		}
	}
	err = storage.Load(toLoad)
	handle(err)
}

func getLabel(storage lifting.Storage, rep *lifting.Repetition, selectLabel *Ask) string {
	labels, err := storage.GetUniqueCategories()
	handle(err)
	selectLabel.Items = labels
	rep.Category, err = selectLabel.Run()
	handle(err)
	selectLabel.Items = append(selectLabel.Items, rep.Category)
	return rep.Category
}
