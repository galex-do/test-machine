import axios from 'axios'

// Create axios instance with base configuration
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const message = error.response?.data?.message || error.message || 'An error occurred'
    return Promise.reject(new Error(message))
  }
)

// API service methods
export const api = {
  // Projects
  getProjects: () => apiClient.get('/projects'),
  getProject: (id) => apiClient.get(`/projects/${id}`),
  createProject: (data) => apiClient.post('/projects', data),
  updateProject: (id, data) => apiClient.put(`/projects/${id}`, data),
  deleteProject: (id) => apiClient.delete(`/projects/${id}`),

  // Test Suites
  getTestSuites: (projectId) => apiClient.get(`/test-suites${projectId ? `?project_id=${projectId}` : ''}`),
  getTestSuite: (id) => apiClient.get(`/test-suites/${id}`),
  createTestSuite: (data) => apiClient.post('/test-suites', data),
  updateTestSuite: (id, data) => apiClient.put(`/test-suites/${id}`, data),
  deleteTestSuite: (id) => apiClient.delete(`/test-suites/${id}`),

  // Test Cases
  getTestCases: (testSuiteId) => apiClient.get(`/test-cases${testSuiteId ? `?test_suite_id=${testSuiteId}` : ''}`),
  getTestCase: (id) => apiClient.get(`/test-cases/${id}`),
  createTestCase: (data) => apiClient.post('/test-cases', data),
  updateTestCase: (id, data) => apiClient.put(`/test-cases/${id}`, data),
  deleteTestCase: (id) => apiClient.delete(`/test-cases/${id}`),
  searchTestCases: (query) => apiClient.get(`/test-cases/search?q=${encodeURIComponent(query)}`),

  // Test Steps
  getTestSteps: (testCaseId) => apiClient.get(`/test-cases/${testCaseId}/steps`),
  createTestStep: (testCaseId, data) => apiClient.post(`/test-cases/${testCaseId}/steps`, data),
  updateTestStep: (id, data) => apiClient.put(`/test-steps/${id}`, data),
  deleteTestStep: (id) => apiClient.delete(`/test-steps/${id}`),

  // Test Runs
  getTestRuns: () => apiClient.get('/test-runs'),
  getTestRun: (id) => apiClient.get(`/test-runs/${id}`),
  createTestRun: (data) => apiClient.post('/test-runs', data),
  updateTestRun: (id, data) => apiClient.put(`/test-runs/${id}`, data),
  deleteTestRun: (id) => apiClient.delete(`/test-runs/${id}`),
  updateTestRunCase: (runId, caseId, data) => apiClient.put(`/test-runs/${runId}/cases/${caseId}`, data),

  // Helper methods for test runs
  getProjectsWithRepositories: () => apiClient.get('/projects'),
  getTestSuitesByProject: (projectId) => apiClient.get(`/test-suites?project_id=${projectId}`),

  // Keys
  getKeys: () => apiClient.get('/keys'),
  getKey: (id) => apiClient.get(`/keys/${id}`),
  getKeyData: (id) => apiClient.get(`/keys/${id}/data`),
  createKey: (data) => apiClient.post('/keys', data),
  updateKey: (id, data) => apiClient.put(`/keys/${id}`, data),
  deleteKey: (id) => apiClient.delete(`/keys/${id}`),

  // Stats and Reports
  getStats: () => apiClient.get('/stats'),
  getReports: () => apiClient.get('/reports'),

  // Repositories
  getRepositories: () => apiClient.get('/repositories'),
  getRepository: (id) => apiClient.get(`/repositories/${id}`),
  getRepositoryDetails: (id) => apiClient.get(`/repositories/${id}/details`),
  createRepository: (data) => apiClient.post('/repositories', data),
  updateRepository: (id, data) => apiClient.put(`/repositories/${id}`, data),
  deleteRepository: (id) => apiClient.delete(`/repositories/${id}`),
  syncRepository: (repositoryId) => apiClient.post(`/repositories/${repositoryId}/sync`),

  // Sync
  syncProject: (projectId) => apiClient.post(`/sync/projects/${projectId}/sync`)
}

export default api