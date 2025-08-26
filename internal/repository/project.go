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
                SELECT p.id, p.name, p.description, p.repository_id, p.created_at, p.updated_at, 
                       COALESCE(COUNT(ts.id), 0) as test_suites_count,
                       r.id, r.name, r.remote_url, k.id, k.name, k.key_type
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                LEFT JOIN repositories r ON p.repository_id = r.id
                LEFT JOIN keys k ON r.key_id = k.id
                GROUP BY p.id, p.name, p.description, p.repository_id, p.created_at, p.updated_at,
                         r.id, r.name, r.remote_url, k.id, k.name, k.key_type
                ORDER BY p.created_at DESC
        `)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var projects []models.Project
        for rows.Next() {
                var p models.Project
                var repoID, keyID sql.NullInt64
                var repoName, repoURL, keyName, keyType sql.NullString
                err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.RepositoryID, &p.CreatedAt, &p.UpdatedAt, &p.TestSuitesCount, &repoID, &repoName, &repoURL, &keyID, &keyName, &keyType)
                if err != nil {
                        return nil, err
                }
                
                // Set repository information if available
                if repoID.Valid {
                        p.Repository = &models.Repository{
                                ID:        int(repoID.Int64),
                                Name:      repoName.String,
                                RemoteURL: repoURL.String,
                        }
                        
                        // Set key information if available
                        if keyID.Valid {
                                p.Repository.Key = &models.Key{
                                        ID:      int(keyID.Int64),
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
        var repoID, keyID sql.NullInt64
        var repoName, repoURL, keyName, keyType sql.NullString
        err := r.db.QueryRow(`
                SELECT p.id, p.name, p.description, p.repository_id, p.created_at, p.updated_at,
                       COALESCE(COUNT(ts.id), 0) as test_suites_count,
                       r.id, r.name, r.remote_url, k.id, k.name, k.key_type
                FROM projects p
                LEFT JOIN test_suites ts ON p.id = ts.project_id
                LEFT JOIN repositories r ON p.repository_id = r.id
                LEFT JOIN keys k ON r.key_id = k.id
                WHERE p.id = $1
                GROUP BY p.id, p.name, p.description, p.repository_id, p.created_at, p.updated_at,
                         r.id, r.name, r.remote_url, k.id, k.name, k.key_type
        `, id).Scan(&project.ID, &project.Name, &project.Description, &project.RepositoryID, &project.CreatedAt, &project.UpdatedAt, &project.TestSuitesCount, &repoID, &repoName, &repoURL, &keyID, &keyName, &keyType)

        if err == nil && repoID.Valid {
                project.Repository = &models.Repository{
                        ID:        int(repoID.Int64),
                        Name:      repoName.String,
                        RemoteURL: repoURL.String,
                }
                
                // Set key information if available
                if keyID.Valid {
                        project.Repository.Key = &models.Key{
                                ID:      int(keyID.Int64),
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

// CountProjectsByRepositoryID returns the count of projects using a specific repository
func (r *ProjectRepository) CountProjectsByRepositoryID(repositoryID int) (int, error) {
        var count int
        err := r.db.QueryRow("SELECT COUNT(*) FROM projects WHERE repository_id = $1", repositoryID).Scan(&count)
        return count, err
}

// Create creates a new project
func (r *ProjectRepository) Create(req *models.CreateProjectRequest) (*models.Project, error) {
        var project models.Project
        err := r.db.QueryRow(
                "INSERT INTO projects (name, description, repository_id) VALUES ($1, $2, $3) RETURNING id, name, description, repository_id, created_at, updated_at",
                req.Name, req.Description, req.RepositoryID,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.RepositoryID, &project.CreatedAt, &project.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
        var project models.Project
        err := r.db.QueryRow(
                "UPDATE projects SET name = $1, description = $2, repository_id = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING id, name, description, repository_id, created_at, updated_at",
                req.Name, req.Description, req.RepositoryID, id,
        ).Scan(&project.ID, &project.Name, &project.Description, &project.RepositoryID, &project.CreatedAt, &project.UpdatedAt)

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