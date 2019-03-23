package main

import (
	"cloud.google.com/go/civil"
	"errors"
	"github.com/manifoldco/promptui"
	"log"
	"strconv"
)

func handle(err error) {
	if err != nil {
		if err == promptui.ErrAbort {
			log.Panic("aborted", err)
		}

		log.Panic(err)
	}
}

func validateInt(arg string) error {
	_, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	return nil
}

func validateEffort(effort string) error {
	i, err := strconv.Atoi(effort)
	if err != nil {
		return err
	}

	if i < 0 || i > 100 {
		return errors.New("Must be between 0 and 100")
	}

	return nil
}

func validateDuration(duration string) error {
	_, err := civil.ParseTime(duration)
	return err
}
