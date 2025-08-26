package service

import (
        "fmt"
        "time"

        "github.com/go-git/go-git/v5"
        "github.com/go-git/go-git/v5/plumbing/transport"
        "github.com/go-git/go-git/v5/plumbing/transport/http"
        "github.com/go-git/go-git/v5/plumbing/transport/ssh"
        "github.com/go-git/go-git/v5/storage/memory"
        cryptossh "golang.org/x/crypto/ssh"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/repository"
)

type GitService struct {
        projectRepo    *repository.ProjectRepository
        repositoryRepo *repository.RepositoryRepository
        keyRepo        *repository.KeyRepository
        encryptionSvc  *EncryptionService
}

func NewGitService(projectRepo *repository.ProjectRepository, repositoryRepo *repository.RepositoryRepository, keyRepo *repository.KeyRepository, encryptionSvc *EncryptionService) *GitService {
        return &GitService{
                projectRepo:    projectRepo,
                repositoryRepo: repositoryRepo,
                keyRepo:        keyRepo,
                encryptionSvc:  encryptionSvc,
        }
}

// SyncProjectRepository syncs a project's Git repository and stores branches/tags
func (s *GitService) SyncProjectRepository(projectID int) (*models.SyncResponse, error) {
        // Get project details with repository information
        project, err := s.projectRepo.GetByID(projectID)
        if err != nil {
                return nil, fmt.Errorf("failed to get project: %w", err)
        }
        if project == nil {
                return nil, fmt.Errorf("project not found")
        }
        if project.Repository == nil {
                return nil, fmt.Errorf("project has no Git repository configured")
        }

        return s.SyncRepository(project.Repository.ID)
}

// SyncRepository syncs a repository and stores branches/tags
func (s *GitService) SyncRepository(repositoryID int) (*models.SyncResponse, error) {
        // Get repository details
        repository, err := s.repositoryRepo.GetByID(repositoryID)
        if err != nil {
                return nil, fmt.Errorf("failed to get repository: %w", err)
        }
        if repository == nil {
                return nil, fmt.Errorf("repository not found")
        }

        // Get authentication if key is configured
        var auth transport.AuthMethod
        if repository.KeyID != nil {
                auth, err = s.getAuthMethod(*repository.KeyID, repository.RemoteURL)
                if err != nil {
                        return &models.SyncResponse{
                                Success: false,
                                Message: fmt.Sprintf("Authentication failed: %v", err),
                        }, nil
                }
        }

        // Clone repository to memory to list branches and tags
        repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
                URL:  repository.RemoteURL,
                Auth: auth,
        })
        if err != nil {
                return &models.SyncResponse{
                        Success: false,
                        Message: fmt.Sprintf("Failed to access Git repository: %v", err),
                }, nil
        }

        // Get remote references
        remote, err := repo.Remote("origin")
        if err != nil {
                return &models.SyncResponse{
                        Success: false,
                        Message: fmt.Sprintf("Failed to get remote: %v", err),
                }, nil
        }

        refs, err := remote.List(&git.ListOptions{Auth: auth})
        if err != nil {
                return &models.SyncResponse{
                        Success: false,
                        Message: fmt.Sprintf("Failed to list remote references: %v", err),
                }, nil
        }

        // Parse branches and tags
        var branches []models.Branch
        var tags []models.Tag
        var defaultBranch string

        for _, ref := range refs {
                if ref.Name().IsBranch() {
                        branchName := ref.Name().Short()
                        branches = append(branches, models.Branch{
                                Name:       branchName,
                                CommitHash: stringPtr(ref.Hash().String()),
                                IsDefault:  branchName == "main" || branchName == "master",
                        })
                        if branchName == "main" || (branchName == "master" && defaultBranch == "") {
                                defaultBranch = branchName
                        }
                } else if ref.Name().IsTag() {
                        tags = append(tags, models.Tag{
                                Name:       ref.Name().Short(),
                                CommitHash: stringPtr(ref.Hash().String()),
                        })
                }
        }

        // Update repository with sync data
        repository.DefaultBranch = stringPtr(defaultBranch)
        repository.SyncedAt = timePtr(time.Now())

        updatedRepository, err := s.repositoryRepo.CreateOrUpdateSync(repository, branches, tags)
        if err != nil {
                return &models.SyncResponse{
                        Success: false,
                        Message: fmt.Sprintf("Failed to store repository data: %v", err),
                }, nil
        }

        return &models.SyncResponse{
                Success:     true,
                Message:     "Repository synced successfully",
                Repository:  updatedRepository,
                BranchCount: len(branches),
                TagCount:    len(tags),
        }, nil
}

// getAuthMethod creates authentication method based on key type
func (s *GitService) getAuthMethod(keyID int, repoURL string) (transport.AuthMethod, error) {
        key, err := s.keyRepo.GetByID(keyID)
        if err != nil {
                return nil, fmt.Errorf("failed to get key: %w", err)
        }
        if key == nil {
                return nil, fmt.Errorf("key not found")
        }

        // Decrypt the key data
        secretData, err := s.encryptionSvc.Decrypt(key.EncryptedData)
        if err != nil {
                return nil, fmt.Errorf("failed to decrypt key: %w", err)
        }

        switch key.KeyType {
        case "RSA":
                // SSH authentication
                sshAuth, err := ssh.NewPublicKeys("git", []byte(secretData), "")
                if err != nil {
                        return nil, fmt.Errorf("failed to create SSH auth: %w", err)
                }
                sshAuth.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()
                return sshAuth, nil

        case "Login":
                // HTTPS authentication
                if key.Username == nil {
                        return nil, fmt.Errorf("username is required for Login type keys")
                }
                return &http.BasicAuth{
                        Username: *key.Username,
                        Password: secretData,
                }, nil

        default:
                return nil, fmt.Errorf("unsupported key type: %s", key.KeyType)
        }
}

// Helper functions
func stringPtr(s string) *string {
        return &s
}

func timePtr(t time.Time) *time.Time {
        return &t
}