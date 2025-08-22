package service

import (
	"errors"

	"github.com/galex-do/test-machine/internal/models"
	"github.com/galex-do/test-machine/internal/repository"
)

// KeyService handles business logic for keys
type KeyService struct {
	repo            *repository.KeyRepository
	encryptionService *EncryptionService
}

// NewKeyService creates a new key service
func NewKeyService(repo *repository.KeyRepository, encryptionService *EncryptionService) *KeyService {
	return &KeyService{
		repo:              repo,
		encryptionService: encryptionService,
	}
}

// GetAll returns all keys
func (s *KeyService) GetAll() ([]models.Key, error) {
	return s.repo.GetAll()
}

// GetByID returns a key by ID
func (s *KeyService) GetByID(id int) (*models.Key, error) {
	return s.repo.GetByID(id)
}

// GetDecryptedData returns the decrypted secret data for a key
func (s *KeyService) GetDecryptedData(id int) (string, error) {
	encryptedData, err := s.repo.GetEncryptedData(id)
	if err != nil {
		return "", err
	}
	
	if encryptedData == "" {
		return "", errors.New("key not found")
	}

	return s.encryptionService.Decrypt(encryptedData)
}

// Create creates a new key
func (s *KeyService) Create(req *models.CreateKeyRequest) (*models.Key, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.KeyType == "" {
		return nil, errors.New("key type is required")
	}
	if req.KeyType != "RSA" && req.KeyType != "Login" {
		return nil, errors.New("key type must be 'RSA' or 'Login'")
	}
	if req.SecretData == "" {
		return nil, errors.New("secret data is required")
	}
	if req.KeyType == "Login" && (req.Username == nil || *req.Username == "") {
		return nil, errors.New("username is required for Login type keys")
	}

	// Encrypt the secret data
	encryptedData, err := s.encryptionService.Encrypt(req.SecretData)
	if err != nil {
		return nil, errors.New("failed to encrypt secret data")
	}

	return s.repo.Create(req, encryptedData)
}

// Update updates an existing key
func (s *KeyService) Update(id int, req *models.UpdateKeyRequest) (*models.Key, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	var encryptedData *string
	if req.SecretData != nil && *req.SecretData != "" {
		// Encrypt new secret data
		encrypted, err := s.encryptionService.Encrypt(*req.SecretData)
		if err != nil {
			return nil, errors.New("failed to encrypt secret data")
		}
		encryptedData = &encrypted
	}

	return s.repo.Update(id, req, encryptedData)
}

// Delete deletes a key
func (s *KeyService) Delete(id int) error {
	return s.repo.Delete(id)
}