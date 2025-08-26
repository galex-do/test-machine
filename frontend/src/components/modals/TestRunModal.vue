<template>
  <div class="modal fade" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" @click.self="close">
    <div class="modal-dialog modal-xl">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            <i class="fas fa-play-circle text-primary"></i>
            {{ isEditing ? 'Edit Test Run' : 'Create New Test Run' }}
          </h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>

        <div class="modal-body">
          <!-- Loading State -->
          <div v-if="loading" class="text-center py-4">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>

          <!-- Form -->
          <form v-else @submit.prevent="saveTestRun">
            <div class="row">
              <!-- Basic Information -->
              <div class="col-md-6">
                <div class="mb-3">
                  <label for="name" class="form-label">Name</label>
                  <input
                    type="text"
                    class="form-control"
                    id="name"
                    v-model="form.name"
                    placeholder="Auto-generated if left empty"
                  >
                  <div class="form-text">If left empty, name will be auto-generated as: &lt;project&gt;-&lt;branch/tag&gt;-&lt;datetime&gt;</div>
                </div>

                <div class="mb-3">
                  <label for="description" class="form-label">Description</label>
                  <textarea
                    class="form-control"
                    id="description"
                    v-model="form.description"
                    rows="3"
                  ></textarea>
                </div>

                <div class="mb-3">
                  <label for="project" class="form-label">Project *</label>
                  <select
                    class="form-select"
                    id="project"
                    v-model="form.project_id"
                    @change="onProjectChange"
                    required
                  >
                    <option value="">Select a project</option>
                    <option v-for="project in projects" :key="project.id" :value="project.id">
                      {{ project.name }}
                    </option>
                  </select>
                </div>

                <div class="mb-3">
                  <label for="created_by" class="form-label">Created By</label>
                  <input
                    type="text"
                    class="form-control"
                    id="created_by"
                    v-model="form.created_by"
                    placeholder="Enter your name"
                  >
                </div>
              </div>

              <!-- Git Reference -->
              <div class="col-md-6">
                <div v-if="selectedProject && selectedProject.repository" class="mb-3">
                  <label class="form-label">Git Reference</label>
                  <div class="card">
                    <div class="card-body">
                      <h6 class="card-title">
                        <i class="fas fa-code-branch"></i> {{ selectedProject.repository.name }}
                      </h6>
                      
                      <!-- Branch/Tag Selection -->
                      <div class="mb-3">
                        <div class="form-check">
                          <input
                            class="form-check-input"
                            type="radio"
                            name="gitRef"
                            id="noBranch"
                            :value="'none'"
                            v-model="gitRefType"
                          >
                          <label class="form-check-label" for="noBranch">
                            No specific branch/tag
                          </label>
                        </div>
                        <div class="form-check">
                          <input
                            class="form-check-input"
                            type="radio"
                            name="gitRef"
                            id="branch"
                            :value="'branch'"
                            v-model="gitRefType"
                          >
                          <label class="form-check-label" for="branch">
                            Branch
                          </label>
                        </div>
                        <div class="form-check">
                          <input
                            class="form-check-input"
                            type="radio"
                            name="gitRef"
                            id="tag"
                            :value="'tag'"
                            v-model="gitRefType"
                          >
                          <label class="form-check-label" for="tag">
                            Tag
                          </label>
                        </div>
                      </div>

                      <!-- Branch Selection -->
                      <div v-if="gitRefType === 'branch'" class="mb-3">
                        <label for="branch_name" class="form-label">Branch</label>
                        <select class="form-select" id="branch_name" v-model="form.branch_name">
                          <option value="">Select a branch</option>
                          <option v-for="branch in selectedProject.repository.branches" :key="branch.id" :value="branch.name">
                            {{ branch.name }} <span v-if="branch.is_default">(default)</span>
                          </option>
                        </select>
                      </div>

                      <!-- Tag Selection -->
                      <div v-if="gitRefType === 'tag'" class="mb-3">
                        <label for="tag_name" class="form-label">Tag</label>
                        <select class="form-select" id="tag_name" v-model="form.tag_name">
                          <option value="">Select a tag</option>
                          <option v-for="tag in selectedProject.repository.tags" :key="tag.id" :value="tag.name">
                            {{ tag.name }}
                          </option>
                        </select>
                      </div>
                    </div>
                  </div>
                </div>

                <div v-else-if="selectedProject && !selectedProject.repository" class="alert alert-info">
                  <i class="fas fa-info-circle"></i>
                  This project has no Git repository configured.
                </div>
              </div>
            </div>

            <!-- Test Case Selection -->
            <div v-if="selectedProject" class="mt-4">
              <h6>
                <i class="fas fa-list-check"></i> Select Test Cases
                <span class="badge bg-secondary ms-2">{{ selectedTestCases.length }} selected</span>
              </h6>
              
              <div v-if="testSuites.length === 0" class="alert alert-warning">
                <i class="fas fa-exclamation-triangle"></i>
                No test suites found for this project. Create test suites and test cases first.
              </div>

              <div v-else>
                <div v-for="suite in testSuites" :key="suite.id" class="card mb-3">
                  <div class="card-header">
                    <div class="d-flex justify-content-between align-items-center">
                      <div class="form-check">
                        <input
                          class="form-check-input"
                          type="checkbox"
                          :id="'suite-' + suite.id"
                          @change="toggleSuite(suite)"
                          :checked="isSuiteSelected(suite)"
                          :indeterminate="isSuitePartiallySelected(suite)"
                        >
                        <label class="form-check-label fw-bold" :for="'suite-' + suite.id">
                          {{ suite.name }}
                        </label>
                      </div>
                      <span class="badge bg-info">{{ getActiveTestCases(suite).length }} cases</span>
                    </div>
                    <small v-if="suite.description" class="text-muted">{{ suite.description }}</small>
                  </div>
                  <div class="card-body">
                    <div v-if="getActiveTestCases(suite).length === 0" class="text-muted">
                      No active test cases in this suite.
                    </div>
                    <div v-else class="row">
                      <div v-for="testCase in getActiveTestCases(suite)" :key="testCase.id" class="col-md-6 mb-2">
                        <div class="form-check">
                          <input
                            class="form-check-input"
                            type="checkbox"
                            :id="'case-' + testCase.id"
                            :value="testCase.id"
                            v-model="selectedTestCases"
                          >
                          <label class="form-check-label" :for="'case-' + testCase.id">
                            {{ testCase.title }}
                            <small class="text-muted d-block">Priority: {{ testCase.priority }}</small>
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary" 
            @click="saveTestRun"
            :disabled="!isFormValid || saving"
          >
            <span v-if="saving" class="spinner-border spinner-border-sm me-2" role="status"></span>
            {{ isEditing ? 'Update Test Run' : 'Create Test Run' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../../services/api.js'
import { showAlert } from '../../utils/helpers.js'

export default {
  name: 'TestRunModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    testRun: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
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
      gitRefType: 'none',
      loading: false,
      saving: false
    }
  },
  computed: {
    isEditing() {
      return !!this.testRun
    },
    selectedProject() {
      return this.projects.find(p => p.id === this.form.project_id)
    },
    isFormValid() {
      return this.form.project_id && this.selectedTestCases.length > 0
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.loadData()
        if (this.testRun) {
          this.populateForm()
        } else {
          this.resetForm()
        }
      }
    },
    gitRefType(newVal) {
      if (newVal === 'none') {
        this.form.branch_name = null
        this.form.tag_name = null
      } else if (newVal === 'branch') {
        this.form.tag_name = null
      } else if (newVal === 'tag') {
        this.form.branch_name = null
      }
    }
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        this.projects = await api.getProjectsWithRepositories()
      } catch (error) {
        showAlert('Error loading projects: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },

    async onProjectChange() {
      if (this.form.project_id) {
        try {
          // Load test suites for the selected project
          this.testSuites = await api.getTestSuitesByProject(this.form.project_id)
          this.selectedTestCases = []
          
          // Load repository details with branches/tags if project has a repository
          const selectedProject = this.selectedProject
          if (selectedProject && selectedProject.repository_id) {
            const repositoryDetails = await api.getRepositoryDetails(selectedProject.repository_id)
            // Update the project's repository with detailed info
            const projectIndex = this.projects.findIndex(p => p.id === this.form.project_id)
            if (projectIndex !== -1) {
              this.projects[projectIndex].repository = repositoryDetails
            }
          }
        } catch (error) {
          showAlert('Error loading project data: ' + error.message, 'danger')
        }
      }
    },

    getActiveTestCases(suite) {
      if (!suite || !suite.test_cases || !Array.isArray(suite.test_cases)) {
        return []
      }
      return suite.test_cases.filter(tc => tc && tc.status === 'Active')
    },

    toggleSuite(suite) {
      const activeTestCases = this.getActiveTestCases(suite)
      const suiteTestCaseIds = activeTestCases.map(tc => tc.id)
      const allSelected = suiteTestCaseIds.every(id => this.selectedTestCases.includes(id))

      if (allSelected) {
        // Remove all test cases from this suite
        this.selectedTestCases = this.selectedTestCases.filter(id => !suiteTestCaseIds.includes(id))
      } else {
        // Add all test cases from this suite
        suiteTestCaseIds.forEach(id => {
          if (!this.selectedTestCases.includes(id)) {
            this.selectedTestCases.push(id)
          }
        })
      }
    },

    isSuiteSelected(suite) {
      const activeTestCases = this.getActiveTestCases(suite)
      const suiteTestCaseIds = activeTestCases.map(tc => tc.id)
      return suiteTestCaseIds.length > 0 && suiteTestCaseIds.every(id => this.selectedTestCases.includes(id))
    },

    isSuitePartiallySelected(suite) {
      const activeTestCases = this.getActiveTestCases(suite)
      const suiteTestCaseIds = activeTestCases.map(tc => tc.id)
      const selectedCount = suiteTestCaseIds.filter(id => this.selectedTestCases.includes(id)).length
      return selectedCount > 0 && selectedCount < suiteTestCaseIds.length
    },

    populateForm() {
      this.form = {
        name: this.testRun.name,
        description: this.testRun.description || '',
        project_id: this.testRun.project_id,
        repository_id: this.testRun.repository_id,
        branch_name: this.testRun.branch_name,
        tag_name: this.testRun.tag_name,
        created_by: this.testRun.created_by || ''
      }
      
      if (this.testRun.branch_name) {
        this.gitRefType = 'branch'
      } else if (this.testRun.tag_name) {
        this.gitRefType = 'tag'
      } else {
        this.gitRefType = 'none'
      }

      this.selectedTestCases = this.testRun.test_cases?.map(tc => tc.test_case_id) || []
    },

    resetForm() {
      this.form = {
        name: '',
        description: '',
        project_id: '',
        repository_id: null,
        branch_name: null,
        tag_name: null,
        created_by: ''
      }
      this.selectedTestCases = []
      this.gitRefType = 'none'
      this.testSuites = []
    },

    async saveTestRun() {
      if (!this.isFormValid) return

      this.saving = true
      try {
        const payload = {
          ...this.form,
          test_case_ids: this.selectedTestCases
        }

        if (this.isEditing) {
          await api.updateTestRun(this.testRun.id, payload)
        } else {
          await api.createTestRun(payload)
        }

        this.$emit('saved')
      } catch (error) {
        showAlert('Error saving test run: ' + error.message, 'danger')
      } finally {
        this.saving = false
      }
    },

    close() {
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.modal.show {
  backdrop-filter: blur(3px);
}

.form-check-input:indeterminate {
  background-color: #007bff;
  border-color: #007bff;
}

.form-check-input:indeterminate:before {
  content: '−';
  color: white;
  font-weight: bold;
  display: block;
  text-align: center;
  line-height: 1;
}
</style>