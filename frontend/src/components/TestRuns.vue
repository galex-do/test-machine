<template>
  <div class="test-runs-container">
    <!-- Header -->
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>
        <i class="fas fa-play-circle text-primary"></i> Test Runs
      </h2>
      <router-link 
        to="/test-runs/new" 
        class="btn btn-primary"
      >
        <i class="fas fa-plus"></i> New Test Run
      </router-link>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading test runs...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="alert alert-danger" role="alert">
      <i class="fas fa-exclamation-triangle"></i> {{ error }}
    </div>

    <!-- Test Runs List -->
    <div v-if="!loading && !error">
      <div class="row mb-3">
        <div class="col-md-6">
          <input
            type="text"
            class="form-control"
            placeholder="Search test runs..."
            v-model="searchQuery"
          >
        </div>
        <div class="col-md-3">
          <select class="form-select" v-model="statusFilter">
            <option value="">All Statuses</option>
            <option value="Not Started">Not Started</option>
            <option value="In Progress">In Progress</option>
            <option value="Completed">Completed</option>
            <option value="Cancelled">Cancelled</option>
          </select>
        </div>
      </div>

      <!-- Test Runs Table -->
      <div class="card">
        <div class="card-body">
          <div v-if="filteredTestRuns.length === 0" class="text-center py-4 text-muted">
            <i class="fas fa-inbox fa-3x mb-3"></i>
            <h5>No test runs found</h5>
            <p>Create your first test run to start testing your projects.</p>
            <router-link to="/test-runs/new" class="btn btn-primary mt-2">
              <i class="fas fa-plus"></i> Create Test Run
            </router-link>
          </div>

          <div v-else class="table-responsive">
            <table class="table table-hover">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Project</th>
                  <th>Git Reference</th>
                  <th>Test Cases</th>
                  <th>Status</th>
                  <th>Progress</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="testRun in filteredTestRuns" :key="testRun.id">
                  <td>
                    <router-link 
                      :to="`/test-run/${testRun.id}`" 
                      class="text-decoration-none fw-bold text-primary"
                    >
                      {{ testRun.name }}
                    </router-link>
                    <small v-if="testRun.description" class="text-muted d-block">{{ testRun.description }}</small>
                  </td>
                  <td>
                    <span class="badge bg-info">{{ testRun.project?.name || 'Unknown' }}</span>
                  </td>
                  <td>
                    <a 
                      v-if="testRun.branch_name && testRun.repository?.remote_url" 
                      :href="getGitReferenceUrl(testRun.repository.remote_url, testRun.branch_name)"
                      target="_blank"
                      class="badge bg-primary text-decoration-none"
                      :title="`View ${testRun.branch_name} branch in repository`"
                    >
                      <i class="fas fa-code-branch"></i> {{ testRun.branch_name }}
                      <i class="fas fa-external-link-alt ms-1" style="font-size: 0.75em;"></i>
                    </a>
                    <span v-else-if="testRun.branch_name" class="badge bg-primary">
                      <i class="fas fa-code-branch"></i> {{ testRun.branch_name }}
                    </span>
                    <a 
                      v-else-if="testRun.tag_name && testRun.repository?.remote_url"
                      :href="getGitReferenceUrl(testRun.repository.remote_url, testRun.tag_name)"
                      target="_blank"
                      class="badge bg-warning text-decoration-none"
                      :title="`View ${testRun.tag_name} tag in repository`"
                    >
                      <i class="fas fa-tag"></i> {{ testRun.tag_name }}
                      <i class="fas fa-external-link-alt ms-1" style="font-size: 0.75em;"></i>
                    </a>
                    <span v-else-if="testRun.tag_name" class="badge bg-warning">
                      <i class="fas fa-tag"></i> {{ testRun.tag_name }}
                    </span>
                    <span v-else class="text-muted">-</span>
                  </td>
                  <td>
                    <span class="badge bg-light text-dark">
                      <i class="fas fa-list-check"></i> {{ testRun.test_cases_count || 0 }}
                    </span>
                  </td>
                  <td>
                    <span :class="getStatusBadgeClass(testRun.status)">
                      {{ testRun.status }}
                    </span>
                  </td>
                  <td>
                    <div class="progress" style="height: 20px;">
                      <div 
                        class="progress-bar" 
                        :class="getProgressBarClass(testRun)"
                        :style="{ width: getProgressPercentage(testRun) + '%' }"
                      >
                        {{ getProgressText(testRun) }}
                      </div>
                    </div>
                  </td>
                  <td class="small text-muted">
                    {{ formatDate(testRun.created_at) }}
                  </td>
                  <td>
                    <div class="btn-group btn-group-sm" role="group">
                      <!-- Time Management Buttons -->
                      <button 
                        v-if="testRun.status === 'Not Started'"
                        @click="startTestRun(testRun.id)"
                        class="btn btn-outline-success"
                        title="Start Test Run"
                      >
                        <i class="fas fa-play"></i>
                      </button>
                      <button 
                        v-else-if="testRun.status === 'In Progress'"
                        @click="pauseTestRun(testRun.id)"
                        class="btn btn-outline-warning"
                        title="Pause Test Run"
                      >
                        <i class="fas fa-pause"></i>
                      </button>
                      <button 
                        v-if="testRun.status === 'In Progress'"
                        @click="finishTestRun(testRun.id)"
                        class="btn btn-outline-primary"
                        title="Finish Test Run"
                      >
                        <i class="fas fa-stop"></i>
                      </button>
                      
                      <!-- Management Buttons -->
                      <router-link 
                        v-if="testRun.status !== 'Completed'"
                        :to="`/test-runs/${testRun.id}/edit`"
                        class="btn btn-outline-secondary"
                        title="Edit"
                      >
                        <i class="fas fa-edit"></i>
                      </router-link>
                      <button 
                        v-if="testRun.status === 'Not Started'"
                        @click="deleteTestRun(testRun.id, testRun.name)"
                        class="btn btn-outline-danger"
                        title="Delete Test Run"
                      >
                        <i class="fas fa-trash"></i>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
        <!-- Pagination -->
        <Pagination 
          v-if="!loading && !error && pagination.total > 0"
          :pagination="pagination"
          @page-changed="changePage"
          @page-size-changed="changePageSize"
        />
      </div>
    </div>

  </div>
