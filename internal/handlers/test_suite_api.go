package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/models"
)

// testSuitesAPIHandler handles API requests for test suites collection
func (h *Handler) testSuitesAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                h.getAllTestSuites(w, r)
        case "POST":
                h.createTestSuite(w, r)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// testSuiteAPIHandler handles API requests for individual test suites
func (h *Handler) testSuiteAPIHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/api/test-suites/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                h.writeJSONError(w, "Invalid test suite ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                h.getTestSuite(w, r, id)
        case "PUT":
                h.updateTestSuite(w, r, id)
        case "DELETE":
                h.deleteTestSuite(w, r, id)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getAllTestSuites(w http.ResponseWriter, r *http.Request) {
        var projectID *int
        projectIDStr := r.URL.Query().Get("project_id")
        if projectIDStr != "" {
                id, err := strconv.Atoi(projectIDStr)
                if err != nil {
                        h.writeJSONError(w, "Invalid project_id", http.StatusBadRequest)
                        return
                }
                projectID = &id
        }

        testSuites, err := h.testSuiteService.GetAll(projectID)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        h.writeJSONResponse(w, testSuites)
}

func (h *Handler) getTestSuite(w http.ResponseWriter, r *http.Request, id int) {
        testSuite, err := h.testSuiteService.GetByID(id)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        if testSuite == nil {
                h.writeJSONError(w, "Test suite not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, testSuite)
}

func (h *Handler) createTestSuite(w http.ResponseWriter, r *http.Request) {
        var req models.CreateTestSuiteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        testSuite, err := h.testSuiteService.Create(&req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testSuite)
}

func (h *Handler) updateTestSuite(w http.ResponseWriter, r *http.Request, id int) {
        var req models.UpdateTestSuiteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        testSuite, err := h.testSuiteService.Update(id, &req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        if testSuite == nil {
                h.writeJSONError(w, "Test suite not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, testSuite)
}

func (h *Handler) deleteTestSuite(w http.ResponseWriter, r *http.Request, id int) {
        err := h.testSuiteService.Delete(id)
        if err != nil {
                if err.Error() == "sql: no rows in result set" {
                        h.writeJSONError(w, "Test suite not found", http.StatusNotFound)
                        return
                }
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}