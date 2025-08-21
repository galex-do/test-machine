package service

import (
	"errors"

	"github.com/galex-do/test-machine/internal/models"
	"github.com/galex-do/test-machine/internal/repository"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	repo *repository.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// GetAll returns all projects
func (s *ProjectService) GetAll() ([]models.Project, error) {
	return s.repo.GetAll()
}

// GetByID returns a project by ID
func (s *ProjectService) GetByID(id int) (*models.Project, error) {
	return s.repo.GetByID(id)
}

// Create creates a new project
func (s *ProjectService) Create(req *models.CreateProjectRequest) (*models.Project, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	return s.repo.Create(req)
}

// Update updates an existing project
func (s *ProjectService) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	return s.repo.Update(id, req)
}

// Delete deletes a project
func (s *ProjectService) Delete(id int) error {
	return s.repo.Delete(id)
}