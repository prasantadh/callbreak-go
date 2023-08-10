package cli

import "net/url"

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
)

type Response struct {
	Status `json:"status"`
	Data   any `json:"data"`
}

type Options struct {
	Server url.URL
}
