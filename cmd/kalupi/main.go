package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"

	"github.com/sf9v/kalupi/account"
	accountsvc "github.com/sf9v/kalupi/account/service"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/ledger"
	"github.com/sf9v/kalupi/postgres"
	"github.com/sf9v/kalupi/transaction"
)

const (
	defaultPort = "8000"
	defaultDSN  = "postgres://kalupi:kalupi@localhost:5432/kalupi?sslmode=disable"
)

func main() {
	var (
		addr     = envString("PORT", defaultPort)
		dsn      = envString("DSN", defaultDSN)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		ctx      = context.Background()
	)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	// migrate the database
	err = postgres.Migrate(db)
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	var (
		ledgerRepo  = postgres.NewLedgerRepository(db)
		accountRepo = postgres.NewAccountRepository(db)
		balRepo     = postgres.NewBalanceRepository(db)
		xactRepo    = postgres.NewXactRepository(db)
	)

	ls := ledger.NewService(ledgerRepo)

	// create cash ledgers
	err = ls.CreateCashLedgers(ctx)
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	bs := balance.NewService(balRepo)

	var as account.Service
	as = accountsvc.New(accountRepo, bs)
	as = account.NewLoggingService(logger, as)

	var xs transaction.Service
	xs = transaction.NewService(accountRepo, ledgerRepo, xactRepo, balRepo)
	xs = transaction.NewLoggingService(logger, xs)

	httpLogger := log.With(logger, "component", "http")

	mux := chi.NewMux()

	mux.Mount("/accounts", account.NewHandler(as, httpLogger))
	mux.Mount("/t", transaction.NewHandler(xs, httpLogger))

	srvr := &http.Server{
		Addr:           *httpAddr,
		Handler:        accessControl(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- srvr.ListenAndServe()
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
