package main

import (
	"github.com/awinterman/lifting/sqlite"
	"github.com/spf13/cobra"
)

var storage, err = sqlite.CreateStorage(".lift.sqlite", nil)

func main() {
	if (err != nil) {
		panic(err)
	}
	var root = &cobra.Command{Use: "lift", Short: "Log, view, or edit workouts"}
	var add = &cobra.Command{
		Use:   "add",
		Run:   logWorkout,
		Short: "Log a workout",
	}
	var history = &cobra.Command{
		Use:   "history",
		Run:   history,
		Short: "View recent workouts",
	}

	root.AddCommand(add)
	root.AddCommand(history)
	root.Execute()
}
