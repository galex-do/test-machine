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

// CreateOrUpdate creates a new repository or updates existing one with branches and tags
func (r *RepositoryRepository) CreateOrUpdate(repo *models.Repository, branches []models.Branch, tags []models.Tag) (*models.Repository, error) {
        tx, err := r.db.Begin()
        if err != nil {
                return nil, fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        // Check if repository exists for this project
        var existingRepo models.Repository
        err = tx.QueryRow(`
                SELECT id, project_id, remote_url, default_branch, synced_at, created_at, updated_at 
                FROM repositories WHERE project_id = $1
        `, repo.ProjectID).Scan(
                &existingRepo.ID, &existingRepo.ProjectID, &existingRepo.RemoteURL,
                &existingRepo.DefaultBranch, &existingRepo.SyncedAt,
                &existingRepo.CreatedAt, &existingRepo.UpdatedAt,
        )

        var repositoryID int
        if err == sql.ErrNoRows {
                // Create new repository
                err = tx.QueryRow(`
                        INSERT INTO repositories (project_id, remote_url, default_branch, synced_at)
                        VALUES ($1, $2, $3, $4)
                        RETURNING id, project_id, remote_url, default_branch, synced_at, created_at, updated_at
                `, repo.ProjectID, repo.RemoteURL, repo.DefaultBranch, repo.SyncedAt).Scan(
                        &repositoryID, &repo.ProjectID, &repo.RemoteURL,
                        &repo.DefaultBranch, &repo.SyncedAt,
                        &repo.CreatedAt, &repo.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to create repository: %w", err)
                }
                repo.ID = repositoryID
        } else if err != nil {
                return nil, fmt.Errorf("failed to check existing repository: %w", err)
        } else {
                // Update existing repository
                repositoryID = existingRepo.ID
                err = tx.QueryRow(`
                        UPDATE repositories 
                        SET remote_url = $1, default_branch = $2, synced_at = $3, updated_at = CURRENT_TIMESTAMP
                        WHERE id = $4
                        RETURNING id, project_id, remote_url, default_branch, synced_at, created_at, updated_at
                `, repo.RemoteURL, repo.DefaultBranch, repo.SyncedAt, repositoryID).Scan(
                        &repo.ID, &repo.ProjectID, &repo.RemoteURL,
                        &repo.DefaultBranch, &repo.SyncedAt,
                        &repo.CreatedAt, &repo.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to update repository: %w", err)
                }
        }

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

// GetByProjectID returns repository with branches and tags for a project
func (r *RepositoryRepository) GetByProjectID(projectID int) (*models.Repository, error) {
        var repo models.Repository
        err := r.db.QueryRow(`
                SELECT id, project_id, remote_url, default_branch, synced_at, created_at, updated_at
                FROM repositories WHERE project_id = $1
        `, projectID).Scan(
                &repo.ID, &repo.ProjectID, &repo.RemoteURL,
                &repo.DefaultBranch, &repo.SyncedAt,
                &repo.CreatedAt, &repo.UpdatedAt,
        )
        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, fmt.Errorf("failed to get repository: %w", err)
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

        return &repo, nil
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