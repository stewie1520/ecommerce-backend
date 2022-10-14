package core

import "net/http"

type Handler func(request *http.Request, response http.ResponseWriter) error
