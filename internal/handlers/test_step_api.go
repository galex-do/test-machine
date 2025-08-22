package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/galex-do/test-machine/internal/models"
)

// testStepAPIHandler handles individual test step operations
func (h *Handler) testStepAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Extract test step ID from URL path: /api/test-steps/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/test-steps/")
	
	if path == "" {
		h.writeJSONError(w, "Test step ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		h.writeJSONError(w, "Invalid test step ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		h.updateTestStep(w, r, id)
	case "DELETE":
		h.deleteTestStep(w, r, id)
	default:
		h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) updateTestStep(w http.ResponseWriter, r *http.Request, id int) {
	var req models.UpdateTestStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	testStep, err := h.testCaseService.UpdateTestStep(id, &req)
	if err != nil {
		h.writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if testStep == nil {
		h.writeJSONError(w, "Test step not found", http.StatusNotFound)
		return
	}

	h.writeJSONResponse(w, testStep)
}

func (h *Handler) deleteTestStep(w http.ResponseWriter, r *http.Request, id int) {
	err := h.testCaseService.DeleteTestStep(id)
	if err != nil {
		h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}