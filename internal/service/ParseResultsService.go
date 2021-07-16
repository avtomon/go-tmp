package service

import (
	"github.com/jmoiron/sqlx"
	"parser/internal/domain"
	"parser/internal/repository"
)

func Save(db *sqlx.DB, result *domain.ParseResult) error {
	return repository.Save(db, result.SiteId, &result.PagesParseResults)
}
