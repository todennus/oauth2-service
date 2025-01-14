package gorm

import (
	"context"

	"github.com/todennus/oauth2-service/domain"
	"github.com/todennus/oauth2-service/infras/database/model"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/xcontext"
	"github.com/xybor-x/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OAuth2ConsentRepository struct {
	db *gorm.DB
}

func NewOAuth2ConsentRepository(db *gorm.DB) *OAuth2ConsentRepository {
	return &OAuth2ConsentRepository{db: db}
}

func (repo *OAuth2ConsentRepository) Upsert(ctx context.Context, consent *domain.OAuth2Consent) error {
	model := model.NewOAuth2Consent(consent)
	return errordef.ConvertGormError(
		xcontext.DB(ctx, repo.db).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "client_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"scope", "expires_at", "updated_at"}),
		}).Create(&model).Error,
	)
}

func (repo *OAuth2ConsentRepository) Get(ctx context.Context, userID, clientID snowflake.ID) (*domain.OAuth2Consent, error) {
	model := model.OAuth2ConsentModel{}
	if err := xcontext.DB(ctx, repo.db).Take(&model, "user_id=? AND client_id=?", userID, clientID).Error; err != nil {
		return nil, errordef.ConvertGormError(err)
	}

	return model.To(), nil
}
