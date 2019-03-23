package main

import (
	"github.com/manifoldco/promptui"
)

// Ask is a wrapper for whatever readline/cli ui input library
type Ask struct {
	Label     string
	Items     []string
	AddLabel  string
	Default   string
	AllowEdit bool
	IsConfirm bool
	Validate  func(arg string) error
}

func readlineRun(ask *Ask) (string, error) {

	return "", nil

}

func promptuiRun(ask *Ask) (string, error) {
	if len(ask.Items) == 0 {
		p := promptui.Prompt{
			Label:     ask.Label,
			Default:   ask.Default,
			AllowEdit: ask.AllowEdit,
			IsConfirm: ask.IsConfirm,
			Validate:  ask.Validate,
		}
		val, err := p.Run()
		return val, err
	} else

	if ask.AddLabel != "" {
		p := promptui.SelectWithAdd{
			Label:    ask.Label,
			Validate: ask.Validate,
			Items:    ask.Items,

			AddLabel: ask.AddLabel,
		}
		_, val, err := p.Run()
		return val, err
	} 

	p := promptui.Select{
		Label: ask.Label,
		Items: ask.Items,
	}
	_, val, err := p.Run()
	return val, err
}

// Run actually promps the question on the CLI
func (ask *Ask) Run() (string, error) {
	return promptuiRun(ask)
}
// Confirm prompsthe user to confirm the ask
func (ask *Ask) Confirm() (bool, error) {
	_, err := promptuiRun(ask)

	if err == promptui.ErrAbort {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}
