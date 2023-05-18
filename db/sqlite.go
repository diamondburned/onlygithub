package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/db/sqlitec"

	_ "modernc.org/sqlite"

	sqlite "modernc.org/sqlite"
	libsqlite "modernc.org/sqlite/lib"
)

const sqliteSchemaSeparator = "\n-- MIGRATE --\n"

//go:embed sqlitec/schema.sql
var sqliteSchema string
var sqliteSchemaVersions = strings.Split(sqliteSchema, sqliteSchemaSeparator)

// SQLite implements various database interfaces using SQLite.
type SQLite struct {
	db  *sql.DB
	dbx *sqlx.DB
	q   *sqlitec.Queries
}

var _ Database = (*SQLite)(nil)

// NewSQLite creates a new SQLite database connection.
func NewSQLite(uri string) (*SQLite, error) {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(sqliteSchemaVersions[0]); err != nil {
		return nil, errors.Wrap(err, "failed to execute initial migration")
	}

	var userVersion int
	if err := db.QueryRow("PRAGMA user_version").Scan(&userVersion); err != nil {
		return nil, errors.Wrap(err, "failed to get user_version")
	}

	if userVersion < len(sqliteSchemaVersions) {
		for i, v := range sqliteSchemaVersions[userVersion:] {
			if _, err := db.Exec(v); err != nil {
				return nil, errors.Wrapf(err, "failed to execute migration %d", i+userVersion)
			}
		}
	}

	userVersionStr := strconv.Itoa(len(sqliteSchemaVersions))
	if _, err := db.Exec("PRAGMA user_version = " + userVersionStr); err != nil {
		return nil, errors.Wrap(err, "failed to set user_version")
	}

	return &SQLite{
		db:  db,
		dbx: sqlx.NewDb(db, "sqlite"),
		q:   sqlitec.New(db),
	}, nil
}

// Close closes the database connection.
func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) SaveToken(ctx context.Context, token, provider string, oauthToken *oauth2.Token) error {
	err := s.q.SaveToken(ctx, sqlitec.SaveTokenParams{
		Token:        token,
		Provider:     provider,
		AccessToken:  oauthToken.AccessToken,
		TokenType:    oauthToken.TokenType,
		RefreshToken: oauthToken.RefreshToken,
		ExpiresIn:    oauthToken.Expiry,
	})
	return sqliteErr(err)
}

