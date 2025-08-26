<template>
  <div class="test-run-form-container">
    <!-- Header -->
    <div class="d-flex justify-content-between align-items-center mb-4">
      <div>
        <nav aria-label="breadcrumb">
          <ol class="breadcrumb">
            <li class="breadcrumb-item">
              <router-link to="/test-runs" class="text-decoration-none">
                <i class="fas fa-play-circle"></i> Test Runs
              </router-link>
            </li>
            <li class="breadcrumb-item active">{{ isEditing ? 'Edit Test Run' : 'New Test Run' }}</li>
          </ol>
        </nav>
        <h2>{{ isEditing ? 'Edit Test Run' : 'Create New Test Run' }}</h2>
      </div>
      <router-link to="/test-runs" class="btn btn-outline-secondary">
        <i class="fas fa-times"></i> Cancel
      </router-link>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading data...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="alert alert-danger" role="alert">
      <i class="fas fa-exclamation-triangle"></i> {{ error }}
    </div>

    <!-- Form -->
    <div v-if="!loading && !error" class="row">
      <!-- Left Column: Basic Info -->
      <div class="col-md-5">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0"><i class="fas fa-info-circle"></i> Test Run Details</h5>
          </div>
          <div class="card-body">
            <form @submit.prevent="saveTestRun">
              <!-- Name -->
              <div class="mb-3">
                <label for="name" class="form-label">Name</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="form.name"
                  placeholder="Leave empty for auto-generation"
                >
                <div class="form-text">Auto-generated format: &lt;project&gt;-&lt;branch/tag&gt;-&lt;datetime&gt;</div>
              </div>

              <!-- Description -->
              <div class="mb-3">
                <label for="description" class="form-label">Description</label>
                <textarea
                  class="form-control"
                  id="description"
                  rows="3"
                  v-model="form.description"
                  placeholder="Describe the purpose of this test run"
                ></textarea>
              </div>

              <!-- Project Selection -->
              <div class="mb-3">
                <label for="project" class="form-label">Project *</label>
                <select
                  class="form-select"
                  id="project"
                  v-model="form.project_id"
                  @change="onProjectChange"
                  required
                >
                  <option value="">Select a project...</option>
                  <option v-for="project in projects" :key="project.id" :value="project.id">
                    {{ project.name }}
                  </option>
                </select>
              </div>

              <!-- Git Reference -->
              <div v-if="selectedProject && selectedProject.repository" class="mb-3">
                <label class="form-label">Git Reference</label>
                <div class="row">
                  <div class="col-6">
                    <select class="form-select form-select-sm" v-model="gitReferenceType">
                      <option value="branch">Branch</option>
                      <option value="tag">Tag</option>
                    </select>
                  </div>
                  <div class="col-6">
                    <select 
                      class="form-select form-select-sm" 
                      :disabled="!gitReferences.length"
                      v-model="selectedGitReference"
                    >
                      <option value="">{{ gitReferenceType === 'branch' ? 'Select branch...' : 'Select tag...' }}</option>
                      <option v-for="ref in gitReferences" :key="ref" :value="ref">{{ ref }}</option>
                    </select>
                  </div>
                </div>
                <div v-if="!gitReferences.length && selectedProject.repository" class="form-text text-muted">
                  Loading {{ gitReferenceType }}s...
                </div>
              </div>

              <!-- Created By -->
              <div class="mb-3">
                <label for="createdBy" class="form-label">Created By *</label>
                <input
                  type="text"
                  class="form-control"
                  id="createdBy"
                  v-model="form.created_by"
                  placeholder="Your name"
                  required
                >
              </div>

              <!-- Action Buttons -->
              <div class="d-flex gap-2">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="!isFormValid || saving"
                >
                  <span v-if="saving" class="spinner-border spinner-border-sm me-2" role="status"></span>
                  {{ isEditing ? 'Update Test Run' : 'Create Test Run' }}
                </button>
                <router-link to="/test-runs" class="btn btn-secondary">Cancel</router-link>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- Right Column: Test Case Selection -->
      <div class="col-md-7">
        <div class="card">
          <div class="card-header">
            <div class="d-flex justify-content-between align-items-center">
              <h5 class="mb-0">
                <i class="fas fa-list-check"></i> Test Cases
                <span v-if="selectedTestCases.length" class="badge bg-primary ms-2">{{ selectedTestCases.length }} selected</span>
              </h5>
              <div class="d-flex gap-2">
                <button 
                  v-if="testSuites.length" 
                  @click="selectAllTestCases" 
                  class="btn btn-sm btn-outline-primary"
                >
                  Select All
                </button>
                <button 
                  v-if="selectedTestCases.length" 
                  @click="clearSelection" 
                  class="btn btn-sm btn-outline-secondary"
                >
                  Clear
                </button>
              </div>
            </div>
          </div>
          <div class="card-body p-0">
            <!-- Search and Filter -->
            <div class="p-3 border-bottom bg-light">
              <div class="row g-2">
                <div class="col-md-8">
                  <input
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="Search test cases..."
                    v-model="testCaseSearch"
                  >
                </div>
                <div class="col-md-4">
                  <select class="form-select form-select-sm" v-model="priorityFilter">
                    <option value="">All Priorities</option>
                    <option value="High">High</option>
                    <option value="Medium">Medium</option>
                    <option value="Low">Low</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Test Cases List -->
            <div class="test-cases-container" style="max-height: 500px; overflow-y: auto;">
              <div v-if="!form.project_id" class="p-4 text-center text-muted">
                <i class="fas fa-arrow-left fa-2x mb-3"></i>
                <p>Select a project to view available test cases</p>
              </div>

              <div v-else-if="testSuites.length === 0" class="p-4 text-center text-muted">
                <i class="fas fa-exclamation-triangle fa-2x mb-3"></i>
                <p>No test suites found for this project</p>
                <small>Create test suites and test cases first</small>
              </div>

              <div v-else>
                <!-- Test Suite Groups -->
                <div v-for="suite in testSuites" :key="suite.id" class="border-bottom">
                  <div class="suite-header p-3 bg-light" @click="toggleSuiteExpansion(suite.id)">
                    <div class="d-flex align-items-center justify-content-between">
                      <div class="d-flex align-items-center">
                        <input
                          type="checkbox"
                          class="form-check-input me-2"
                          @change="toggleSuite(suite)"
                          :checked="isSuiteSelected(suite)"
                          :indeterminate="isSuitePartiallySelected(suite)"
                          @click.stop
                        >
                        <span class="fw-semibold">{{ suite.name }}</span>
                        <small class="text-muted ms-2">({{ getFilteredTestCases(suite).length }} cases)</small>
                      </div>
                      <i class="fas" :class="expandedSuites.has(suite.id) ? 'fa-chevron-down' : 'fa-chevron-right'"></i>
                    </div>
                    <small v-if="suite.description" class="text-muted d-block mt-1">{{ suite.description }}</small>
                  </div>

                  <!-- Test Cases (Collapsible) -->
                  <div v-if="expandedSuites.has(suite.id)" class="test-cases-list">
                    <div v-if="getFilteredTestCases(suite).length === 0" class="p-3 text-muted">
                      <small>No test cases match the current filters</small>
                    </div>
                    <div v-else>
                      <div 
                        v-for="testCase in getFilteredTestCases(suite)" 
                        :key="testCase.id" 
                        class="test-case-item p-2 border-bottom"
                        :class="{ 'bg-light': selectedTestCases.includes(testCase.id) }"
                      >
                        <div class="form-check">
                          <input
                            type="checkbox"
                            class="form-check-input"
                            :id="'case-' + testCase.id"
                            :value="testCase.id"
                            v-model="selectedTestCases"
                          >
                          <label class="form-check-label w-100" :for="'case-' + testCase.id">
                            <div class="d-flex justify-content-between align-items-start">
                              <div class="flex-grow-1">
                                <div class="fw-medium">{{ testCase.title }}</div>
                                <small v-if="testCase.description" class="text-muted d-block">{{ testCase.description }}</small>
                              </div>
                              <span class="badge ms-2" :class="getPriorityClass(testCase.priority)">
                                {{ testCase.priority }}
                              </span>
                            </div>
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api.js'
import { showAlert } from '../utils/helpers.js'

