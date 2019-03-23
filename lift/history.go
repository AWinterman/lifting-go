package main

import (
	"github.com/awinterman/lifting"
	"fmt"
	"bytes"
	"github.com/spf13/cobra"
	"text/template"
)

func history(cmd *cobra.Command, args []string) {
	rs, err := storage.GetLast(100, 0)

	if err != nil {
		panic(err)
	}

	for _, r := range rs {
		fmt.Println(r)
	}
}

const repetitionTemplate = `
	{{.SessionDate}} .Category .Exercise {{- if .Sets > 1}} {{.Sets}} {{.Volume}}
	{{- if .Weight > 0}}
	   {{.Weight}}
	{{.Units}}
	{{.Elapsed}}
	{{.Failure}}
	`

func repr(r lifting.Repetition) string {
	t, err := template.New("workout_element").Parse(repetitionTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	t.Execute(buf, r)
	return buf.String()
}
