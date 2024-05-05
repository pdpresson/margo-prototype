package models

type FileState int

const (
	New FileState = iota + 1
	Updated
	Unchanged
)

const (
	KindSolution = "solutions"
	KindInstance = "instances"
)

type DesiredState struct {
	FileName         string
	CurrentStatePath string
	SourceFile       string
	State            FileState
}
