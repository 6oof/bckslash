package commands

import "github.com/6oof/bckslash/pkg/helpers"

type ExecFinishedMsg struct {
	Err     error
	Content string
}

type EmptyMsg struct{}

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

type ReturnHomeMsg struct{}
