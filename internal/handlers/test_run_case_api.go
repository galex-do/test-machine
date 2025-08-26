package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/galex-do/test-machine/internal/models"
)

// updateTestRunCase handles PUT /api/test-runs/{runId}/cases/{caseId}
func (h *Handler) updateTestRunCase(w http.ResponseWriter, r *http.Request) {
	runIdStr := r.PathValue("runId")
	caseIdStr := r.PathValue("caseId")
	
	runId, err := strconv.Atoi(runIdStr)
	if err != nil {
		h.writeJSONError(w, "Invalid test run ID", http.StatusBadRequest)
		return
	}

	caseId, err := strconv.Atoi(caseIdStr)
	if err != nil {
		h.writeJSONError(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTestRunCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	testRunCase, err := h.testRunService.UpdateTestRunCase(runId, caseId, req)
	if err != nil {
		h.writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSONResponse(w, testRunCase)
}