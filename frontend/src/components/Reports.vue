<template>
  <div>
    <!-- Page Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Reports & Analytics</h1>
    </div>

    <!-- Overview Stats -->
    <div class="row mb-4" v-if="stats">
      <div class="col-md-3">
        <div class="card stats-card">
          <h3>{{ stats.totalProjects }}</h3>
          <p>Total Projects</p>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card" style="background: linear-gradient(135deg, #28a745, #20c997);">
          <h3>{{ stats.totalTestSuites }}</h3>
          <p>Test Suites</p>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card" style="background: linear-gradient(135deg, #ffc107, #fd7e14);">
          <h3>{{ stats.totalTestCases }}</h3>
          <p>Test Cases</p>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stats-card" style="background: linear-gradient(135deg, #17a2b8, #6f42c1);">
          <h3>{{ stats.totalTestRuns }}</h3>
          <p>Test Runs</p>
        </div>
      </div>
    </div>

    <!-- Test Execution Status -->
    <div class="row mb-4">
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">
            <h5><i class="fas fa-chart-pie"></i> Test Case Status Distribution</h5>
          </div>
          <div class="card-body">
            <div v-if="testCaseStats">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Pass</span>
                <span class="badge bg-success">{{ testCaseStats.pass }}</span>
              </div>
              <div class="progress mb-3">
                <div class="progress-bar bg-success" :style="{ width: getPercentage(testCaseStats.pass, testCaseStats.total) + '%' }"></div>
              </div>
              
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Fail</span>
                <span class="badge bg-danger">{{ testCaseStats.fail }}</span>
              </div>
              <div class="progress mb-3">
                <div class="progress-bar bg-danger" :style="{ width: getPercentage(testCaseStats.fail, testCaseStats.total) + '%' }"></div>
              </div>
              
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Blocked</span>
                <span class="badge bg-warning">{{ testCaseStats.blocked }}</span>
              </div>
              <div class="progress mb-3">
                <div class="progress-bar bg-warning" :style="{ width: getPercentage(testCaseStats.blocked, testCaseStats.total) + '%' }"></div>
              </div>
              
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Not Executed</span>
                <span class="badge bg-secondary">{{ testCaseStats.notExecuted }}</span>
              </div>
              <div class="progress">
                <div class="progress-bar bg-secondary" :style="{ width: getPercentage(testCaseStats.notExecuted, testCaseStats.total) + '%' }"></div>
              </div>
            </div>
            <div v-else v-html="showLoading()"></div>
          </div>
        </div>
      </div>
      
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">
            <h5><i class="fas fa-chart-bar"></i> Priority Distribution</h5>
          </div>
          <div class="card-body">
            <div v-if="priorityStats">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>High Priority</span>
                <span class="badge bg-danger">{{ priorityStats.high }}</span>
              </div>
              <div class="progress mb-3">
                <div class="progress-bar bg-danger" :style="{ width: getPercentage(priorityStats.high, priorityStats.total) + '%' }"></div>
              </div>
              
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Medium Priority</span>
                <span class="badge bg-warning">{{ priorityStats.medium }}</span>
              </div>
              <div class="progress mb-3">
                <div class="progress-bar bg-warning" :style="{ width: getPercentage(priorityStats.medium, priorityStats.total) + '%' }"></div>
              </div>
              
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Low Priority</span>
                <span class="badge bg-info">{{ priorityStats.low }}</span>
              </div>
              <div class="progress">
                <div class="progress-bar bg-info" :style="{ width: getPercentage(priorityStats.low, priorityStats.total) + '%' }"></div>
              </div>
            </div>
            <div v-else v-html="showLoading()"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="card">
      <div class="card-header">
        <h5><i class="fas fa-clock"></i> Recent Test Runs</h5>
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="recentTestRuns.length === 0" class="empty-state">
          <i class="fas fa-play-circle"></i>
          <h5>No Recent Test Runs</h5>
          <p>Test runs will appear here once you start executing tests.</p>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>Test Run</th>
                <th>Test Case</th>
                <th>Status</th>
                <th>Executed At</th>
                <th>Duration</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="run in recentTestRuns" :key="run.id">
                <td>
                  <router-link :to="`/test-runs/${run.id}`" class="text-decoration-none">
                    <strong>{{ run.name }}</strong>
                  </router-link>
                </td>
                <td>{{ run.test_case?.title || 'N/A' }}</td>
                <td>
                  <span class="status-badge" :class="getStatusBadgeClass(run.status)">
                    {{ run.status }}
                  </span>
                </td>
                <td>{{ formatDate(run.started_at) }}</td>
                <td>{{ calculateDuration(run.started_at, run.completed_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, getStatusBadgeClass, calculateStats } from '../utils/helpers.js'

export default {
  name: 'Reports',
  data() {
    return {
      stats: null,
      testCaseStats: null,
      priorityStats: null,
      recentTestRuns: [],
      loading: true
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    formatDate,
    showLoading,
    getStatusBadgeClass,
    
    async loadData() {
      this.loading = true
      try {
        const [statsData, testRunsData] = await Promise.all([
          api.getStats().catch(() => null),
          api.getTestRuns().catch(() => [])
        ])
        
        this.stats = statsData
        this.recentTestRuns = testRunsData.slice(0, 10) // Last 10 test runs
        
        // Calculate test case statistics
        await this.calculateTestCaseStats()
      } catch (error) {
        showAlert('Error loading reports data: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },

    async calculateTestCaseStats() {
      try {
        const testCases = await api.getTestCases()
        
        // Status distribution
        this.testCaseStats = calculateStats(testCases, 'status')
        
        // Priority distribution
        this.priorityStats = {
          total: testCases.length,
          high: testCases.filter(tc => tc.priority?.toUpperCase() === 'HIGH').length,
          medium: testCases.filter(tc => tc.priority?.toUpperCase() === 'MEDIUM').length,
          low: testCases.filter(tc => tc.priority?.toUpperCase() === 'LOW').length
        }
      } catch (error) {
        console.error('Error calculating test case stats:', error)
      }
    },

    getPercentage(value, total) {
      if (!total || total === 0) return 0
      return Math.round((value / total) * 100)
    },

    calculateDuration(startedAt, completedAt) {
      if (!startedAt || !completedAt) return 'N/A'
      
      const start = new Date(startedAt)
      const end = new Date(completedAt)
      const durationMs = end - start
      
      if (durationMs < 0) return 'N/A'
      
      const minutes = Math.floor(durationMs / 60000)
      const seconds = Math.floor((durationMs % 60000) / 1000)
      
      if (minutes > 0) {
        return `${minutes}m ${seconds}s`
      } else {
        return `${seconds}s`
      }
    }
  }
}
</script>