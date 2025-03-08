package model

type APIStatus int

const (
	Success APIStatus = iota
	Failure
)

var stateName = map[APIStatus]string{
	Success: "success",
	Failure: "failure",
}
