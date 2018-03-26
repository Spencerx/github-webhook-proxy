package main

import (
	"net/http"
)

type PendingRequest struct {
	Headers http.Header
	UUID    string
	Body    string
	Method  string
}

type RequestCache struct {
	Requests []*PendingRequest
}
