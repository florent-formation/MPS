package router

import "net/http"

type Handlers map[string]http.HandlerFunc
