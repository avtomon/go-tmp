package repository

import (
	//"errors"
	"github.com/jmoiron/sqlx"
	"parser/internal/domain"
)

func GetSiteConfigs(db *sqlx.DB) (*[]domain.SiteConfig, error) {
	var sites []domain.SiteConfig

	// PRIVATE

	return &sites, nil
}
