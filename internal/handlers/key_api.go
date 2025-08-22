package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/models"
)

// KeyAPIHandler handles API requests for keys
func (h *Handler) keyAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                h.getKeys(w, r)
        case "POST":
                h.createKey(w, r)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// keyByIDAPIHandler handles API requests for individual keys
func (h *Handler) keyByIDAPIHandler(w http.ResponseWriter, r *http.Request) {
        // Extract key ID from URL path: /api/keys/{id}
        path := strings.TrimPrefix(r.URL.Path, "/api/keys/")
        
        if path == "" {
                h.writeJSONError(w, "Key ID is required", http.StatusBadRequest)
                return
        }

        // Handle special endpoint for getting decrypted data
        if strings.HasSuffix(path, "/data") {
                idStr := strings.TrimSuffix(path, "/data")
                id, err := strconv.Atoi(idStr)
                if err != nil {
                        h.writeJSONError(w, "Invalid key ID", http.StatusBadRequest)
                        return
                }
                
                if r.Method == "GET" {
                        h.getKeyData(w, r, id)
                } else {
                        h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
                }
                return
        }

        id, err := strconv.Atoi(path)
        if err != nil {
                h.writeJSONError(w, "Invalid key ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                h.getKey(w, r, id)
        case "PUT":
                h.updateKey(w, r, id)
        case "DELETE":
                h.deleteKey(w, r, id)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getKeys(w http.ResponseWriter, r *http.Request) {
        keys, err := h.keyService.GetAll()
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }
        
        h.writeJSONResponse(w, keys)
}

func (h *Handler) getKey(w http.ResponseWriter, r *http.Request, id int) {
        key, err := h.keyService.GetByID(id)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        if key == nil {
                h.writeJSONError(w, "Key not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, key)
}

func (h *Handler) getKeyData(w http.ResponseWriter, r *http.Request, id int) {
        data, err := h.keyService.GetDecryptedData(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        response := map[string]string{"data": data}
        h.writeJSONResponse(w, response)
}

func (h *Handler) createKey(w http.ResponseWriter, r *http.Request) {
        var req models.CreateKeyRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        key, err := h.keyService.Create(&req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, key)
}

func (h *Handler) updateKey(w http.ResponseWriter, r *http.Request, id int) {
        var req models.UpdateKeyRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        key, err := h.keyService.Update(id, &req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        if key == nil {
                h.writeJSONError(w, "Key not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, key)
}

func (h *Handler) deleteKey(w http.ResponseWriter, r *http.Request, id int) {
        err := h.keyService.Delete(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}