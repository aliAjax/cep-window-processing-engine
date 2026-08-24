package admission

import "errors"

type Response struct {
	Status int
	Body   string
}

type Handler struct {
	Service Service
}

func (h Handler) Publish(schemaID string) Response {
	err := h.Service.Admit(schemaID)
	if err == nil {
		return Response{Status: 202, Body: "accepted"}
	}
	var permanent PermanentError
	if errors.As(err, &permanent) {
		return Response{Status: 422, Body: permanent.Error()}
	}
	return Response{Status: 503, Body: err.Error()}
}
