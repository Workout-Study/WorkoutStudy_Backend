package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"workoutstudy_chatting/model"
)

type FitGroupRepository interface {
	GetFitGroupByID(id int) (*model.FitGroup, error)
	GetFitMatesByFitGroupId(id int) ([]int, error)
	SaveFitGroup(fitGroup *model.FitGroup) (*model.FitGroup, error)
	DeleteFitGroup(fitGroupID int) error
	UpdateFitGroup(fitGroup *model.FitGroup) error
}

type FitGroupRepositoryImpl struct {
	DB *sql.DB
}

// 훈기 tip : 인터페이스 메소드 슬라이스 중 구현안되거 있으면 에러 띄워줌
var _ FitGroupRepository = &FitGroupRepositoryImpl{}

func NewFitGroupRepository(db *sql.DB) FitGroupRepository {
	return &FitGroupRepositoryImpl{DB: db}
}

func (repo *FitGroupRepositoryImpl) GetFitGroupByID(id int) (*model.FitGroup, error) {
	query := `SELECT id, fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by FROM fit_group WHERE id = $1`

	log.Printf("Repository layer: Executing query for FitGroupID: %d", id)
	fitGroup := model.FitGroup{}
	err := repo.DB.QueryRow(query, id).Scan(&fitGroup.ID, &fitGroup.FitGroupName, &fitGroup.MaxFitMate, &fitGroup.CreatedAt, &fitGroup.CreatedBy, &fitGroup.UpdatedAt, &fitGroup.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Repository layer: No fit_group found for ID: %v", id)
			return nil, fmt.Errorf("no fit_group found for ID: %d", id)
		}
		log.Printf("Repository layer: Error querying fit_group by ID: %v", err)
		return nil, err
	}

	return &fitGroup, nil
}

func (repo *FitGroupRepositoryImpl) GetFitMatesByFitGroupId(id int) ([]int, error) {
	// SQL 쿼리 정의
	query := `
        SELECT f_mate.fit_mate_id
        FROM fit_group_mate f_mate
        JOIN fit_mate f ON f_mate.fit_mate_id = f.id
        WHERE f_mate.fit_group_id = $1
    `

	// 쿼리 실행을 위한 슬라이스 준비
	var fitMateIds []int

	// 쿼리 실행
	rows, err := repo.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 결과 처리
	for rows.Next() {
		var fitMateId int
		if err := rows.Scan(&fitMateId); err != nil {
			return nil, err
		}
		fitMateIds = append(fitMateIds, fitMateId)
	}

	// 모든 결과 처리 후 에러 확인
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fitMateIds, nil
}

func (repo *FitGroupRepositoryImpl) SaveFitGroup(fitGroup *model.FitGroup) (*model.FitGroup, error) {
	query := `
		INSERT INTO fit_group (fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := repo.DB.QueryRow(query, fitGroup.FitGroupName, fitGroup.MaxFitMate, fitGroup.CreatedAt, fitGroup.CreatedBy, fitGroup.UpdatedAt, fitGroup.UpdatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}

	fitGroup.ID = id
	return fitGroup, nil
}

func (repo *FitGroupRepositoryImpl) DeleteFitGroup(fitGroupID int) error {
	query := `DELETE FROM fit_group WHERE id = $1`

	_, err := repo.DB.Exec(query, fitGroupID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *FitGroupRepositoryImpl) UpdateFitGroup(fitGroup *model.FitGroup) error {
	query := `
		UPDATE fit_group
		SET fit_group_name = $1, category = $2, cycle = $3, frequency = $4, present_fit_mate_count = $5, max_fit_mate = $6, state = $7, created_at = $8, created_by = $9, updated_at = NOW(), updated_by = $11
		WHERE id = $12
	`
	_, err := repo.DB.Exec(query, fitGroup.FitGroupName, fitGroup.Category, fitGroup.Cycle, fitGroup.Frequency, fitGroup.PresentFitMateCount, fitGroup.MaxFitMate, fitGroup.State, fitGroup.CreatedAt, fitGroup.CreatedBy, fitGroup.UpdatedBy, fitGroup.ID)
	if err != nil {
		return err
	}

	return nil
}
