package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteGroupRepository struct {
	db *sql.DB
}

func NewSQLiteGroupRepository(db *sql.DB) *SQLiteGroupRepository {
	return &SQLiteGroupRepository{db: db}
}

func (r *SQLiteGroupRepository) Create(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (
			name, description, difficulty_level, category
		) VALUES (?, ?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		group.Name,
		group.Description,
		group.DifficultyLevel,
		group.Category,
	).Scan(&group.ID, &group.CreatedAt)

	return err
}

func (r *SQLiteGroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	group := &models.Group{}
	query := `
		SELECT id, name, description, difficulty_level, category, created_at
		FROM groups
		WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&group.DifficultyLevel,
		&group.Category,
		&group.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r *SQLiteGroupRepository) List(ctx context.Context, offset, limit int) ([]*models.Group, error) {
	query := `
		SELECT id, name, description, difficulty_level, category, created_at
		FROM groups
		ORDER BY difficulty_level, name
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.DifficultyLevel,
			&group.Category,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, rows.Err()
}

func (r *SQLiteGroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM groups WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SQLiteGroupRepository) Update(ctx context.Context, group *models.Group) error {
	query := `UPDATE groups 
			  SET name = ?, description = ?, difficulty_level = ?, category = ? 
			  WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query,
		group.Name, group.Description, group.DifficultyLevel, group.Category, group.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
