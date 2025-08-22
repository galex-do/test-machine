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
                SELECT p.id, p.name, p.description, p.git_project, p.key_id, p.created_at, p.updated_at, 
                       COALESCE(COUNT(ts.id), 0) as test_suites_count,
                       k.id, k.name, k.key_type
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                LEFT JOIN keys k ON p.key_id = k.id
                GROUP BY p.id, p.name, p.description, p.git_project, p.key_id, p.created_at, p.updated_at,
                         k.id, k.name, k.key_type
                ORDER BY p.created_at DESC
        `)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var projects []models.Project
        for rows.Next() {
                var p models.Project
                var keyID, keyName, keyType sql.NullString
                err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.GitProject, &p.KeyID, &p.CreatedAt, &p.UpdatedAt, &p.TestSuitesCount, &keyID, &keyName, &keyType)
                if err != nil {
                        return nil, err
                }
                
                // Set key information if available
                if keyID.Valid {
                        var kid int
                        if err := keyID.Scan(&kid); err == nil {
                                p.Key = &models.Key{
                                        ID:      kid,
                                        Name:    keyName.String,
                                        KeyType: keyType.String,
                                }
                        }
                }
                
                projects = append(projects, p)
        }

        return projects, nil
}

// GetByID returns a project by ID with test suite count
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
        var project models.Project
        var keyID, keyName, keyType sql.NullString
        err := r.db.QueryRow(`
                SELECT p.id, p.name, p.description, p.git_project, p.key_id, p.created_at, p.updated_at,
                       COALESCE(COUNT(ts.id), 0) as test_suites_count,
                       k.id, k.name, k.key_type
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                LEFT JOIN keys k ON p.key_id = k.id
                WHERE p.id = $1
                GROUP BY p.id, p.name, p.description, p.git_project, p.key_id, p.created_at, p.updated_at,
                         k.id, k.name, k.key_type
        `, id).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.KeyID, &project.CreatedAt, &project.UpdatedAt, &project.TestSuitesCount, &keyID, &keyName, &keyType)

        if err == nil && keyID.Valid {
                var kid int
                if err := keyID.Scan(&kid); err == nil {
                        project.Key = &models.Key{
                                ID:      kid,
                                Name:    keyName.String,
                                KeyType: keyType.String,
                        }
                }
        }

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
                "INSERT INTO projects (name, description, git_project, key_id) VALUES ($1, $2, $3, $4) RETURNING id, name, description, git_project, key_id, created_at, updated_at",
                req.Name, req.Description, req.GitProject, req.KeyID,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.KeyID, &project.CreatedAt, &project.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
        var project models.Project
        err := r.db.QueryRow(
                "UPDATE projects SET name = $1, description = $2, git_project = $3, key_id = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 RETURNING id, name, description, git_project, key_id, created_at, updated_at",
                req.Name, req.Description, req.GitProject, req.KeyID, id,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.GitProject, &project.KeyID, &project.CreatedAt, &project.UpdatedAt)

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