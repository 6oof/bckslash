package commands

import "github.com/6oof/bckslash/pkg/helpers"

type ExecFinishedMsg struct {
	Content string
}

type ExecStartMsg struct{}

type ProjectFoundMsg struct {
	Project helpers.Project
}

type ProjectViewData struct {
	GitLog string
}

type ProgramErrMsg struct {
	Err error
}

type ProjectListChangedMsg struct {
	ProjectList []helpers.Project
}
