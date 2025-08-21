package repository

import (
	"database/sql"

	"github.com/galex-do/test-machine/internal/models"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetAll returns all projects
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetByID returns a project by ID
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	var project models.Project
	err := r.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM projects WHERE id = $1",
		id,
	).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Create creates a new project
func (r *ProjectRepository) Create(req *models.CreateProjectRequest) (*models.Project, error) {
	var project models.Project
	err := r.db.QueryRow(
		"INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at",
		req.Name, req.Description,
	).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
	var project models.Project
	err := r.db.QueryRow(
		"UPDATE projects SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, description, created_at, updated_at",
		req.Name, req.Description, id,
	).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}