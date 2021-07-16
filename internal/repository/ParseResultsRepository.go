package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"parser/internal/domain"
)

func Save(db *sqlx.DB, siteId uint16, parseResults *[]domain.PageResponse) error {
	stmt, err := db.Prepare("INSERT INTO parse_result(site_id, url, data) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	for _, parseResult := range *parseResults {
		data, err := json.Marshal(parseResult.Data)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(siteId, parseResult.PageUrl, data)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