export default {
  name: 'TestRunForm',
  props: {
    id: {
      type: String,
      default: null
    }
  },
  data() {
    return {
      loading: true,
      saving: false,
      error: null,
      form: {
        name: '',
        description: '',
        project_id: '',
        repository_id: null,
        branch_name: null,
        tag_name: null,
        created_by: ''
      },
      projects: [],
      testSuites: [],
      selectedTestCases: [],
      testCaseSearch: '',
      priorityFilter: '',
      gitReferenceType: 'branch',
      gitReferences: [],
      expandedSuites: new Set()
    }
  },
  computed: {
    isEditing() {
      return !!this.id
    },
    selectedProject() {
      return this.projects.find(p => p.id == this.form.project_id)
    },
    selectedGitReference: {
      get() {
        return this.gitReferenceType === 'branch' ? this.form.branch_name : this.form.tag_name
      },
      set(value) {
        if (this.gitReferenceType === 'branch') {
          this.form.branch_name = value
          this.form.tag_name = null
        } else {
          this.form.tag_name = value
          this.form.branch_name = null
        }
      }
    },
    isFormValid() {
      return this.form.project_id && this.form.created_by && this.selectedTestCases.length > 0
    }
  },
  async mounted() {
    await this.loadData()
  },
  watch: {
    gitReferenceType() {
      this.form.branch_name = null
      this.form.tag_name = null
      this.loadGitReferences()
    }
  },
  methods: {
    async loadData() {
      try {
        this.loading = true
        this.error = null

        // Load projects
        const projectsResponse = await api.getProjects()
        this.projects = projectsResponse || []

        // If editing, load test run data
        if (this.isEditing) {
          await this.loadTestRun()
        }

      } catch (error) {
        console.error('Error loading data:', error)
        this.error = 'Failed to load data. Please try again.'
      } finally {
        this.loading = false
      }
    },

    async loadTestRun() {
      try {
        const response = await api.getTestRun(this.id)
        const testRun = response

        this.form = {
          name: testRun.name || '',
          description: testRun.description || '',
          project_id: testRun.project_id,
          repository_id: testRun.repository_id,
          branch_name: testRun.branch_name,
          tag_name: testRun.tag_name,
          created_by: testRun.created_by || ''
        }

        this.selectedTestCases = testRun.test_cases?.map(tc => tc.id) || []

        // Load project-specific data
        if (this.form.project_id) {
          await this.onProjectChange()
        }

      } catch (error) {
        console.error('Error loading test run:', error)
        this.error = 'Failed to load test run data.'
      }
    },

    async onProjectChange() {
      if (!this.form.project_id) {
        this.testSuites = []
        this.selectedTestCases = []
        return
      }

      try {
        // Load test suites for project
        const response = await api.getTestSuites(this.form.project_id)
        this.testSuites = response || []

        // Expand all suites by default when loading
        this.testSuites.forEach(suite => {
          this.expandedSuites.add(suite.id)
        })

        // Load git references if project has repository
        if (this.selectedProject?.repository) {
          this.form.repository_id = this.selectedProject.repository.id
          await this.loadGitReferences()
        } else {
          this.form.repository_id = null
          this.gitReferences = []
        }

      } catch (error) {
        console.error('Error loading project data:', error)
        showAlert('Failed to load project data', 'error')
      }
    },

    async loadGitReferences() {
      if (!this.selectedProject?.repository) return

      try {
        const response = await api.getRepository(this.selectedProject.repository.id)
        const repo = response
        
        if (this.gitReferenceType === 'branch') {
          this.gitReferences = repo.branches || []
        } else {
          this.gitReferences = repo.tags || []
        }
      } catch (error) {
        console.error('Error loading git references:', error)
        this.gitReferences = []
      }
    },

    getFilteredTestCases(suite) {
      if (!suite.test_cases) return []

      return suite.test_cases.filter(testCase => {
        // Search filter
        if (this.testCaseSearch) {
          const searchLower = this.testCaseSearch.toLowerCase()
          if (!testCase.title.toLowerCase().includes(searchLower) &&
              !testCase.description?.toLowerCase().includes(searchLower)) {
            return false
          }
        }

        // Priority filter
        if (this.priorityFilter && testCase.priority !== this.priorityFilter) {
          return false
        }

        return true
      })
    },

    toggleSuiteExpansion(suiteId) {
      if (this.expandedSuites.has(suiteId)) {
        this.expandedSuites.delete(suiteId)
      } else {
        this.expandedSuites.add(suiteId)
      }
    },

    toggleSuite(suite) {
      const testCases = this.getFilteredTestCases(suite)
      const testCaseIds = testCases.map(tc => tc.id)
      
      if (this.isSuiteSelected(suite)) {
        // Remove all test cases from this suite
        this.selectedTestCases = this.selectedTestCases.filter(id => !testCaseIds.includes(id))
      } else {
        // Add all test cases from this suite
        const newSelections = testCaseIds.filter(id => !this.selectedTestCases.includes(id))
        this.selectedTestCases.push(...newSelections)
      }
    },

    isSuiteSelected(suite) {
      const testCases = this.getFilteredTestCases(suite)
      if (testCases.length === 0) return false
      return testCases.every(tc => this.selectedTestCases.includes(tc.id))
    },

    isSuitePartiallySelected(suite) {
      const testCases = this.getFilteredTestCases(suite)
      if (testCases.length === 0) return false
      const selectedCount = testCases.filter(tc => this.selectedTestCases.includes(tc.id)).length
      return selectedCount > 0 && selectedCount < testCases.length
    },

    selectAllTestCases() {
      this.selectedTestCases = []
      this.testSuites.forEach(suite => {
        const testCases = this.getFilteredTestCases(suite)
        testCases.forEach(tc => {
          if (!this.selectedTestCases.includes(tc.id)) {
            this.selectedTestCases.push(tc.id)
          }
        })
      })
    },

    clearSelection() {
      this.selectedTestCases = []
    },

    getPriorityClass(priority) {
      switch (priority?.toLowerCase()) {
        case 'high': return 'bg-danger'
        case 'medium': return 'bg-warning'
        case 'low': return 'bg-success'
        default: return 'bg-secondary'
      }
    },

    async saveTestRun() {
      if (!this.isFormValid) {
        showAlert('Please fill in all required fields and select at least one test case', 'error')
        return
      }

      try {
        this.saving = true

        const payload = {
          name: this.form.name.trim() || null, // Auto-generate if empty
          description: this.form.description.trim(),
          project_id: parseInt(this.form.project_id),
          repository_id: this.form.repository_id,
          branch_name: this.form.branch_name,
          tag_name: this.form.tag_name,
          created_by: this.form.created_by.trim(),
          test_case_ids: this.selectedTestCases
        }

        if (this.isEditing) {
          await api.updateTestRun(this.id, payload)
          showAlert('Test run updated successfully', 'success')
        } else {
          await api.createTestRun(payload)
          showAlert('Test run created successfully', 'success')
        }

        this.$router.push('/test-runs')

      } catch (error) {
        console.error('Error saving test run:', error)
        showAlert('Failed to save test run. Please try again.', 'error')
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style scoped>
.test-run-form-container {
  padding: 20px;
}

.suite-header {
  cursor: pointer;
  transition: background-color 0.2s;
}

.suite-header:hover {
  background-color: #f8f9fa !important;
}

.test-case-item {
  transition: background-color 0.2s;
}

.test-case-item:hover {
  background-color: #f8f9fa;
}

.test-cases-container::-webkit-scrollbar {
  width: 6px;
}

.test-cases-container::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.test-cases-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.test-cases-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>