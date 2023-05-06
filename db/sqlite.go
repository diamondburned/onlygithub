package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub/db/sqlite"
	"libdb.so/onlygithub/onlygithub/auth"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed sqlite/schema.sql
var sqliteSchema string

// SQLite implements various database interfaces using SQLite.
type SQLite struct {
	db *sql.DB
	q  *sqlite.Queries
}

// NewSQLite creates a new SQLite database connection.
func NewSQLite(uri string) (*SQLite, error) {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(sqliteSchema); err != nil {
		return nil, errors.Wrap(err, "failed to execute schema")
	}

	return &SQLite{
		db: db,
		q:  sqlite.New(db),
	}, nil
}

// Close closes the database connection.
func (s *SQLite) Close() error {
	return s.db.Close()
}

// AsOAuthTokenService returns the database as an OAuthTokenService.
func (s *SQLite) AsOAuthTokenService() auth.OAuthTokenService {
	return (*sqliteOAuthTokenService)(s)
}

type sqliteOAuthTokenService SQLite

func (s *sqliteOAuthTokenService) SaveToken(ctx context.Context, token, provider string, oauthToken *oauth2.Token) error {
	err := s.q.SaveToken(ctx, sqlite.SaveTokenParams{
		Token:        token,
		Provider:     provider,
		AccessToken:  oauthToken.AccessToken,
		TokenType:    oauthToken.TokenType,
		RefreshToken: oauthToken.RefreshToken,
		ExpiresIn:    oauthToken.Expiry,
	})
	return sqliteErr(err)
}

func (s *sqliteOAuthTokenService) RetrieveToken(ctx context.Context, token, provider string) (*oauth2.Token, error) {
	v, err := s.q.RestoreToken(ctx, sqlite.RestoreTokenParams{
		Token:    token,
		Provider: provider,
	})
	if err != nil {
		return nil, sqliteErr(err)
	}
	return &oauth2.Token{
		AccessToken:  v.AccessToken,
		TokenType:    v.TokenType,
		RefreshToken: v.RefreshToken,
		Expiry:       v.ExpiresIn,
	}, nil
}

func sqliteErr(err error) error {
	return err
}
