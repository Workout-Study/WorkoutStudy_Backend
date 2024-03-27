package persistence

import (
	"database/sql"
	"log"
	"workoutstudy_chatting/model"
)

func GetFitGroupByID(db *sql.DB, id int) (*model.FitGroup, error) {
	query := `SELECT id, fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by FROM fit_group WHERE id = $1`
	fitGroup := &model.FitGroup{}

	err := db.QueryRow(query, id).Scan(&fitGroup.ID, &fitGroup.FitGroupName, &fitGroup.MaxFitMate, &fitGroup.CreatedAt, &fitGroup.CreatedBy, &fitGroup.UpdatedAt, &fitGroup.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows is not necessarily an error
		}
		log.Printf("Error querying fit_group by ID: %v", err)
		return nil, err
	}

	return fitGroup, nil
}
