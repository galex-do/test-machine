package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/models"
)

// testRunsAPIHandler handles API requests for test runs collection
func (h *Handler) testRunsAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                h.getAllTestRuns(w, r)
        case "POST":
                h.createTestRun(w, r)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// testRunAPIHandler handles API requests for individual test runs
func (h *Handler) testRunAPIHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/api/test-runs/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                h.writeJSONError(w, "Invalid test run ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                h.getTestRun(w, r, id)
        case "PUT":
                h.updateTestRun(w, r, id)
        case "DELETE":
                h.deleteTestRun(w, r, id)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getAllTestRuns(w http.ResponseWriter, r *http.Request) {
        testRuns, err := h.testRunService.GetAllTestRuns()
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        h.writeJSONResponse(w, testRuns)
}

func (h *Handler) getTestRun(w http.ResponseWriter, r *http.Request, id int) {
        testRun, err := h.testRunService.GetTestRunByID(id)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        if testRun == nil {
                h.writeJSONError(w, "Test run not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, testRun)
}

func (h *Handler) createTestRun(w http.ResponseWriter, r *http.Request) {
        var req models.CreateTestRunRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        testRun, err := h.testRunService.CreateTestRun(req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testRun)
}

func (h *Handler) updateTestRun(w http.ResponseWriter, r *http.Request, id int) {
        var req models.UpdateTestRunRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        testRun, err := h.testRunService.UpdateTestRun(id, req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        if testRun == nil {
                h.writeJSONError(w, "Test run not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, testRun)
}

func (h *Handler) deleteTestRun(w http.ResponseWriter, r *http.Request, id int) {
        err := h.testRunService.DeleteTestRun(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) startTestRun(w http.ResponseWriter, r *http.Request, id int) {
        testRun, err := h.testRunService.StartTestRun(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testRun)
}

func (h *Handler) pauseTestRun(w http.ResponseWriter, r *http.Request, id int) {
        testRun, err := h.testRunService.PauseTestRun(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testRun)
}

func (h *Handler) finishTestRun(w http.ResponseWriter, r *http.Request, id int) {
        testRun, err := h.testRunService.FinishTestRun(id)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testRun)
}