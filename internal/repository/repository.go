package repository

import (
        "database/sql"
        "fmt"

        "github.com/galex-do/test-machine/internal/models"
)

type RepositoryRepository struct {
        db *sql.DB
}

func NewRepositoryRepository(db *sql.DB) *RepositoryRepository {
        return &RepositoryRepository{db: db}
}

// Create creates a new repository
func (r *RepositoryRepository) Create(req *models.CreateRepositoryRequest) (*models.Repository, error) {
        var repo models.Repository
        err := r.db.QueryRow(`
                INSERT INTO repositories (name, description, remote_url, key_id)
                VALUES ($1, $2, $3, $4)
                RETURNING id, name, description, remote_url, key_id, default_branch, synced_at, created_at, updated_at
        `, req.Name, req.Description, req.RemoteURL, req.KeyID).Scan(
                &repo.ID, &repo.Name, &repo.Description, &repo.RemoteURL, &repo.KeyID,
                &repo.DefaultBranch, &repo.SyncedAt, &repo.CreatedAt, &repo.UpdatedAt,
        )

        if err != nil {
                return nil, fmt.Errorf("failed to create repository: %w", err)
        }

        return &repo, nil
}

// Update updates an existing repository (name, description, key only - remote_url is immutable)
func (r *RepositoryRepository) Update(id int, req *models.UpdateRepositoryRequest) (*models.Repository, error) {
        var repo models.Repository
        err := r.db.QueryRow(`
                UPDATE repositories 
                SET name = $1, description = $2, key_id = $3, updated_at = CURRENT_TIMESTAMP
                WHERE id = $4
                RETURNING id, name, description, remote_url, key_id, default_branch, synced_at, created_at, updated_at
        `, req.Name, req.Description, req.KeyID, id).Scan(
                &repo.ID, &repo.Name, &repo.Description, &repo.RemoteURL, &repo.KeyID,
                &repo.DefaultBranch, &repo.SyncedAt, &repo.CreatedAt, &repo.UpdatedAt,
        )

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, fmt.Errorf("failed to update repository: %w", err)
        }

        return &repo, nil
}

// GetAll returns all repositories with key information
func (r *RepositoryRepository) GetAll() ([]models.Repository, error) {
        rows, err := r.db.Query(`
                SELECT r.id, r.name, r.description, r.remote_url, r.key_id, r.default_branch, r.synced_at, r.created_at, r.updated_at,
                       k.id, k.name, k.key_type
                FROM repositories r
                LEFT JOIN keys k ON r.key_id = k.id
                ORDER BY r.created_at DESC
        `)
        if err != nil {
                return nil, fmt.Errorf("failed to query repositories: %w", err)
        }
        defer rows.Close()

        var repositories []models.Repository
        for rows.Next() {
                var repo models.Repository
                var keyID, keyName, keyType sql.NullString
                err := rows.Scan(
                        &repo.ID, &repo.Name, &repo.Description, &repo.RemoteURL, &repo.KeyID,
                        &repo.DefaultBranch, &repo.SyncedAt, &repo.CreatedAt, &repo.UpdatedAt,
                        &keyID, &keyName, &keyType,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan repository: %w", err)
                }

                // Set key information if available
                if keyID.Valid {
                        var kid int
                        if err := keyID.Scan(&kid); err == nil {
                                repo.Key = &models.Key{
                                        ID:      kid,
                                        Name:    keyName.String,
                                        KeyType: keyType.String,
                                }
                        }
                }

                repositories = append(repositories, repo)
        }

        return repositories, nil
}

// GetByID returns a repository by ID with key information
func (r *RepositoryRepository) GetByID(id int) (*models.Repository, error) {
        var repo models.Repository
        var keyID, keyName, keyType sql.NullString
        err := r.db.QueryRow(`
                SELECT r.id, r.name, r.description, r.remote_url, r.key_id, r.default_branch, r.synced_at, r.created_at, r.updated_at,
                       k.id, k.name, k.key_type
                FROM repositories r
                LEFT JOIN keys k ON r.key_id = k.id
                WHERE r.id = $1
        `, id).Scan(
                &repo.ID, &repo.Name, &repo.Description, &repo.RemoteURL, &repo.KeyID,
                &repo.DefaultBranch, &repo.SyncedAt, &repo.CreatedAt, &repo.UpdatedAt,
                &keyID, &keyName, &keyType,
        )

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, fmt.Errorf("failed to get repository: %w", err)
        }

        // Set key information if available
        if keyID.Valid {
                var kid int
                if err := keyID.Scan(&kid); err == nil {
                        repo.Key = &models.Key{
                                ID:      kid,
                                Name:    keyName.String,
                                KeyType: keyType.String,
                        }
                }
        }

        return &repo, nil
}

