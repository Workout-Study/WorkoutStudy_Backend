package service

import (
	"database/sql"
	"workoutstudy_chatting/persistence"
)

func CheckFitGroupExistence(db *sql.DB, fitGroupID int) (bool, error) {
	fitGroup, err := persistence.GetFitGroupByID(db, fitGroupID)
	if err != nil {
		return false, err
	}
	return fitGroup != nil, nil
}
