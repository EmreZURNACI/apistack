package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/EmreZURNACI/apistack/domain"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresHandler struct {
	Db     *gorm.DB
	Tracer trace.Tracer
}

func GetPostgresHandler(tracer trace.Tracer) (*PostgresHandler, error) {

	var dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable",
		os.Getenv("HOST"), os.Getenv("PORT"),
		os.Getenv("USER"), os.Getenv("PASSWORD"),
		os.Getenv("DB"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("failed to connect to gorm postgres database")
		return nil, err

	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		zap.L().Error("failed to tracing gorm postgres database")
		return nil, err
	}

	return &PostgresHandler{
		Db:     db,
		Tracer: tracer,
	}, nil
}

func (h *PostgresHandler) GetActors(ctx context.Context, search string, offset, limit int, orderBy bool) ([]domain.Actor, error) {
	ctx, span := h.Tracer.Start(ctx, "GetActors")
	defer span.End()

	db := h.Db.WithContext(ctx).Table("actor")

	if search != "" {
		db = db.Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if orderBy {
		db = db.Order("actor_id DESC")
	}

	if offset > 0 {
		db = db.Offset(offset)
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	var actors []domain.Actor
	if err := db.Find(&actors).Error; err != nil {
		zap.L().Error("failed to query actors", zap.Error(err))
		return nil, errors.New("aktörler getirilirken bir sorun oluştu")
	}

	if len(actors) == 0 {
		zap.L().Info("kayıtlı aktör bulunamadı")
		return nil, errors.New("kayıtlı aktör bulunamadı")
	}

	return actors, nil
}

func (h *PostgresHandler) CreateActor(ctx context.Context, firstName, lastName string) (int64, error) {
	ctx, span := h.Tracer.Start(ctx, "CreateActor")
	defer span.End()

	tx := h.Db.WithContext(ctx).Table("actor").Begin()
	if tx.Error != nil {
		zap.L().Error("failed to start transaction", zap.Error(tx.Error))
		return 0, errors.New("transaction başlatılamadı")
	}

	defer func() {
		//recover == panic durumlarını yakalar
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var exist domain.Actor
	err := tx.Where("first_name = ? AND last_name = ?", firstName, lastName).
		First(&exist).Error

	if err == nil {
		tx.Rollback()
		return -1, errors.New("bu bilgilere ait kullanıcı zaten mevcut")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		zap.L().Error("veritabanı sorgu hatası", zap.Error(err))
		return -1, errors.New("sorgu hatası")
	}

	var count int64
	tx.Select("actor_id").Order("actor_id DESC").Limit(1).Scan(&count)
	actor := domain.Actor{ActorID: count + 1, FirstName: firstName, LastName: lastName}
	if err := tx.Create(&actor).Error; err != nil {
		tx.Rollback()
		zap.L().Error("kayıt eklenirken hata oluştu", zap.Error(err))
		return -1, errors.New("kayıt eklenirken hata oluştu")
	}

	if err := tx.Commit().Error; err != nil {
		zap.L().Error("transaction commit hatası", zap.Error(err))
		return -1, err
	}

	return actor.ActorID, nil
}

func (h *PostgresHandler) DeleteActor(ctx context.Context, id string) error {
	ctx, span := h.Tracer.Start(ctx, "DeleteActor")
	defer span.End()

	tx := h.Db.WithContext(ctx).Table("actor").Begin()
	if tx.Error != nil {
		zap.L().Error("transaction başlatılamadı", zap.Error(tx.Error))
		return errors.New("transaction başlatılamadı")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var actor domain.Actor
	if err := tx.Where("actor_id = ?", id).First(&actor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("bu id'ye ait kullanıcı bulunmamaktadır")
		}
		tx.Rollback()
		zap.L().Error("actor sorgusunda hata", zap.Error(err))
		return errors.New("ilgili id'li actor bulunurken hata oluştu")
	}

	if err := tx.Where("actor_id = ?", actor.ActorID).Delete(&actor).Error; err != nil {
		tx.Rollback()
		zap.L().Error("actor silinirken hata oluştu", zap.Error(err))
		return errors.New("actor silme sorgusu çalıştırılırken hata oluştu")
	}

	if err := tx.Commit().Error; err != nil {
		zap.L().Error("transaction commit hatası", zap.Error(err))
		return err
	}

	return nil
}

func (h *PostgresHandler) GetActor(ctx context.Context, id string) (*domain.Actor, error) {
	ctx, span := h.Tracer.Start(ctx, "GetActor")
	defer span.End()

	var actor domain.Actor
	err := h.Db.WithContext(ctx).Table("actor").Where("actor_id = ?", id).First(&actor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Info("Bu id'li kullanıcı bulunmamaktadır", zap.String("id", id))
		return nil, errors.New("bu id'li kullanıcı bulunmamaktadır")
	}
	if err != nil {
		zap.L().Error("Sorgu çalıştırılırken hata oluştu", zap.Error(err))
		return nil, errors.New("sorgu çalıştırılırken hata oluştu")
	}

	zap.L().Info("Veriler getirildi", zap.String("id", id))
	return &actor, nil
}

func (h *PostgresHandler) UpdateActor(ctx context.Context, id, firstname, lastname string) error {
	ctx, span := h.Tracer.Start(ctx, "UpdateActor")
	defer span.End()

	tx := h.Db.Table("actor").WithContext(ctx).Begin()
	if tx.Error != nil {
		zap.L().Error("transaction başlatılamadı", zap.Error(tx.Error))
		return errors.New("transaction başlatılamadı")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var actor domain.Actor
	if err := tx.Where("actor_id = ?", id).First(&actor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("bu id'ye ait kullanıcı bulunmamaktadır")
		}
		tx.Rollback()
		zap.L().Error("actor sorgusu hatası", zap.Error(err))
		return errors.New("actor sorgusu hatası")
	}

	if actor.FirstName == firstname && actor.LastName == lastname {
		tx.Rollback()
		return errors.New("bu bilgilere ait kullanıcı zaten mevcut")
	}

	actor.FirstName = firstname
	actor.LastName = lastname

	if err := tx.Where("actor_id = ?", actor.ActorID).Updates(&actor).Error; err != nil {
		tx.Rollback()
		zap.L().Error("güncelleme yapılamadı", zap.Error(err))
		return errors.New("güncelleme yapılamadı")
	}

	if err := tx.Commit().Error; err != nil {
		zap.L().Error("transaction commit hatası", zap.Error(err))
		return err
	}

	zap.L().Info("actor güncellendi", zap.String("id", id))
	return nil
}
