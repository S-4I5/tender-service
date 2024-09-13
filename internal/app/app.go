package app

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"tender-service/internal/config"
	"tender-service/internal/middleware"
)

type App struct {
	provider *serviceProvider
	server   http.Server
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	a := &App{}
	a.provider = newServiceProvider(cfg)

	err := a.setup(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) setup(ctx context.Context) error {

	funcs := []func(context.Context) error{
		a.runMigrationsForPostgres,
		a.setupHttpServer,
	}

	for _, f := range funcs {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) setupHttpServer(ctx context.Context) error {

	tenderMux := http.NewServeMux()
	tenderMux.HandleFunc("POST /new", a.provider.TenderController().PostNewTender(ctx))
	tenderMux.HandleFunc("GET /my", a.provider.TenderController().GetUserTenders(ctx))
	tenderMux.HandleFunc("GET /{tenderId}/status", a.provider.TenderController().GetTenderStatus(ctx))
	tenderMux.HandleFunc("PUT /{tenderId}/status", a.provider.TenderController().PutTenderStatus(ctx))
	tenderMux.HandleFunc("PATCH /{tenderId}/edit", a.provider.TenderController().PatchTender(ctx))
	tenderMux.HandleFunc("PUT /{tenderId}/rollback/{version}", a.provider.TenderController().PutTenderRollback(ctx))

	bidMux := http.NewServeMux()
	bidMux.HandleFunc("POST /new", a.provider.BidController().PostNewBid(ctx))
	bidMux.HandleFunc("GET /my", a.provider.BidController().GetUserBids(ctx))
	bidMux.HandleFunc("GET /{tenderId}/list", a.provider.BidController().GetTenderBids(ctx))
	bidMux.HandleFunc("GET /{bidId}/status", a.provider.BidController().GetBidStatus(ctx))
	bidMux.HandleFunc("PUT /{bidId}/status", a.provider.BidController().PutBidStatus(ctx))
	bidMux.HandleFunc("PATCH /{bidId}/edit", a.provider.BidController().PatchBid(ctx))
	bidMux.HandleFunc("PUT /{bidId}/submit_decision", a.provider.BidController().PutBidSubmitDecision(ctx))
	bidMux.HandleFunc("PUT /{bidId}/feedback", a.provider.BidController().PutBidFeedback(ctx))
	bidMux.HandleFunc("PUT /{bidId}/rollback/{version}", a.provider.BidController().PutBidRollback(ctx))
	bidMux.HandleFunc("GET /{tenderId}/reviews", a.provider.BidController().GetBidReviews(ctx))

	api := http.NewServeMux()

	api.Handle("GET /ping", a.provider.PingController().GetPing(ctx))
	api.Handle("GET /tenders", a.provider.TenderController().GetTenders(ctx))

	api.Handle("/bids/", http.StripPrefix("/bids", bidMux))
	api.Handle("/tenders/", http.StripPrefix("/tenders", tenderMux))

	main := http.NewServeMux()

	main.Handle("/api/", http.StripPrefix("/api", api))

	log.Println("starting on:" + a.provider.config.Server.Address)

	a.server = http.Server{
		Addr:    a.provider.config.Server.Address,
		Handler: middleware.GetLoggerMiddleware(main),
	}
	return nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop() error {
	return a.server.Shutdown(context.Background())
}

func (a *App) runMigrationsForPostgres(_ context.Context) error {
	log.Println("running migrations in:", a.provider.config.Postgres.MigrationsDir)

	conn := fmt.Sprintf(a.provider.config.Postgres.Conn)

	dsn := flag.String("dsn", conn, "PostgreSQL")

	sql, err := goose.OpenDBWithDriver("postgres", *dsn)
	if err != nil {
		return err
	}

	err = goose.Up(sql, a.provider.config.Postgres.MigrationsDir)
	if err != nil {
		return err
	}

	return nil
}
