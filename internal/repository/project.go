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

// GetAll returns all projects with test suite counts
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
        rows, err := r.db.Query(`
                SELECT p.id, p.name, p.description, p.git_project, p.created_at, p.updated_at, 
                       COALESCE(COUNT(ts.id), 0) as test_suites_count
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                GROUP BY p.id, p.name, p.description, p.git_project, p.created_at, p.updated_at
                ORDER BY p.created_at DESC
        `)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var projects []models.Project
        for rows.Next() {
                var p models.Project
                err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.GitProject, &p.CreatedAt, &p.UpdatedAt, &p.TestSuitesCount)
                if err != nil {
                        return nil, err
                }
                projects = append(projects, p)
        }

        return projects, nil
}

// GetByID returns a project by ID with test suite count
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
        var project models.Project
        err := r.db.QueryRow(`
                SELECT p.id, p.name, p.description, p.git_project, p.created_at, p.updated_at,
                       COALESCE(COUNT(ts.id), 0) as test_suites_count
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                WHERE p.id = $1
                GROUP BY p.id, p.name, p.description, p.git_project, p.created_at, p.updated_at
        `, id).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.CreatedAt, &project.UpdatedAt, &project.TestSuitesCount)

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
                "INSERT INTO projects (name, description, git_project) VALUES ($1, $2, $3) RETURNING id, name, description, git_project, created_at, updated_at",
                req.Name, req.Description, req.GitProject,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.CreatedAt, &project.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
        var project models.Project
        err := r.db.QueryRow(
                "UPDATE projects SET name = $1, description = $2, git_project = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING id, name, description, git_project, created_at, updated_at",
                req.Name, req.Description, req.GitProject, id,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.CreatedAt, &project.UpdatedAt)

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