package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/models"
)

// testCasesAPIHandler handles API requests for test cases collection
func (h *Handler) testCasesAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                h.getAllTestCases(w, r)
        case "POST":
                h.createTestCase(w, r)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// testCaseAPIHandler handles API requests for individual test cases
func (h *Handler) testCaseAPIHandler(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimPrefix(r.URL.Path, "/api/test-cases/")
        
        // Handle test steps endpoint: /api/test-cases/{id}/steps
        if strings.Contains(path, "/steps") {
                parts := strings.Split(path, "/")
                if len(parts) >= 1 {
                        id, err := strconv.Atoi(parts[0])
                        if err != nil {
                                h.writeJSONError(w, "Invalid test case ID", http.StatusBadRequest)
                                return
                        }
                        h.handleTestSteps(w, r, id)
                        return
                }
        }
        
        // Handle regular test case requests: /api/test-cases/{id}
        id, err := strconv.Atoi(path)
        if err != nil {
                h.writeJSONError(w, "Invalid test case ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                h.getTestCase(w, r, id)
        case "PUT":
                h.updateTestCase(w, r, id)
        case "DELETE":
                h.deleteTestCase(w, r, id)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getAllTestCases(w http.ResponseWriter, r *http.Request) {
        var testSuiteID *int
        testSuiteIDStr := r.URL.Query().Get("test_suite_id")
        if testSuiteIDStr != "" {
                id, err := strconv.Atoi(testSuiteIDStr)
                if err != nil {
                        h.writeJSONError(w, "Invalid test_suite_id", http.StatusBadRequest)
                        return
                }
                testSuiteID = &id
        }

        testCases, err := h.testCaseService.GetAll(testSuiteID)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        h.writeJSONResponse(w, testCases, http.StatusOK)
}

func (h *Handler) getTestCase(w http.ResponseWriter, r *http.Request, id int) {
        testCase, err := h.testCaseService.GetByID(id)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        if testCase == nil {
                h.writeJSONError(w, "Test case not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, testCase, http.StatusOK)
}

func (h *Handler) createTestCase(w http.ResponseWriter, r *http.Request) {
        var req models.CreateTestCaseRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        testCase, err := h.testCaseService.Create(&req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, testCase, http.StatusCreated)
}

// handleTestSteps handles test steps requests for a specific test case
func (h *Handler) handleTestSteps(w http.ResponseWriter, r *http.Request, testCaseID int) {
        switch r.Method {
        case "GET":
                h.getTestSteps(w, r, testCaseID)
        case "POST":
                h.createTestStep(w, r, testCaseID)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getTestSteps(w http.ResponseWriter, r *http.Request, testCaseID int) {
        testSteps, err := h.testCaseService.GetTestSteps(testCaseID)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }
        
        h.writeJSONResponse(w, testSteps, http.StatusOK)
}

func (h *Handler) createTestStep(w http.ResponseWriter, r *http.Request, testCaseID int) {
        var req models.CreateTestStepRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        req.TestCaseID = testCaseID
        
        testStep, err := h.testCaseService.CreateTestStep(&req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }
        
        h.writeJSONResponse(w, testStep, http.StatusCreated)
}

func (h *Handler) updateTestCase(w http.ResponseWriter, r *http.Request, id int) {
        // TODO: Implement update functionality when service method is available
        h.writeJSONError(w, "Update not implemented", http.StatusNotImplemented)
}

func (h *Handler) deleteTestCase(w http.ResponseWriter, r *http.Request, id int) {
        err := h.testCaseService.Delete(id)
        if err != nil {
                if err.Error() == "sql: no rows in result set" {
                        h.writeJSONError(w, "Test case not found", http.StatusNotFound)
                        return
                }
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}