package model

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	logging "fetch-me-if-you-read-me/logger"
)

//go:embed _migrations/*.sql
var migrations embed.FS

var (
	insertImage = strings.Join([]string{
		"INSERT INTO mafiyrm.images(",
		"  used_in",
		")",
		"VALUES ($1)",
		"ON CONFLICT ON CONSTRAINT images_pkey",
		"DO UPDATE",
		"SET",
		"  last_update_date = CURRENT_TIMESTAMP",
		"WHERE images.used_in = $1",
		"RETURNING id::varchar AS image_fk",
	}, " ")
	insertWhoIsFetching = strings.Join([]string{
		"INSERT INTO mafiyrm.who(",
		"  remote_addr,",
		"  meta",
		")",
		"VALUES ($1, $2::jsonb)",
		"ON CONFLICT ON CONSTRAINT who_pkey",
		"DO UPDATE",
		"SET",
		"  last_update_date = CURRENT_TIMESTAMP",
		"WHERE who.remote_addr = $1",
		"RETURNING id::varchar AS who_fk",
	}, " ")
	boundWhoIsFetchingWithImage = strings.Join([]string{
		"INSERT INTO mafiyrm.images_accessed(",
		"  image_fk,",
		"  who_fk",
		")",
		"VALUES ($1, $2)",
	}, " ")
)

type PostgresqlConfigurations struct {
	Administrator         *string
	AdministratorPassword *string
	Host                  *string
	Username              *string
	Password              *string
	Database              *string
	Threads               *int
	ApplicationName       string
	Schema                *string
	MigrationTable        *string
}

type Model struct {
	logger          *zap.SugaredLogger
	keepAliveTicker *time.Ticker
	keepAliveDone   chan bool

	connectionString         string
	postgresqlConfigurations *PostgresqlConfigurations
	pool                     *pgxpool.Pool
	txOpts                   *pgx.TxOptions
}

func (model *Model) ImageFetched(imageFk uuid.UUID, remoteAddr string, meta map[string][]string) error {
	model.logger.Debugf("Storing %s remote address for %s imageFk", remoteAddr, imageFk)
	metaJSON, err := json.Marshal(meta)
	if err != nil {

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tx, err := model.pool.BeginTx(ctx, *model.txOpts)
	if err != nil {

		return err
	}

	defer tx.Rollback(ctx)

	var whoFk string
	if err := tx.QueryRow(ctx, insertWhoIsFetching, remoteAddr, metaJSON).Scan(&whoFk); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, boundWhoIsFetchingWithImage, imageFk, whoFk); err != nil {
		return err
	}

	tx.Commit(ctx)

	select {
	case <-ctx.Done():
		model.logger.Errorf("Registering image %s fetch from %s went in error: %s", imageFk, remoteAddr, ctx.Err().Error())
		return ctx.Err()
	default:
		model.logger.Infof("Registering image %s fetch from %s done", imageFk, remoteAddr)

		return nil
	}
}

