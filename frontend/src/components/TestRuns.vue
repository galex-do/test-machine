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
                  <th>Status</th>
                  <th>Progress</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="testRun in filteredTestRuns" :key="testRun.id" @click="viewTestRun(testRun.id)" style="cursor: pointer;">
                  <td>
                    <strong>{{ testRun.name }}</strong>
                    <small v-if="testRun.description" class="text-muted d-block">{{ testRun.description }}</small>
                  </td>
                  <td>
                    <span class="badge bg-info">{{ testRun.project?.name || 'Unknown' }}</span>
                  </td>
                  <td>
                    <span v-if="testRun.branch_name" class="badge bg-primary">
                      <i class="fas fa-code-branch"></i> {{ testRun.branch_name }}
                    </span>
                    <span v-else-if="testRun.tag_name" class="badge bg-warning">
                      <i class="fas fa-tag"></i> {{ testRun.tag_name }}
                    </span>
                    <span v-else class="text-muted">-</span>
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
                      <button 
                        @click.stop="viewTestRun(testRun.id)" 
                        class="btn btn-outline-primary"
                        title="View Details"
                      >
                        <i class="fas fa-eye"></i>
                      </button>
                      <router-link 
                        :to="`/test-runs/${testRun.id}/edit`"
                        class="btn btn-outline-secondary"
                        title="Edit"
                        @click.stop
                      >
                        <i class="fas fa-edit"></i>
                      </router-link>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import api from '../services/api.js'
import { formatDate, showAlert } from '../utils/helpers.js'

export default {
  name: 'TestRuns',
  data() {
    return {
      testRuns: [],
      loading: true,
      error: null,
      searchQuery: '',
      statusFilter: ''
    }
  },
  computed: {
    filteredTestRuns() {
      if (!this.testRuns || !Array.isArray(this.testRuns)) {
        return []
      }

      let filtered = [...this.testRuns]

      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(run => 
          (run.name || '').toLowerCase().includes(query) ||
          (run.description || '').toLowerCase().includes(query) ||
          (run.project?.name || '').toLowerCase().includes(query)
        )
      }

      if (this.statusFilter) {
        filtered = filtered.filter(run => run.status === this.statusFilter)
      }

      return filtered
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
        this.testRuns = await api.getTestRuns()
      } catch (error) {
        this.error = 'Error loading test runs: ' + error.message
      } finally {
        this.loading = false
      }
    },

    viewTestRun(id) {
      this.$router.push(`/test-run/${id}`)
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