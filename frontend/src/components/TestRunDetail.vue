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

    <!-- Main Content Layout -->
    <div class="row" v-if="testRun">
      <!-- Left Sidebar: Test Cases List -->
      <div class="col-md-3">
        <div class="card">
          <div class="card-header">
            <h6 class="mb-0"><i class="fas fa-list-check"></i> Test Cases</h6>
          </div>
          <div class="card-body p-0">
            <div v-if="loading" class="text-center py-3">
              <div class="spinner-border spinner-border-sm text-primary"></div>
            </div>
            <div v-else-if="!testCases || testCases.length === 0" class="text-center py-3 text-muted">
              <small>No test cases</small>
            </div>
            <div v-else class="list-group list-group-flush">
              <button
                v-for="(testCase, index) in testCases" 
                :key="testCase.id"
                @click="selectTestCase(index)"
                class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
                :class="{ 
                  'active': testRun.status === 'In Progress' && index === currentTestCaseIndex,
                  'list-group-item-success': testCase.Status === 'Pass',
                  'list-group-item-danger': testCase.Status === 'Fail',
                  'list-group-item-warning': testCase.Status === 'Skip'
                }"
              >
                <div>
                  <div class="fw-bold">{{ index + 1 }}</div>
                  <small class="text-muted">{{ truncateText(testCase.TestCase?.Title || testCase.title, 25) }}</small>
                </div>
                <span class="badge" :class="{
                  'bg-success': testCase.Status === 'Pass',
                  'bg-danger': testCase.Status === 'Fail', 
                  'bg-warning text-dark': testCase.Status === 'Skip',
                  'bg-secondary': testCase.Status === 'Not Executed' || !testCase.Status
                }">
                  {{ testCase.Status || 'Not Executed' }}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Main Panel: Test Running -->
      <div class="col-md-9">
        <div class="card" v-if="testRun?.status === 'In Progress' && currentTestCase">
          <div class="card-header bg-primary text-white">
            <div class="row align-items-center">
              <div class="col">
                <h5 class="mb-0">
                  <i class="fas fa-play-circle"></i>
                  Test Running - Case {{ currentTestCaseIndex + 1 }} of {{ testCases?.length || 0 }}
                </h5>
              </div>
              <div class="col-auto">
                <div class="btn-group btn-group-sm" role="group">
                  <button 
                    @click="previousTestCase" 
                    :disabled="currentTestCaseIndex <= 0"
                    class="btn btn-outline-light"
                    title="Previous Test Case"
                  >
                    <i class="fas fa-chevron-left"></i>
                  </button>
                  <button 
                    @click="nextTestCase" 
                    :disabled="currentTestCaseIndex >= (testCases?.length || 0) - 1"
                    class="btn btn-outline-light"
                    title="Next Test Case"
                  >
                    <i class="fas fa-chevron-right"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div class="card-body">
            <!-- Test Case Details -->
            <div class="mb-4">
              <h5>{{ currentTestCase.TestCase?.Title || currentTestCase.title }}</h5>
              <p class="text-muted">{{ currentTestCase.TestCase?.Description || currentTestCase.description }}</p>
              
              <!-- Test Steps -->
              <div v-if="currentTestSteps && currentTestSteps.length > 0" class="mt-3">
                <h6><i class="fas fa-list-ol"></i> Test Steps:</h6>
                <div class="test-steps">
                  <div 
                    v-for="step in currentTestSteps" 
                    :key="step.id"
                    class="step-item mb-3 p-3 border rounded"
                  >
                    <div class="fw-bold text-primary">Step {{ step.step_number }}</div>
                    <div class="step-description mb-2">{{ step.description }}</div>
                    <div class="step-expected">
                      <small class="text-muted fw-bold">Expected Result:</small>
                      <small class="d-block">{{ step.expected_result }}</small>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Result Selection -->
            <div class="row">
              <div class="col-md-8">
                <!-- Elapsed Time Display -->
                <div v-if="elapsedTime" class="mb-3">
                  <small class="text-muted">
                    <i class="fas fa-clock"></i> Elapsed Time: {{ elapsedTime }}
                  </small>
                </div>
                
                <!-- Result Selection -->
                <div class="mb-3">
                  <label class="form-label fw-bold">Test Result *</label>
                  <div class="btn-group d-flex" role="group">
                    <input type="radio" class="btn-check" :id="`pass-${currentTestCase.id}`" :value="'Pass'" v-model="currentResult.status">
                    <label class="btn btn-outline-success" :for="`pass-${currentTestCase.id}`">
                      <i class="fas fa-check"></i> Pass
                    </label>
                    
                    <input type="radio" class="btn-check" :id="`fail-${currentTestCase.id}`" :value="'Fail'" v-model="currentResult.status">
                    <label class="btn btn-outline-danger" :for="`fail-${currentTestCase.id}`">
                      <i class="fas fa-times"></i> Fail
                    </label>
                    
                    <input type="radio" class="btn-check" :id="`skip-${currentTestCase.id}`" :value="'Skip'" v-model="currentResult.status">
                    <label class="btn btn-outline-warning" :for="`skip-${currentTestCase.id}`">
                      <i class="fas fa-forward"></i> Skip
                    </label>
                  </div>
                </div>
              </div>
              <div class="col-md-4">
                <!-- Save Button -->
                <div class="d-grid">
                  <button 
                    @click="saveTestResult" 
                    class="btn btn-primary btn-lg"
                    :disabled="!currentResult.status || saving"
                  >
                    <i class="fas fa-save"></i> 
                    {{ saving ? 'Saving...' : 'Save Result' }}
                  </button>
                </div>
              </div>
            </div>
            
            <!-- Comments -->
            <div class="mb-3">
              <label class="form-label">Comments</label>
              <textarea 
                class="form-control" 
                rows="3" 
                placeholder="Add comments about test execution..."
                v-model="currentResult.notes"
              ></textarea>
            </div>
          </div>
        </div>
        
        <!-- Not Started State -->
        <div class="card" v-else-if="testRun?.status === 'Not Started'">
          <div class="card-body text-center py-5">
            <i class="fas fa-play-circle fa-4x text-muted mb-3"></i>
            <h5>Test Run Not Started</h5>
            <p class="text-muted">Click "Start" to begin test execution</p>
          </div>
        </div>
        
        <!-- Completed State -->
        <div class="card" v-else-if="testRun?.status === 'Completed'">
          <div class="card-body text-center py-5">
            <i class="fas fa-check-circle fa-4x text-success mb-3"></i>
            <h5>Test Run Completed</h5>
            <p class="text-muted">All tests have been executed</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass, getTestCaseStatusBadgeClass, calculateDuration } from '../utils/helpers.js'