// Delete deletes a repository and all associated branches/tags
func (r *RepositoryRepository) Delete(id int) error {
        // Foreign key constraints will handle cascading deletes for branches/tags
        result, err := r.db.Exec("DELETE FROM repositories WHERE id = $1", id)
        if err != nil {
                return fmt.Errorf("failed to delete repository: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get rows affected: %w", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("repository not found")
        }

        return nil
}

// CreateOrUpdateSync creates or updates a repository with sync data (branches and tags)
func (r *RepositoryRepository) CreateOrUpdateSync(repo *models.Repository, branches []models.Branch, tags []models.Tag) (*models.Repository, error) {
        tx, err := r.db.Begin()
        if err != nil {
                return nil, fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        // Update repository sync information
        err = tx.QueryRow(`
                UPDATE repositories 
                SET default_branch = $1, synced_at = $2, updated_at = CURRENT_TIMESTAMP
                WHERE id = $3
                RETURNING id, name, description, remote_url, key_id, default_branch, synced_at, created_at, updated_at
        `, repo.DefaultBranch, repo.SyncedAt, repo.ID).Scan(
                &repo.ID, &repo.Name, &repo.Description, &repo.RemoteURL, &repo.KeyID,
                &repo.DefaultBranch, &repo.SyncedAt, &repo.CreatedAt, &repo.UpdatedAt,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to update repository sync info: %w", err)
        }

        repositoryID := repo.ID

        // Clear existing branches and tags
        _, err = tx.Exec("DELETE FROM branches WHERE repository_id = $1", repositoryID)
        if err != nil {
                return nil, fmt.Errorf("failed to delete existing branches: %w", err)
        }

        _, err = tx.Exec("DELETE FROM tags WHERE repository_id = $1", repositoryID)
        if err != nil {
                return nil, fmt.Errorf("failed to delete existing tags: %w", err)
        }

        // Insert new branches
        for i := range branches {
                branches[i].RepositoryID = repositoryID
                err = tx.QueryRow(`
                        INSERT INTO branches (repository_id, name, commit_hash, is_default)
                        VALUES ($1, $2, $3, $4)
                        RETURNING id, created_at, updated_at
                `, branches[i].RepositoryID, branches[i].Name, branches[i].CommitHash, branches[i].IsDefault).Scan(
                        &branches[i].ID, &branches[i].CreatedAt, &branches[i].UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to insert branch %s: %w", branches[i].Name, err)
                }
        }

        // Insert new tags
        for i := range tags {
                tags[i].RepositoryID = repositoryID
                err = tx.QueryRow(`
                        INSERT INTO tags (repository_id, name, commit_hash)
                        VALUES ($1, $2, $3)
                        RETURNING id, created_at, updated_at
                `, tags[i].RepositoryID, tags[i].Name, tags[i].CommitHash).Scan(
                        &tags[i].ID, &tags[i].CreatedAt, &tags[i].UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to insert tag %s: %w", tags[i].Name, err)
                }
        }

        // Commit transaction
        err = tx.Commit()
        if err != nil {
                return nil, fmt.Errorf("failed to commit transaction: %w", err)
        }

        // Set branches and tags on repository
        repo.Branches = branches
        repo.Tags = tags

        return repo, nil
}

// GetWithBranchesAndTags returns repository with branches and tags by repository ID
func (r *RepositoryRepository) GetWithBranchesAndTags(repositoryID int) (*models.Repository, error) {
        repo, err := r.GetByID(repositoryID)
        if err != nil || repo == nil {
                return repo, err
        }

        // Get branches
        branchRows, err := r.db.Query(`
                SELECT id, repository_id, name, commit_hash, is_default, created_at, updated_at
                FROM branches WHERE repository_id = $1 ORDER BY is_default DESC, name
        `, repo.ID)
        if err != nil {
                return nil, fmt.Errorf("failed to get branches: %w", err)
        }
        defer branchRows.Close()

        for branchRows.Next() {
                var branch models.Branch
                err = branchRows.Scan(
                        &branch.ID, &branch.RepositoryID, &branch.Name,
                        &branch.CommitHash, &branch.IsDefault,
                        &branch.CreatedAt, &branch.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan branch: %w", err)
                }
                repo.Branches = append(repo.Branches, branch)
        }

        // Get tags
        tagRows, err := r.db.Query(`
                SELECT id, repository_id, name, commit_hash, created_at, updated_at
                FROM tags WHERE repository_id = $1 ORDER BY name DESC
        `, repo.ID)
        if err != nil {
                return nil, fmt.Errorf("failed to get tags: %w", err)
        }
        defer tagRows.Close()

        for tagRows.Next() {
                var tag models.Tag
                err = tagRows.Scan(
                        &tag.ID, &tag.RepositoryID, &tag.Name,
                        &tag.CommitHash, &tag.CreatedAt, &tag.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan tag: %w", err)
                }
                repo.Tags = append(repo.Tags, tag)
        }

        return repo, nil
}

// GetBranchesByRepositoryID returns all branches for a repository
func (r *RepositoryRepository) GetBranchesByRepositoryID(repositoryID int) ([]models.Branch, error) {
        rows, err := r.db.Query(`
                SELECT id, repository_id, name, commit_hash, is_default, created_at, updated_at
                FROM branches WHERE repository_id = $1 ORDER BY is_default DESC, name
        `, repositoryID)
        if err != nil {
                return nil, fmt.Errorf("failed to get branches: %w", err)
        }
        defer rows.Close()

        var branches []models.Branch
        for rows.Next() {
                var branch models.Branch
                err = rows.Scan(
                        &branch.ID, &branch.RepositoryID, &branch.Name,
                        &branch.CommitHash, &branch.IsDefault,
                        &branch.CreatedAt, &branch.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan branch: %w", err)
                }
                branches = append(branches, branch)
        }

        return branches, nil
}

// GetTagsByRepositoryID returns all tags for a repository
func (r *RepositoryRepository) GetTagsByRepositoryID(repositoryID int) ([]models.Tag, error) {
        rows, err := r.db.Query(`
                SELECT id, repository_id, name, commit_hash, created_at, updated_at
                FROM tags WHERE repository_id = $1 ORDER BY name DESC
        `, repositoryID)
        if err != nil {
                return nil, fmt.Errorf("failed to get tags: %w", err)
        }
        defer rows.Close()

        var tags []models.Tag
        for rows.Next() {
                var tag models.Tag
                err = rows.Scan(
                        &tag.ID, &tag.RepositoryID, &tag.Name,
                        &tag.CommitHash, &tag.CreatedAt, &tag.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan tag: %w", err)
                }
                tags = append(tags, tag)
        }

        return tags, nil
}