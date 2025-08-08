package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/EmreZURNACI/apistack.git/domain"
	_ "github.com/lib/pq"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type PostgresHandler struct {
	Db     *sql.DB
	Tracer trace.Tracer
}

func GetPostgresHandler(tracer trace.Tracer) (*PostgresHandler, error) {

	driverName, err := otelsql.Register("postgres",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
	)
	if err != nil {
		return nil, err
	}

	var dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable",
		os.Getenv("HOST"), os.Getenv("PORT"),
		os.Getenv("USER"), os.Getenv("PASSWORD"),
		os.Getenv("DB"))

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		zap.L().Error("Error connecting to postgres", zap.Error(err))
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		zap.L().Error("Error pinging postgres", zap.Error(err))
		return nil, err
	}

	if err := otelsql.RecordStats(db); err != nil {
		return nil, err
	}

	return &PostgresHandler{
		Db:     db,
		Tracer: tracer,
	}, nil
}

func (h *PostgresHandler) GetActors(ctx context.Context, search string, offset, limit int, orderBy bool) ([]domain.Actor, error) {

	ctx, span := h.Tracer.Start(ctx, "GetUsers")
	defer span.End()

	var query = "SELECT * FROM public.actor"

	if search != "" {
		query += " WHERE first_name ILIKE '%" + search + "%' OR last_name ILIKE '%" + search + "%'"
	}

	if orderBy {
		query += " ORDER BY actor_id DESC"
	}

	if offset > 0 {
		query += " OFFSET " + strconv.Itoa(offset)
	}

	if limit > 0 {
		query += " LIMIT " + strconv.Itoa(limit)
	}

	rows, err := h.Db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Error("No record found", zap.Error(err))
			return nil, errors.New("kayıt bulunamadı")
		}
		zap.L().Error("Error executing query", zap.Error(err))
		return nil, errors.New("sorgu çalıştırılırken hata ile karşılaşıldı")
	}

	defer func() {
		if err := rows.Close(); err != nil {
			zap.L().Error("Error closing rows", zap.Error(err))
		}
	}()

	var (
		actors []domain.Actor
		actor  domain.Actor
	)

	for rows.Next() {
		if err := rows.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName, &actor.LastUpdate); err != nil {
			zap.L().Error("Error scanning row", zap.Error(err))
			return nil, errors.New("satır scan edilirken hata ile karşılaşıldı")
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func (h *PostgresHandler) CreateActor(ctx context.Context, firstName, lastName string) (int64, error) {

	ctx, span := h.Tracer.Start(ctx, "CreateActor")
	defer span.End()

	tx, err := h.Db.BeginTx(ctx, nil)

	if err != nil {
		zap.L().Error("Error starting transaction", zap.Error(err))
		return -1, errors.New("error starting transaction")
	}

	var count int64
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM public.actor WHERE first_name = $1 OR  last_name = $2", firstName, lastName).Scan(&count)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return -1, err
	}

	if count >= 1 {
		zap.L().Error("Bu bilgilere ait kullanıcı bulunmakta", zap.Error(err))
		return -1, errors.New("bu bilgilere ait kullanıcı bulunmakta")
	}

	var query = "INSERT INTO public.actor (first_name,last_name) VALUES ($1,$2)"

	_, err = tx.ExecContext(ctx, query, firstName, lastName)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return -1, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				zap.L().Error("Error rolling back", zap.Error(err))
				return
			}
		}
	}()

	if err := tx.Commit(); err != nil {
		zap.L().Error("Error committing transaction", zap.Error(err))
		return -1, err
	}

	var id int64
	err = h.Db.QueryRowContext(ctx, "SELECT actor_id FROM public.actor WHERE first_name = $1 AND last_name = $2", firstName, lastName).Scan(&id)
	if err != nil {
		zap.L().Error("Error scanning row", zap.Error(err))
		return -1, err
	}

	return id, nil
}

func (h *PostgresHandler) DeleteActor(ctx context.Context, id string) error {

	ctx, span := h.Tracer.Start(ctx, "DeleteActor")
	defer span.End()

	tx, err := h.Db.BeginTx(ctx, nil)

	if err != nil {
		zap.L().Error("Error starting transaction", zap.Error(err))
		return err
	}

	var count int64
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM public.actor WHERE actor_id= $1", id).Scan(&count)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return err
	}

	if count <= 0 {
		zap.L().Error("Bu id'ye ait kullanıcı bulunmamaktadır", zap.String("id", id))
		return errors.New("bu id'ye ait kullanıcı bulunmamaktadır")
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM public.actor WHERE actor_id = $1", id); err != nil {
		zap.L().Error("Error deleting actor", zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				zap.L().Error("Error rolling back", zap.Error(err))
				return
			}
		}
	}()

	if err := tx.Commit(); err != nil {
		zap.L().Error("Error committing transaction", zap.Error(err))
		return err
	}

	return nil
}

func (h *PostgresHandler) GetActor(ctx context.Context, id string) (*domain.Actor, error) {

	ctx, span := h.Tracer.Start(ctx, "GetActor")
	defer span.End()

	row := h.Db.QueryRowContext(ctx, "SELECT * FROM public.actor WHERE actor_id = $1", id)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			zap.L().Error("No record with this id", zap.Error(row.Err()))
			return nil, row.Err()
		}
		zap.L().Error("Error querying row", zap.Error(row.Err()))
		return nil, row.Err()
	}

	var actor domain.Actor
	if err := row.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName, &actor.LastUpdate); err != nil {
		zap.L().Error("Error scanning row", zap.Error(err))
		return nil, err
	}

	return &domain.Actor{
		ActorID:    actor.ActorID,
		FirstName:  actor.FirstName,
		LastName:   actor.LastName,
		LastUpdate: actor.LastUpdate,
	}, nil
}

func (h *PostgresHandler) UpdateActor(ctx context.Context, id, firstname, lastname string) error {

	ctx, span := h.Tracer.Start(ctx, "UpdateActor")
	defer span.End()

	tx, err := h.Db.BeginTx(ctx, nil)

	if err != nil {
		zap.L().Error("Error starting transaction", zap.Error(err))
		return err
	}

	var count int64
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM public.actor WHERE actor_id = $1", id).Scan(&count)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return err
	}

	if count <= 0 {
		zap.L().Error("bu id'ye ait veri bulunmamaktadır")
		return errors.New("bu id'ye ait veri bulunmamaktadır")
	}

	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM public.actor WHERE first_name = $1 OR  last_name = $2", firstname, lastname).Scan(&count)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return err
	}
	if count >= 1 {
		zap.L().Error("bu bilgilere ait kullanıcı bulunmakta")
		return errors.New("bu bilgilere ait kullanıcı bulunmakta")
	}

	_, err = tx.ExecContext(ctx, "UPDATE public.actor SET first_name = $1, last_name = $2 WHERE actor_id = $3", firstname, lastname, id)
	if err != nil {
		zap.L().Error("Error updating actor", zap.Error(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("Error committing transaction", zap.Error(err))
		return err
	}

	zap.L().Info("Updated")
	return nil

}