func (s *SQLite) RetrieveToken(ctx context.Context, token, provider string) (*oauth2.Token, error) {
	v, err := s.q.RestoreToken(ctx, sqlitec.RestoreTokenParams{
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

func (s *SQLite) DeleteToken(ctx context.Context, token, provider string) error {
	err := s.q.DeleteToken(ctx, sqlitec.DeleteTokenParams{
		Token:    token,
		Provider: provider,
	})
	return sqliteErr(err)
}

func (s *SQLite) User(ctx context.Context, id onlygithub.GitHubID) (*onlygithub.User, error) {
	const q = `
		SELECT
			users.id, users.username, users.email, users.real_name, users.pronouns, users.avatar_url, users.joined_at, users.is_owner,
			user_tiers.price AS tier_price,
			user_tiers.is_one_time AS tier_is_one_time,
			user_tiers.is_custom_amount AS tier_is_custom_amount,
			user_tiers.started_at AS tier_started_at,
			user_tiers.renewed_at AS tier_renewed_at,
			tiers.id AS tier_id,
			tiers.name AS tier_name,
			tiers.description AS tier_description
		FROM users AS users -- https://github.com/kyleconroy/sqlc/issues/2271
		LEFT JOIN user_tiers ON users.id = user_tiers.user_id
		LEFT JOIN tiers ON user_tiers.tier_id = tiers.id
		WHERE users.id = ?
	`

	var u struct {
		ID                 string         `db:"id"`
		Username           string         `db:"username"`
		Email              string         `db:"email"`
		RealName           string         `db:"real_name"`
		Pronouns           string         `db:"pronouns"`
		AvatarUrl          string         `db:"avatar_url"`
		JoinedAt           time.Time      `db:"joined_at"`
		IsOwner            bool           `db:"is_owner"`
		TierPrice          sql.NullInt64  `db:"tier_price"`
		TierIsOneTime      sql.NullBool   `db:"tier_is_one_time"`
		TierIsCustomAmount sql.NullBool   `db:"tier_is_custom_amount"`
		TierStartedAt      sql.NullTime   `db:"tier_started_at"`
		TierRenewedAt      sql.NullTime   `db:"tier_renewed_at"`
		TierID             sql.NullString `db:"tier_id"`
		TierName           sql.NullString `db:"tier_name"`
		TierDescription    sql.NullString `db:"tier_description"`
	}

	if err := sqlx.GetContext(ctx, s.dbx, &u, q, string(id)); err != nil {
		return nil, sqliteErr(err)
	}

	user := &onlygithub.User{
		ID:        onlygithub.GitHubID(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		RealName:  u.RealName,
		Pronouns:  u.Pronouns,
		AvatarURL: u.AvatarUrl,
		JoinedAt:  u.JoinedAt,
		IsOwner:   u.IsOwner,
	}

	if u.TierPrice.Valid {
		user.Sponsorship = &onlygithub.Sponsorship{
			Price:          onlygithub.Cents(u.TierPrice.Int64),
			StartedAt:      u.TierStartedAt.Time,
			RenewedAt:      u.TierRenewedAt.Time,
			IsOneTime:      u.TierIsOneTime.Bool,
			IsCustomAmount: u.TierIsCustomAmount.Bool,
		}

		if u.TierID.Valid {
			user.Sponsorship.Tier = &onlygithub.Tier{
				ID:          onlygithub.GitHubID(u.TierID.String),
				Name:        u.TierName.String,
				Price:       onlygithub.Cents(u.TierPrice.Int64),
				Description: template.HTML(u.TierDescription.String),
			}
		}
	}

	return user, nil
}

func (s *SQLite) UpdateUser(ctx context.Context, user *onlygithub.User) error {
	err := s.q.UpdateUser(ctx, sqlitec.UpdateUserParams{
		ID:        string(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		RealName:  user.RealName,
		Pronouns:  user.Pronouns,
		AvatarUrl: user.AvatarURL,
	})
	return sqliteErr(err)
}

func (s *SQLite) Owner(ctx context.Context) (*onlygithub.User, error) {
	u, err := s.q.Owner(ctx)
	if err != nil {
		return nil, sqliteErr(err)
	}

	return &onlygithub.User{
		ID:        onlygithub.GitHubID(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		RealName:  u.RealName,
		Pronouns:  u.Pronouns,
		AvatarURL: u.AvatarUrl,
		JoinedAt:  u.JoinedAt,
		IsOwner:   u.IsOwner,
	}, nil
}

func (s *SQLite) MakeOwner(ctx context.Context, username string) error {
	err := s.q.MakeOwner(ctx, username)
	return sqliteErr(err)
}

func (s *SQLite) Image(ctx context.Context, id onlygithub.ID) (*onlygithub.ImageAsset, error) {
	a, err := s.q.ImageAsset(ctx, id.String())
	if err != nil {
		return nil, sqliteErr(err)
	}

	return &onlygithub.ImageAsset{
		Asset: onlygithub.Asset{
			ID:          id,
			Visibility:  onlygithub.Visibility(a.Visibility),
			MinimumCost: onlygithub.Cents(a.MinimumCost),
			LastUpdated: nullTimePtr(a.LastUpdated),
		},
		Filename: a.Filename,
	}, nil
}

func (s *SQLite) ImageData(ctx context.Context, id onlygithub.ID) (io.ReadCloser, error) {
	b, err := s.q.ImageData(ctx, id.String())
	if err != nil {
		return nil, sqliteErr(err)
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func (s *SQLite) UploadImage(ctx context.Context, req onlygithub.UploadImageRequest, r io.Reader) (*onlygithub.ImageAsset, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, sqliteErr(err)
	}

	id := onlygithub.GenerateID()

	updated, err := s.q.CreateImageAsset(ctx, sqlitec.CreateImageAssetParams{
		ID:          id.String(),
		Data:        b,
		Visibility:  string(req.Visibility),
		MinimumCost: int64(req.MinimumCost),
		Filename:    req.Filename,
	})
	if err != nil {
		return nil, sqliteErr(err)
	}

	return &onlygithub.ImageAsset{
		Asset: onlygithub.Asset{
			ID:          id,
			Visibility:  req.Visibility,
			MinimumCost: req.MinimumCost,
			LastUpdated: nullTimePtr(updated),
		},
		Filename: req.Filename,
	}, nil
}

func (s *SQLite) DeleteImage(ctx context.Context, id onlygithub.ID) error {
	err := s.q.DeleteImageAsset(ctx, id.String())
	return sqliteErr(err)
}

func (s *SQLite) SetImageVisibility(ctx context.Context, id onlygithub.ID, visibility onlygithub.Visibility) error {
	err := s.q.SetAssetVisibility(ctx, sqlitec.SetAssetVisibilityParams{
		ID:         id.String(),
		Visibility: string(visibility),
	})
	return sqliteErr(err)
}

func (s *SQLite) SiteConfig(ctx context.Context) (*onlygithub.SiteConfig, error) {
	cfg, _ := s.q.SiteConfig(ctx)
	if cfg != nil {
		var obj *onlygithub.SiteConfig
		if err := json.Unmarshal([]byte(cfg), &obj); err != nil {
			return nil, sqliteErr(err)
		}
		if obj != nil {
			return obj, nil
		}
	}

	return onlygithub.DefaultSiteConfig(), nil
}

func (s *SQLite) SetSiteConfig(ctx context.Context, cfg *onlygithub.SiteConfig) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return sqliteErr(err)
	}

	err = s.q.SetSiteConfig(ctx, []byte(b))
	return sqliteErr(err)
}

func (s *SQLite) UserConfig(ctx context.Context, id onlygithub.GitHubID) (*onlygithub.UserConfig, error) {
	cfg, err := s.q.UserConfig(ctx, string(id))
	if err == nil {
		var obj *onlygithub.UserConfig
		if err := json.Unmarshal([]byte(cfg), &obj); err != nil {
			return nil, sqliteErr(err)
		}
		if obj != nil {
			return obj, nil
		}
	}

	return onlygithub.DefaultUserConfig(), nil
}

func (s *SQLite) SetUserConfig(ctx context.Context, id onlygithub.GitHubID, cfg *onlygithub.UserConfig) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return sqliteErr(err)
	}

	err = s.q.SetUserConfig(ctx, sqlitec.SetUserConfigParams{
		UserConfig: []byte(b),
		ID:         string(id),
	})
	return sqliteErr(err)
}

func (s *SQLite) Tiers(ctx context.Context) ([]onlygithub.Tier, error) {
	v, err := s.q.Tiers(ctx)
	if err != nil {
		return nil, sqliteErr(err)
	}

	tiers := make([]onlygithub.Tier, len(v))
	for i, t := range v {
		tiers[i] = onlygithub.Tier{
			ID:          onlygithub.GitHubID(t.ID),
			Name:        t.Name,
			Price:       onlygithub.Cents(t.Price),
			Description: template.HTML(t.Description),
		}
	}

	return tiers, nil
}

func (s *SQLite) UpdateTiers(ctx context.Context, tiers []onlygithub.Tier) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return sqliteErr(err)
	}
	defer tx.Rollback()

	q := sqlitec.New(tx)

	if err := q.DeleteAllTiers(ctx); err != nil {
		return errors.Wrap(sqliteErr(err), "delete tiers")
	}

	for _, t := range tiers {
		if err := q.CreateTier(ctx, sqlitec.CreateTierParams{
			ID:          string(t.ID),
			Name:        t.Name,
			Price:       int64(t.Price),
			Description: string(t.Description),
		}); err != nil {
			return errors.Wrap(sqliteErr(err), "create tier")
		}
	}

	return sqliteErr(tx.Commit())
}

func sqliteErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return onlygithub.ErrNotFound
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code() {
		case libsqlite.SQLITE_CONSTRAINT:
			return errors.New("already exists")
		case libsqlite.SQLITE_TOOBIG:
			return errors.New("data too big")
		}
	}

	return onlygithub.WrapInternalError(err)
}

func nullTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