func (model *Model) PrepareImage(usedIn string) (*uuid.UUID, error) {
	model.logger.Debugf("Creating image reference used in %s", usedIn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tx, err := model.pool.BeginTx(ctx, *model.txOpts)
	if err != nil {

		return nil, err
	}

	defer tx.Rollback(ctx)

	var imageFk string
	if err := tx.QueryRow(ctx, insertImage, usedIn).Scan(&imageFk); err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	select {
	case <-ctx.Done():
		model.logger.Errorf("Insert image for %s went in error: %s", usedIn, ctx.Err().Error())
		return nil, ctx.Err()
	default:
		model.logger.Infof("Insert image for %s done", usedIn)

		imageFkUUID, err := uuid.Parse(imageFk)
		if err != nil {
			return nil, err
		}
		return &imageFkUUID, nil
	}
}

func (model *Model) Dispose() {
	model.keepAliveDone <- true
	model.keepAliveTicker.Stop()
}

func New(logger *logging.Logger, postgresqlConfigurations *PostgresqlConfigurations) (*Model, error) {
	toReturn := &Model{
		postgresqlConfigurations: postgresqlConfigurations,
		connectionString: fmt.Sprintf(strings.Join([]string{
			"postgres://%s:%s@%s:5432/%s?",
			"application_name=%s",
			"&connect_timeout=20",
			"&pool_max_conns=%d",
			"&pool_min_conns=%d",
		}, ""),
			*postgresqlConfigurations.Username,
			*postgresqlConfigurations.Password,
			*postgresqlConfigurations.Host,
			*postgresqlConfigurations.Database,
			postgresqlConfigurations.ApplicationName,
			*postgresqlConfigurations.Threads,
			*postgresqlConfigurations.Threads,
		),
		txOpts: &pgx.TxOptions{
			IsoLevel:       pgx.ReadUncommitted,
			DeferrableMode: pgx.NotDeferrable,
			AccessMode:     pgx.ReadWrite,
		},
		logger: logger.Log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := toReturn.migrate(ctx); err != nil {
		return nil, err
	}

	if err := toReturn.initPool(ctx); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		logger.Log.Errorf("Model initialization went in error: %s", ctx.Err().Error())
		return nil, ctx.Err()
	default:
		logger.Log.Info("Model initialization done")

		return toReturn, nil
	}
}

func setupDatabase(username, password, database, schema string) (string, error) {

	//TODO sanitize
	sqlComment := "--"
	if strings.Contains(username, sqlComment) ||
		strings.Contains(password, sqlComment) ||
		strings.Contains(database, sqlComment) ||
		strings.Contains(schema, sqlComment) {
		return " ", errors.New("invalid characters in username, password or database")
	}

	return strings.Join([]string{
		"DO",
		"$do$",
		"BEGIN",
		"  IF (",
		"    SELECT COUNT(*)",
		"    FROM pg_catalog.pg_user",
		"    WHERE usename = '" + username + "'",
		"  ) = 0 THEN",
		"    CREATE ROLE \"" + username + "\" LOGIN PASSWORD '" + password + "';",
		"  END IF;",
		"END",
		"$do$;",
		"CREATE SCHEMA IF NOT EXISTS " + schema + " AUTHORIZATION \"" + username + "\";",
		"GRANT CONNECT ON DATABASE \"" + database + "\" TO \"" + username + "\";",
		"GRANT USAGE ON ALL SEQUENCES IN SCHEMA " + schema + " TO \"" + username + "\";",
		"GRANT CREATE ON SCHEMA " + schema + " TO \"" + username + "\";",
		"GRANT SELECT, INSERT, UPDATE, REFERENCES ON ALL TABLES IN SCHEMA " + schema + " TO \"" + username + "\";",
		"GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO \"" + username + "\";",
		"GRANT USAGE ON SCHEMA public TO \"" + username + "\";",
		"GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA " + schema + " TO \"" + username + "\";",
		"GRANT USAGE ON SCHEMA " + schema + " TO \"" + username + "\";",
	}, " "), nil
}

func (m *Model) migrate(ctx context.Context) error {
	adminConn, err := pgx.Connect(ctx, fmt.Sprintf(strings.Join([]string{
		"postgres://%s:%s@%s:5432/%s?",
		"application_name=%s",
		"&connect_timeout=20",
	}, ""),
		*m.postgresqlConfigurations.Administrator,
		*m.postgresqlConfigurations.AdministratorPassword,
		*m.postgresqlConfigurations.Host,
		*m.postgresqlConfigurations.Database,
		m.postgresqlConfigurations.ApplicationName+"-admin",
	))
	if err != nil {
		return err
	}

	defer adminConn.Close(ctx)

	sourceInstance, err := httpfs.New(http.FS(migrations), "_migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("httpfs", sourceInstance, fmt.Sprintf(strings.Join([]string{
		"postgres://%s:%s@%s:5432/%s?",
		"application_name=%s",
		"&connect_timeout=20",
		"&x-migrations-table=%s.%s",
		"&sslmode=disable",
	}, ""),
		*m.postgresqlConfigurations.Administrator,
		*m.postgresqlConfigurations.AdministratorPassword,
		*m.postgresqlConfigurations.Host,
		*m.postgresqlConfigurations.Database,
		m.postgresqlConfigurations.ApplicationName,
		*m.postgresqlConfigurations.Schema,
		*m.postgresqlConfigurations.MigrationTable,
	))
	if err != nil {
		return err
	}

	setupDatabaseStr, err := setupDatabase(*m.postgresqlConfigurations.Username,
		*m.postgresqlConfigurations.Password, *m.postgresqlConfigurations.Database, *m.postgresqlConfigurations.Schema)
	if err != nil {
		return err
	}

	_, err = adminConn.Exec(ctx, setupDatabaseStr)
	if err != nil {
		return err
	}

	if err := migrator.Up(); errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info(err)
	} else if err != nil {

		return err
	}

	sourceErr, databaseErr := migrator.Close()
	if sourceErr != nil {
		return sourceErr
	}

	if databaseErr != nil {
		return databaseErr
	}

	return nil
}

func (model *Model) CheckStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := model.pool.Ping(ctx)

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		errorStr := ctx.Err().Error()
		model.logger.Errorf("Connection check to pg went in error: %s", errorStr)
		return ctx.Err()
	default:
		model.logger.Debug("Still connected to postgresql")
		return nil
	}
}

func (model *Model) keepAlive() {
	for {
		select {
		case <-model.keepAliveDone:
			return
		case _ = <-model.keepAliveTicker.C:
			err := model.CheckStatus()

			if err != nil {

				panic(err.Error())
			}
		}
	}
}

func (model *Model) initPool(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, model.connectionString)
	if err != nil {
		return err
	}

	model.pool = pool
	model.keepAliveTicker = time.NewTicker(time.Minute)
	model.keepAliveDone = make(chan bool)
	go model.keepAlive()
	return nil
}
