package middleware

import (
	"worker-nicepay/infrastructure/database/repositories"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Middlewares struct {
	log *zap.Logger
	// HttpResponse *models.HttpResponse // This model likely doesn't exist yet, I might need to create it or remove it. The user code had it.
	// user code had `http *models.HttpResponse` in constructor but struct had `HttpResponse`.
	// I'll keep it but comment it out if I don't see it defined.
	// Actually, the user might have expected me to fix it.
	// Let's assume there is a response utility. I'll mock it or check if it exists.
	// For now I will remove HttpResponse dependency effectively or just keep it as interface if possible.
	// But `models` imported above is from database. `HttpResponse` usually is in `application/dto` or `pkg/utils`.
	// I will remove it for now and implement error response directly in middleware or use a simple struct.
	// user code: `h.HttpResponse.ErrorResponseV2`
	repo *repositories.MasterDataRepositoryYugabyteDB // Using simpler dependency since Repositories struct is missing
	DB   *gorm.DB
}

func NewMiddlewares(log *zap.Logger, repo *repositories.MasterDataRepositoryYugabyteDB, db *gorm.DB) *Middlewares {
	return &Middlewares{
		log:  log,
		repo: repo,
		DB:   db,
	}
}
