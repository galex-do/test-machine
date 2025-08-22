package repository

import (
	"database/sql"

	"github.com/galex-do/test-machine/internal/models"
)

// KeyRepository handles database operations for keys
type KeyRepository struct {
	db *sql.DB
}

// NewKeyRepository creates a new key repository
func NewKeyRepository(db *sql.DB) *KeyRepository {
	return &KeyRepository{db: db}
}

// GetAll returns all keys (without decrypted data)
func (r *KeyRepository) GetAll() ([]models.Key, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, key_type, username, created_at, updated_at
		FROM keys
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.Key
	for rows.Next() {
		var k models.Key
		err := rows.Scan(&k.ID, &k.Name, &k.Description, &k.KeyType, &k.Username, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	return keys, nil
}

// GetByID returns a key by ID (without decrypted data)
func (r *KeyRepository) GetByID(id int) (*models.Key, error) {
	var key models.Key
	err := r.db.QueryRow(`
		SELECT id, name, description, key_type, username, created_at, updated_at
		FROM keys
		WHERE id = $1
	`, id).Scan(&key.ID, &key.Name, &key.Description, &key.KeyType, &key.Username, &key.CreatedAt, &key.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &key, nil
}

// GetEncryptedData returns the encrypted data for a specific key
func (r *KeyRepository) GetEncryptedData(id int) (string, error) {
	var encryptedData string
	err := r.db.QueryRow("SELECT encrypted_data FROM keys WHERE id = $1", id).Scan(&encryptedData)
	
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

// Create creates a new key with encrypted data
func (r *KeyRepository) Create(req *models.CreateKeyRequest, encryptedData string) (*models.Key, error) {
	var key models.Key
	err := r.db.QueryRow(`
		INSERT INTO keys (name, description, key_type, username, encrypted_data) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, name, description, key_type, username, created_at, updated_at`,
		req.Name, req.Description, req.KeyType, req.Username, encryptedData,
	).Scan(&key.ID, &key.Name, &key.Description, &key.KeyType, &key.Username, &key.CreatedAt, &key.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &key, nil
}

// Update updates an existing key
func (r *KeyRepository) Update(id int, req *models.UpdateKeyRequest, encryptedData *string) (*models.Key, error) {
	var key models.Key
	var err error

	if encryptedData != nil {
		// Update with new encrypted data
		err = r.db.QueryRow(`
			UPDATE keys SET name = $1, description = $2, username = $3, encrypted_data = $4, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $5 
			RETURNING id, name, description, key_type, username, created_at, updated_at`,
			req.Name, req.Description, req.Username, *encryptedData, id,
		).Scan(&key.ID, &key.Name, &key.Description, &key.KeyType, &key.Username, &key.CreatedAt, &key.UpdatedAt)
	} else {
		// Update without changing encrypted data
		err = r.db.QueryRow(`
			UPDATE keys SET name = $1, description = $2, username = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $4 
			RETURNING id, name, description, key_type, username, created_at, updated_at`,
			req.Name, req.Description, req.Username, id,
		).Scan(&key.ID, &key.Name, &key.Description, &key.KeyType, &key.Username, &key.CreatedAt, &key.UpdatedAt)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &key, nil
}

// Delete deletes a key
func (r *KeyRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM keys WHERE id = $1", id)
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