</template>

<script>
import api from '../services/api.js'
import { formatDate, showAlert } from '../utils/helpers.js'
import Pagination from './Pagination.vue'

export default {
  name: 'TestRuns',
  components: {
    Pagination
  },
  data() {
    return {
      testRuns: [],
      loading: true,
      error: null,
      searchQuery: '',
      statusFilter: '',
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      },
      allTestRuns: [] // Store all test runs for client-side filtering
    }
  },
  computed: {
    filteredTestRuns() {
      // For paginated results, we return the current page data  
      // Filtering is now handled server-side via the pagination API
      return this.testRuns || []
    }
  },
  watch: {
    searchQuery() {
      // Reset to first page when searching
      this.pagination.page = 1
      this.loadTestRuns()
    },
    statusFilter() {
      // Reset to first page when filtering
      this.pagination.page = 1
      this.loadTestRuns()
    }
  },
  mounted() {
    this.loadTestRuns()
  },
  methods: {
    formatDate,

    async loadTestRuns() {
      this.loading = true
      this.error = null
      try {
        // For now, load all test runs and handle pagination client-side
        // TODO: Update API to support server-side pagination
        const allTestRuns = await api.getTestRuns()
        
        // Simulate pagination for now
        const startIndex = (this.pagination.page - 1) * this.pagination.page_size
        const endIndex = startIndex + this.pagination.page_size
        
        // Apply filters first
        let filteredRuns = allTestRuns
        
        if (this.searchQuery) {
          const query = this.searchQuery.toLowerCase()
          filteredRuns = filteredRuns.filter(run => 
            (run.name || '').toLowerCase().includes(query) ||
            (run.description || '').toLowerCase().includes(query) ||
            (run.project?.name || '').toLowerCase().includes(query)
          )
        }

        if (this.statusFilter) {
          filteredRuns = filteredRuns.filter(run => run.status === this.statusFilter)
        }
        
        // Update pagination info
        this.pagination.total = filteredRuns.length
        this.pagination.total_pages = Math.ceil(filteredRuns.length / this.pagination.page_size)
        this.pagination.has_next = this.pagination.page < this.pagination.total_pages
        this.pagination.has_prev = this.pagination.page > 1
        
        // Get current page data
        this.testRuns = filteredRuns.slice(startIndex, endIndex)
        
      } catch (error) {
        this.error = 'Error loading test runs: ' + error.message
      } finally {
        this.loading = false
      }
    },

    changePage(page) {
      this.pagination.page = page
      this.loadTestRuns()
    },

    changePageSize(pageSize) {
      this.pagination.page_size = pageSize
      this.pagination.page = 1 // Reset to first page
      this.loadTestRuns()
    },


    getStatusBadgeClass(status) {
      const classes = {
        'Not Started': 'badge bg-secondary',
        'In Progress': 'badge bg-warning',
        'Completed': 'badge bg-success',
        'Cancelled': 'badge bg-danger'
      }
      return classes[status] || 'badge bg-secondary'
    },

    getProgressBarClass(testRun) {
      if (!testRun || !testRun.test_cases || !Array.isArray(testRun.test_cases) || testRun.test_cases.length === 0) return 'bg-secondary'
      
      const total = testRun.test_cases.length
      const passed = testRun.test_cases.filter(tc => tc && tc.status === 'Pass').length
      const failed = testRun.test_cases.filter(tc => tc && tc.status === 'Fail').length
      
      if (failed > 0) return 'bg-danger'
      if (passed === total) return 'bg-success'
      return 'bg-warning'
    },

    getProgressPercentage(testRun) {
      if (!testRun || !testRun.test_cases || !Array.isArray(testRun.test_cases) || testRun.test_cases.length === 0) return 0
      
      const total = testRun.test_cases.length
      const executed = testRun.test_cases.filter(tc => tc && tc.status && tc.status !== 'Not Executed').length
      
      return Math.round((executed / total) * 100)
    },

    getProgressText(testRun) {
      if (!testRun || !testRun.test_cases || !Array.isArray(testRun.test_cases) || testRun.test_cases.length === 0) return '0/0'
      
      const total = testRun.test_cases.length
      const executed = testRun.test_cases.filter(tc => tc && tc.status && tc.status !== 'Not Executed').length
      
      return `${executed}/${total}`
    },

    async deleteTestRun(id, name) {
      const confirmed = confirm(`Are you sure you want to delete the test run "${name}"? This action cannot be undone.`)
      
      if (confirmed) {
        try {
          await api.deleteTestRun(id)
          showAlert('Test run deleted successfully', 'success')
          await this.loadTestRuns() // Refresh the list
        } catch (error) {
          console.error('Error deleting test run:', error)
          showAlert('Failed to delete test run. Please try again.', 'error')
        }
      }
    },

    async startTestRun(id) {
      try {
        await api.startTestRun(id)
        showAlert('Test run started successfully', 'success')
        await this.loadTestRuns() // Refresh the list
      } catch (error) {
        console.error('Error starting test run:', error)
        showAlert('Failed to start test run: ' + error.message, 'error')
      }
    },

    async pauseTestRun(id) {
      try {
        await api.pauseTestRun(id)
        showAlert('Test run paused successfully', 'success')
        await this.loadTestRuns() // Refresh the list
      } catch (error) {
        console.error('Error pausing test run:', error)
        showAlert('Failed to pause test run: ' + error.message, 'error')
      }
    },

    async finishTestRun(id) {
      const confirmed = confirm('Are you sure you want to finish this test run? This will mark it as completed.')
      
      if (confirmed) {
        try {
          await api.finishTestRun(id)
          showAlert('Test run finished successfully', 'success')
          await this.loadTestRuns() // Refresh the list
        } catch (error) {
          console.error('Error finishing test run:', error)
          showAlert('Failed to finish test run: ' + error.message, 'error')
        }
      }
    },

    getGitReferenceUrl(remoteUrl, reference) {
      if (!remoteUrl || !reference) return '#'
      
      // Handle different Git hosting providers
      let baseUrl = remoteUrl
      
      // Convert SSH to HTTPS if needed
      if (baseUrl.startsWith('git@')) {
        // git@github.com:user/repo.git -> https://github.com/user/repo
        baseUrl = baseUrl
          .replace(/^git@/, 'https://')
          .replace(':', '/')
          .replace(/\.git$/, '')
      }
      
      // Remove .git extension if present
      if (baseUrl.endsWith('.git')) {
        baseUrl = baseUrl.slice(0, -4)
      }
      
      // Add tree path for branch/tag
      return `${baseUrl}/tree/${reference}`
    }
  }
}
</script>

<style scoped>
.test-runs-container {
  padding: 20px;
}

.table tbody tr:hover {
  background-color: #f8f9fa;
}

.progress {
  min-width: 80px;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
}
</style>