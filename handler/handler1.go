package handler

import (
	"hiteshkotian/ssl-tunnel/logging"
)

type CustomHandler struct {
}

func (handler *CustomHandler) HandleRequest(request *Request) error {
	logging.Info("Proxy request handled by handler 1")
	return nil
}
