package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"libdb.so/hserve"
	"libdb.so/onlygithub/db"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/routes"
	"libdb.so/onlygithub/internal/auth"
	"libdb.so/onlygithub/internal/reflectutil"
)

// Flags.
var (
	httpAddr = "localhost:8080"
	mkOwner  = ""
)

// Environment variables.
var (
	databaseURL = requireEnv("DATABASE_URL")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	flag.StringVar(&httpAddr, "http", httpAddr, "HTTP address")
	flag.StringVar(&mkOwner, "mkowner", mkOwner, "Make a user the owner of the site")
	flag.Parse()

	switch {
	case mkOwner != "":
		makeOwner(ctx, mkOwner)
	default:
		serve()
	}
}

func makeOwner(ctx context.Context, username string) {
	db, err := openDB()
	if err != nil {
		log.Fatalln("error opening database:", err)
	}
	defer db.Close()

	owner, err := db.Owner(ctx)
	if err == nil {
		if owner.Username == username {
			log.Println("site already has", username, "as owner")
			return
		}
		log.Fatalln("site already has an owner:", owner.Username)
	}

	log.Println("making", username, "the owner of the site")

	if err := db.MakeOwner(ctx, username); err != nil {
		log.Fatalln("error making", username, "the owner:", err)
	}
}

func serve() {
	db, err := openDB()
	if err != nil {
		log.Fatalln("error opening database:", err)
	}
	defer db.Close()

	ghConfig := auth.GitHubConfig{
		ID:          os.Getenv("GITHUB_CLIENT_ID"),
		Secret:      os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL: os.Getenv("GITHUB_REDIRECT_URL"),
	}
	reflectutil.AssertNoZeroFields(ghConfig, func(field string) {
		log.Fatalln("missing required GitHub config:", field)
	})

	ghOAuth := auth.NewGitHubAuthorizer(ghConfig, db)

	h := routes.New(frontend.Deps{
		GitHubOAuth: ghOAuth,
		Config:      db,
		Images:      db,
		Users:       db,
		Tiers:       db,
		Posts:       db,
	})

	log.Println("listening on", httpAddr)
	hserve.MustListenAndServe(context.Background(), httpAddr, h)
}

func openDB() (db.Database, error) {
	var store db.Database

	databaseURL, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalln("error parsing DATABASE_URL:", err)
	}

	switch databaseURL.Scheme {
	case "sqlite", "sqlite3":
		d, err := db.NewSQLite(databaseURL.Host + databaseURL.Path)
		if err != nil {
			log.Fatalln("error opening SQLite database:", err)
		}
		store = d
	default:
		log.Fatalln("invalid DATABASE_URL scheme:", databaseURL.Scheme)
	}

	return store, nil
}

func requireEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalln("missing required environment variable:", name)
	}
	return v
}
