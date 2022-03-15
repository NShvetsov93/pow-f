package main

import (
	"context"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"

	"pow-f/cmd/config"
	"pow-f/internal/app"
	"pow-f/internal/pkg/back"
	"pow-f/internal/pkg/solve"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	spew.Dump(cfg)

	backService := back.New(cfg.BackAddr, cfg.BackTimeout)
	solveService := solve.New(backService, cfg.TargetBits)
	impl := app.New(solveService)

	public := chi.NewRouter()

	public.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		impl.Solve(writer, request)
	})

	if err := http.ListenAndServe(":8080", public); err != nil {
		log.Fatal(ctx, err)
	}
}