export default {
  name: 'TestRunDetail',
  props: {
    id: {
      type: [String, Number],
      required: true
    }
  },
  data() {
    return {
      testRun: null,
      testCases: [],
      currentTestSteps: [], // Store current test case steps
      loading: false,
      saving: false,
      currentTestCaseIndex: 0,
      currentResult: {
        status: '',
        notes: ''
      },
      elapsedTime: null,
      elapsedTimer: null
    }
  },
  async mounted() {
    if (this.id) {
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
    getTestCaseStatusBadgeClass,
    calculateDuration,
    
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
          this.loadCurrentTestSteps()
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
    
    async loadCurrentTestSteps() {
      if (!this.currentTestCase?.TestCase?.ID) return
      
      try {
        const steps = await api.getTestCaseSteps(this.currentTestCase.TestCase.ID)
        this.currentTestSteps = steps
      } catch (error) {
        console.error('Error loading test steps:', error)
        this.currentTestSteps = []
      }
    },
    
    loadCurrentTestResult() {
      if (!this.currentTestCase) return
      
      this.currentResult = {
        status: this.currentTestCase.Status || '',
        notes: this.currentTestCase.ResultNotes || ''
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

    // Navigation Methods
    selectTestCase(index) {
      if (index >= 0 && index < this.testCases.length) {
        this.currentTestCaseIndex = index
        this.loadCurrentTestResult()
        this.loadCurrentTestSteps()
      }
    },
    
    nextTestCase() {
      if (this.currentTestCaseIndex < this.testCases.length - 1) {
        this.selectTestCase(this.currentTestCaseIndex + 1)
      }
    },
    
    previousTestCase() {
      if (this.currentTestCaseIndex > 0) {
        this.selectTestCase(this.currentTestCaseIndex - 1)
      }
    },

    // Test Result Management
    async saveTestResult() {
      if (!this.currentResult.status) {
        showAlert('Please select a test result status', 'warning')
        return
      }
      
      try {
        this.saving = true
        
        const request = {
          status: this.currentResult.status,
          result_notes: this.currentResult.notes || null,
          executed_by: 'Current User', // TODO: Get from auth
          completed_at: new Date().toISOString()
        }
        
        await api.updateTestRunCase(this.testRun.id, this.currentTestCase.TestCase.ID, request)
        
        // Update local data
        this.currentTestCase.Status = this.currentResult.status
        this.currentTestCase.ResultNotes = this.currentResult.notes
        this.currentTestCase.UpdatedAt = new Date().toISOString()
        
        showAlert('Test result saved successfully!', 'success')
        
        // Auto-advance to next test case if not the last one
        if (this.currentTestCaseIndex < this.testCases.length - 1) {
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

    // Timer Methods
    startElapsedTimer() {
      this.elapsedTimer = setInterval(() => {
        if (this.testRun?.started_at) {
          const start = new Date(this.testRun.started_at)
          const now = new Date()
          const diff = now - start
          
          const hours = Math.floor(diff / (1000 * 60 * 60))
          const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
          const seconds = Math.floor((diff % (1000 * 60)) / 1000)
          
          this.elapsedTime = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
        }
      }, 1000)
    },
    
    stopElapsedTimer() {
      if (this.elapsedTimer) {
        clearInterval(this.elapsedTimer)
        this.elapsedTimer = null
      }
    }
  },
  
  beforeUnmount() {
    this.stopElapsedTimer()
  }
}
</script>

<style scoped>
.test-steps {
  max-height: 400px;
  overflow-y: auto;
}

.step-item {
  background-color: #f8f9fa;
}

.step-item:hover {
  background-color: #e9ecef;
}

.list-group-item.active {
  background-color: #0d6efd;
  border-color: #0d6efd;
}

.test-case-list {
  max-height: 600px;
  overflow-y: auto;
}
</style>