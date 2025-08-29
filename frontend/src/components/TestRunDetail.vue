<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link to="/test-runs">Test Runs</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(testRun?.name, 30) }}</li>
      </ol>
    </nav>

    <!-- Test Run Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">
          <i class="fas fa-play-circle text-primary"></i>
          {{ testRun?.name }}
        </h1>
        <p class="text-muted">{{ testRun?.description || 'No description available' }}</p>
        <div class="mt-2" v-if="testRun">
          <span class="status-badge me-2" :class="getStatusBadgeClass(testRun.status)">
            {{ testRun.status }}
          </span>
          <span v-if="testRun.project" class="badge bg-info me-2">
            <i class="fas fa-project-diagram"></i> {{ testRun.project.name }}
          </span>
          <span v-if="testRun.branch_name" class="badge bg-primary me-2">
            <i class="fas fa-code-branch"></i> {{ testRun.branch_name }}
          </span>
          <small class="text-muted d-block mt-1">
            Created: {{ formatDate(testRun.created_at) }}
          </small>
        </div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0" v-if="testRun">
        <div class="btn-group me-2" role="group">
          <!-- Time Management Buttons -->
          <button 
            v-if="testRun.status === 'Not Started'"
            @click="startTestRun"
            class="btn btn-success"
            :disabled="loading"
            title="Start Test Run"
          >
            <i class="fas fa-play"></i> Start
          </button>
          <button 
            v-else-if="testRun.status === 'In Progress'"
            @click="pauseTestRun"
            class="btn btn-warning"
            :disabled="loading"
            title="Pause Test Run"
          >
            <i class="fas fa-pause"></i> Pause
          </button>
          <button 
            v-if="testRun.status === 'In Progress'"
            @click="finishTestRun"
            class="btn btn-primary"
            :disabled="loading"
            title="Finish Test Run"
          >
            <i class="fas fa-stop"></i> Finish
          </button>
        </div>
        <div class="btn-group" role="group">
          <!-- Management Buttons -->
          <router-link 
            v-if="testRun.status !== 'Completed'"
            :to="`/test-runs/${testRun.id}/edit`"
            class="btn btn-outline-secondary"
            title="Edit Test Run"
          >
            <i class="fas fa-edit"></i> Edit
          </router-link>
          <button 
            v-if="testRun.status === 'Not Started'"
            @click="deleteTestRun"
            class="btn btn-outline-danger"
            :disabled="loading"
            title="Delete Test Run"
          >
            <i class="fas fa-trash"></i> Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Test Run Statistics -->
    <div class="row mb-4" v-if="testRun">
      <div class="col-md-2">
        <div class="card">
          <div class="card-body text-center">
            <h5>Status</h5>
            <span class="status-badge" :class="getStatusBadgeClass(testRun.status)">
              {{ testRun.status }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-2">
        <div class="card">
          <div class="card-body text-center">
            <h5>Test Cases</h5>
            <h3 class="text-primary">{{ testCases?.length || 0 }}</h3>
          </div>
        </div>
      </div>
      <div class="col-md-2">
        <div class="card">
          <div class="card-body text-center">
            <h5>Progress</h5>
            <div class="progress mb-2" style="height: 6px;">
              <div 
                class="progress-bar" 
                :class="getProgressBarClass()"
                :style="{ width: getProgressPercentage() + '%' }"
              ></div>
            </div>
            <small>{{ getProgressText() }}</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Duration</h5>
            <p class="mb-0">{{ calculateDuration(testRun.started_at, testRun.completed_at) }}</p>
            <small v-if="testRun.status === 'In Progress'" class="text-muted">
              {{ getElapsedTime() }}
            </small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Executed By</h5>
            <p class="mb-0">{{ testRun.executed_by || 'Not started' }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Running Interface -->
    <div v-if="testRun && testRun.status === 'In Progress'" class="card mb-4">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5 class="mb-0"><i class="fas fa-play-circle text-success"></i> Test Running</h5>
        <div class="d-flex gap-2">
          <button 
            @click="previousTestCase" 
            class="btn btn-outline-secondary btn-sm"
            :disabled="currentTestCaseIndex === 0"
          >
            <i class="fas fa-chevron-left"></i> Previous
          </button>
          <span class="badge bg-primary align-self-center">
            {{ currentTestCaseIndex + 1 }} of {{ testCases?.length || 0 }}
          </span>
          <button 
            @click="nextTestCase" 
            class="btn btn-outline-secondary btn-sm"
            :disabled="currentTestCaseIndex >= (testCases?.length || 0) - 1"
          >
            Next <i class="fas fa-chevron-right"></i>
          </button>
        </div>
      </div>
      <div class="card-body" v-if="currentTestCase">
        <div class="row">
          <div class="col-md-8">
            <h5>{{ currentTestCase.title }}</h5>
            <p class="text-muted">{{ currentTestCase.description }}</p>
            
            <!-- Test Steps -->
            <div v-if="currentTestCase.test_steps && currentTestCase.test_steps.length > 0" class="mb-3">
              <h6><i class="fas fa-list-ol"></i> Test Steps:</h6>
              <div v-for="step in currentTestCase.test_steps" :key="step.id" class="test-step-item mb-2">
                <div class="d-flex">
                  <span class="step-number me-2">{{ step.step_number }}</span>
                  <div class="flex-grow-1">
                    <p class="mb-1"><strong>Step:</strong> {{ step.description }}</p>
                    <p class="mb-0 text-muted"><strong>Expected:</strong> {{ step.expected_result }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="col-md-4">
            <div class="border rounded p-3 bg-light">
              <h6><i class="fas fa-clipboard-check"></i> Test Result</h6>
              
              <!-- Status Selection -->
              <div class="mb-3">
                <label class="form-label">Status</label>
                <div class="btn-group d-flex" role="group">
                  <input type="radio" class="btn-check" :id="`pass-${currentTestCase.id}`" :value="'Pass'" v-model="currentResult.status">
                  <label class="btn btn-outline-success btn-sm" :for="`pass-${currentTestCase.id}`">
                    <i class="fas fa-check"></i> Pass
                  </label>
                  
                  <input type="radio" class="btn-check" :id="`fail-${currentTestCase.id}`" :value="'Fail'" v-model="currentResult.status">
                  <label class="btn btn-outline-danger btn-sm" :for="`fail-${currentTestCase.id}`">
                    <i class="fas fa-times"></i> Fail
                  </label>
                  
                  <input type="radio" class="btn-check" :id="`skip-${currentTestCase.id}`" :value="'Skip'" v-model="currentResult.status">
                  <label class="btn btn-outline-warning btn-sm" :for="`skip-${currentTestCase.id}`">
                    <i class="fas fa-forward"></i> Skip
                  </label>
                </div>
              </div>
              
              <!-- Comments -->
              <div class="mb-3">
                <label class="form-label">Comments</label>
                <textarea 
                  class="form-control" 
                  rows="4" 
                  placeholder="Add comments about test execution..."
                  v-model="currentResult.notes"
                ></textarea>
              </div>
              
              <!-- Save Button -->
              <button 
                @click="saveTestResult" 
                class="btn btn-primary w-100"
                :disabled="!currentResult.status || saving"
              >
                <i class="fas fa-save"></i> 
                {{ saving ? 'Saving...' : 'Save Result' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Test Cases List -->
    <div class="card" v-if="testRun">
      <div class="card-header">
        <h5><i class="fas fa-list-check"></i> Test Cases ({{ testCases?.length || 0 }})</h5>
      </div>
      <div class="card-body">
        <div v-if="loading" class="text-center py-4">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
          <p class="mt-2">Loading test cases...</p>
        </div>
        
        <div v-else-if="!testCases || testCases.length === 0" class="text-center py-4 text-muted">
          <i class="fas fa-inbox fa-3x mb-3"></i>
          <h5>No Test Cases</h5>
          <p>No test cases are associated with this test run.</p>
        </div>
        
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th width="5%">#</th>
                <th>Test Case</th>
                <th width="15%">Priority</th>
                <th width="15%">Status</th>
                <th width="20%">Last Updated</th>
                <th width="10%">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(testCase, index) in testCases" 
                :key="testCase.id"
                :class="{ 'table-active': testRun.status === 'In Progress' && index === currentTestCaseIndex }"
              >
                <td>{{ index + 1 }}</td>
                <td>
                  <div>
                    <strong>{{ testCase.title }}</strong>
                    <div v-if="testCase.description" class="text-muted small">{{ truncateText(testCase.description, 80) }}</div>
                  </div>
                </td>
                <td>
                  <span class="priority-badge" :class="getPriorityBadgeClass(testCase.priority)">
                    {{ testCase.priority }}
                  </span>
                </td>
                <td>
                  <span class="status-badge" :class="getTestCaseStatusBadgeClass(testCase.result_status || 'Not Executed')">
                    {{ testCase.result_status || 'Not Executed' }}
                  </span>
                </td>
                <td class="small text-muted">
                  {{ testCase.result_updated_at ? formatDate(testCase.result_updated_at) : '-' }}
                </td>
                <td>
                  <div class="btn-group btn-group-sm" role="group">
                    <button 
                      v-if="testRun.status === 'In Progress'"
                      @click="goToTestCase(index)"
                      class="btn btn-outline-primary"
                      :class="{ 'btn-primary': index === currentTestCaseIndex }"
                      title="Execute Test"
                    >
                      <i class="fas fa-play"></i>
                    </button>
                    <router-link 
                      :to="`/project/${testCase.test_suite?.project_id}/test-suite/${testCase.test_suite_id}/test-case/${testCase.id}`"
                      class="btn btn-outline-info"
                      title="View Details"
                    >
                      <i class="fas fa-eye"></i>
                    </router-link>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Test Run Information -->
    <div class="card mt-4" v-if="testRun">
      <div class="card-header">
        <h5><i class="fas fa-info-circle"></i> Test Run Information</h5>
      </div>
      <div class="card-body">
        <div class="row">
          <div class="col-md-6">
            <p><strong>Project:</strong> 
              <router-link 
                v-if="testRun.project" 
                :to="`/project/${testRun.project.id}`"
                class="text-decoration-none"
              >
                {{ testRun.project.name }}
              </router-link>
              <span v-else>N/A</span>
            </p>
            <p><strong>Git Reference:</strong> 
              <span v-if="testRun.branch_name" class="badge bg-primary">
                <i class="fas fa-code-branch"></i> {{ testRun.branch_name }}
              </span>
              <span v-else-if="testRun.tag_name" class="badge bg-warning">
                <i class="fas fa-tag"></i> {{ testRun.tag_name }}
              </span>
              <span v-else class="text-muted">No git reference</span>
            </p>
            <p><strong>Started At:</strong> {{ formatDate(testRun.started_at) }}</p>
            <p><strong>Completed At:</strong> {{ formatDate(testRun.completed_at) }}</p>
          </div>
          <div class="col-md-6">
            <p><strong>Description:</strong></p>
            <div class="border rounded p-3 bg-light">
              {{ testRun.description || 'No description available' }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading test run details...</p>
    </div>
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass } from '../utils/helpers.js'

export default {
  name: 'TestRunDetail',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      testRun: null,
      testCases: [],
      currentTestCaseIndex: 0,
      currentResult: {
        status: '',
        notes: ''
      },
      loading: true,
      saving: false,
      elapsedTimer: null,
      elapsedTime: ''
    }
  },
  mounted() {
    this.loadData()
  },
  watch: {
    id() {
      this.loadData()
    }
  },
  computed: {
    currentTestCase() {
      return this.testCases && this.testCases[this.currentTestCaseIndex] || null
    }
  },
  methods: {
    formatDate,
    showLoading,
    truncateText,
    getStatusBadgeClass,
    getPriorityBadgeClass,
    
    async loadData() {
      this.loading = true
      try {
        const [testRunData] = await Promise.all([
          api.getTestRun(this.id)
        ])
        
        this.testRun = testRunData
        
        // Load test cases associated with this test run
        if (testRunData?.test_cases) {
          this.testCases = testRunData.test_cases
        }
        
        // Initialize current result if in progress
        if (this.testRun?.status === 'In Progress' && this.currentTestCase) {
          this.loadCurrentTestResult()
        }
        
        // Start elapsed time timer if in progress
        if (this.testRun?.status === 'In Progress') {
          this.startElapsedTimer()
        }
        
      } catch (error) {
        showAlert('Error loading test run data: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },
    
    loadCurrentTestResult() {
      if (!this.currentTestCase) return
      
      this.currentResult = {
        status: this.currentTestCase.result_status || '',
        notes: this.currentTestCase.result_notes || ''
      }
    },

    // Test Run Control Methods
    async startTestRun() {
      try {
        this.loading = true
        await api.startTestRun(this.id)
        showAlert('Test run started successfully!', 'success')
        await this.loadData()
      } catch (error) {
        showAlert('Error starting test run: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },
    
    async pauseTestRun() {
      if (!confirm('Are you sure you want to pause this test run?')) return
      
      try {
        this.loading = true
        await api.pauseTestRun(this.id)
        showAlert('Test run paused successfully!', 'success')
        this.stopElapsedTimer()
        await this.loadData()
      } catch (error) {
        showAlert('Error pausing test run: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },
    
    async finishTestRun() {
      if (!confirm('Are you sure you want to finish this test run? This will mark it as completed.')) return
      
      try {
        this.loading = true
        await api.finishTestRun(this.id)
        showAlert('Test run finished successfully!', 'success')
        this.stopElapsedTimer()
        await this.loadData()
      } catch (error) {
        showAlert('Error finishing test run: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },
    
    async deleteTestRun() {
      if (!confirm(`Are you sure you want to delete test run "${this.testRun?.name}"? This action cannot be undone.`)) return
      
      try {
        this.loading = true
        await api.deleteTestRun(this.id)
        showAlert('Test run deleted successfully!', 'success')
        this.$router.push('/test-runs')
      } catch (error) {
        showAlert('Error deleting test run: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },
    
    // Test Case Navigation
    previousTestCase() {
      if (this.currentTestCaseIndex > 0) {
        this.currentTestCaseIndex--
        this.loadCurrentTestResult()
      }
    },
    
    nextTestCase() {
      if (this.currentTestCaseIndex < (this.testCases?.length || 0) - 1) {
        this.currentTestCaseIndex++
        this.loadCurrentTestResult()
      }
    },
    
    goToTestCase(index) {
      this.currentTestCaseIndex = index
      this.loadCurrentTestResult()
    },
    
    // Test Result Management
    async saveTestResult() {
      if (!this.currentTestCase || !this.currentResult.status) return
      
      try {
        this.saving = true
        const updateData = {
          status: this.currentResult.status,
          notes: this.currentResult.notes,
          executed_by: 'Current User' // TODO: Get from auth context
        }
        
        await api.updateTestRunCase(this.id, this.currentTestCase.id, updateData)
        
        // Update local test case data
        this.currentTestCase.result_status = this.currentResult.status
        this.currentTestCase.result_notes = this.currentResult.notes
        this.currentTestCase.result_updated_at = new Date().toISOString()
        
        showAlert('Test result saved successfully!', 'success')
        
        // Auto-advance to next test case if not the last one
        if (this.currentTestCaseIndex < (this.testCases?.length || 0) - 1) {
          setTimeout(() => {
            this.nextTestCase()
          }, 1000)
        }
        
      } catch (error) {
        showAlert('Error saving test result: ' + error.message, 'danger')
      } finally {
        this.saving = false
      }
    },
    
    // Progress and Status Methods
    getProgressPercentage() {
      if (!this.testCases || this.testCases.length === 0) return 0
      
      const executed = this.testCases.filter(tc => tc.result_status && tc.result_status !== 'Not Executed').length
      return Math.round((executed / this.testCases.length) * 100)
    },
    
    getProgressText() {
      if (!this.testCases || this.testCases.length === 0) return '0/0'
      
      const executed = this.testCases.filter(tc => tc.result_status && tc.result_status !== 'Not Executed').length
      return `${executed}/${this.testCases.length}`
    },
    
    getProgressBarClass() {
      if (!this.testCases || this.testCases.length === 0) return 'bg-secondary'
      
      const passed = this.testCases.filter(tc => tc.result_status === 'Pass').length
      const failed = this.testCases.filter(tc => tc.result_status === 'Fail').length
      const total = this.testCases.length
      
      if (failed > 0) return 'bg-danger'
      if (passed === total) return 'bg-success'
      return 'bg-warning'
    },
    
    getTestCaseStatusBadgeClass(status) {
      const classes = {
        'Pass': 'badge bg-success',
        'Fail': 'badge bg-danger', 
        'Skip': 'badge bg-warning',
        'Not Executed': 'badge bg-secondary'
      }
      return classes[status] || 'badge bg-secondary'
    },
    
    // Timer Methods
    startElapsedTimer() {
      if (!this.testRun?.started_at) return
      
      this.elapsedTimer = setInterval(() => {
        this.updateElapsedTime()
      }, 1000)
      
      this.updateElapsedTime()
    },
    
    stopElapsedTimer() {
      if (this.elapsedTimer) {
        clearInterval(this.elapsedTimer)
        this.elapsedTimer = null
      }
    },
    
    updateElapsedTime() {
      if (!this.testRun?.started_at) return
      
      const start = new Date(this.testRun.started_at)
      const now = new Date()
      const elapsed = now - start
      
      const hours = Math.floor(elapsed / 3600000)
      const minutes = Math.floor((elapsed % 3600000) / 60000)
      const seconds = Math.floor((elapsed % 60000) / 1000)
      
      if (hours > 0) {
        this.elapsedTime = `${hours}h ${minutes}m ${seconds}s`
      } else if (minutes > 0) {
        this.elapsedTime = `${minutes}m ${seconds}s`
      } else {
        this.elapsedTime = `${seconds}s`
      }
    },
    
    getElapsedTime() {
      return this.elapsedTime || '0s'
    },
    
    calculateDuration(startedAt, completedAt) {
      if (!startedAt || !completedAt) return 'N/A'
      
      const start = new Date(startedAt)
      const end = new Date(completedAt)
      const durationMs = end - start
      
      if (durationMs < 0) return 'N/A'
      
      const hours = Math.floor(durationMs / 3600000)
      const minutes = Math.floor((durationMs % 3600000) / 60000)
      const seconds = Math.floor((durationMs % 60000) / 1000)
      
      if (hours > 0) {
        return `${hours}h ${minutes}m ${seconds}s`
      } else if (minutes > 0) {
        return `${minutes}m ${seconds}s`
      } else {
        return `${seconds}s`
      }
    }
  },
  
  beforeUnmount() {
    this.stopElapsedTimer()
  }
}
</script>

<style scoped>
.status-badge {
  font-size: 0.85rem;
  padding: 0.4em 0.8em;
}

.priority-badge {
  font-size: 0.75rem;
  padding: 0.3em 0.6em;
}

.test-step-item {
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  padding: 10px;
}

.step-number {
  background-color: var(--bs-primary);
  color: white;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8em;
  font-weight: bold;
  flex-shrink: 0;
}

.table-active {
  background-color: rgba(13, 110, 253, 0.1) !important;
  border-left: 3px solid var(--bs-primary);
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
}

.progress {
  border-radius: 10px;
}

.card {
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  border: 1px solid rgba(0, 0, 0, 0.125);
  margin-bottom: 1rem;
}

.card-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid rgba(0, 0, 0, 0.125);
}

.card-header h5 {
  margin-bottom: 0;
  font-weight: 600;
}

.btn-check:checked + .btn-outline-success {
  background-color: var(--bs-success);
  border-color: var(--bs-success);
  color: white;
}

.btn-check:checked + .btn-outline-danger {
  background-color: var(--bs-danger);
  border-color: var(--bs-danger);
  color: white;
}

.btn-check:checked + .btn-outline-warning {
  background-color: var(--bs-warning);
  border-color: var(--bs-warning);
  color: white;
}

.text-muted {
  font-size: 0.9em;
}

.badge {
  font-size: 0.75rem;
}

@media (max-width: 768px) {
  .btn-toolbar {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .btn-group {
    width: 100%;
  }
  
  .d-flex.gap-2 {
    flex-direction: column;
    gap: 0.25rem;
  }
}
</style>