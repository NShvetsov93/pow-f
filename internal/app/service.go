package app

import (
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
)

type Implementation struct {
	service service
}

type service interface {
	Solve(ctx context.Context) (string, error)
}

type Response struct {
	Phrase string
}

func New(s service) *Implementation {
	return &Implementation{
		service: s,
	}
}

func (i *Implementation) Solve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "app/Solve")
	defer span.Finish()

	phrase, err := i.service.Solve(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res := &Response{
		Phrase: phrase,
	}

	if err := jsoniter.NewEncoder(w).Encode(res); err != nil {
		err = fmt.Errorf("error encoding response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
