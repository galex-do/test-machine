package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/galex-do/test-machine/internal/models"
	"github.com/galex-do/test-machine/internal/service"
)

// TestRunHandler handles test run HTTP requests
type TestRunHandler struct {
	service *service.TestRunService
}

// NewTestRunHandler creates a new test run handler
func NewTestRunHandler(service *service.TestRunService) *TestRunHandler {
	return &TestRunHandler{service: service}
}

// GetAll handles GET /api/test-runs
func (h *TestRunHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	testRuns, err := h.service.GetAllTestRuns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testRuns)
}

// GetByID handles GET /api/test-runs/{id}
func (h *TestRunHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test run ID", http.StatusBadRequest)
		return
	}

	testRun, err := h.service.GetTestRunByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if testRun == nil {
		http.Error(w, "Test run not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testRun)
}

// Create handles POST /api/test-runs
func (h *TestRunHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTestRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	testRun, err := h.service.CreateTestRun(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testRun)
}

// Update handles PUT /api/test-runs/{id}
func (h *TestRunHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test run ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTestRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	testRun, err := h.service.UpdateTestRun(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testRun)
}

// Delete handles DELETE /api/test-runs/{id}
func (h *TestRunHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test run ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTestRun(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTestCase handles PUT /api/test-runs/{runId}/cases/{caseId}
func (h *TestRunHandler) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	runIdStr := r.PathValue("runId")
	caseIdStr := r.PathValue("caseId")
	
	runId, err := strconv.Atoi(runIdStr)
	if err != nil {
		http.Error(w, "Invalid test run ID", http.StatusBadRequest)
		return
	}

	caseId, err := strconv.Atoi(caseIdStr)
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTestRunCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	testRunCase, err := h.service.UpdateTestRunCase(runId, caseId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testRunCase)
}