package admission

import (
	"context"
	"errors"
)

type Response struct {
	Status int
	Body   string
}

type Handler struct {
	Service Service
}

func (h Handler) Publish(_ context.Context, schemaID string) Response {
	err := h.Service.Admit(context.Background(), schemaID)
	if err == nil {
		return Response{Status: 202, Body: "accepted"}
	}
	var permanent PermanentError
	if errors.As(err, &permanent) {
		return Response{Status: 422, Body: permanent.Error()}
	}
	return Response{Status: 503, Body: err.Error()}
}
