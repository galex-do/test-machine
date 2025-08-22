<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(testRun?.name, 30) }}</li>
      </ol>
    </nav>

    <!-- Test Run Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">{{ testRun?.name }}</h1>
        <p class="text-muted">{{ testRun?.description || 'No description available' }}</p>
        <div class="mt-2" v-if="testRun">
          <span class="status-badge me-2" :class="getStatusBadgeClass(testRun.status)">
            {{ testRun.status }}
          </span>
          <small class="text-muted">
            Created: {{ formatDate(testRun.created_at) }}
          </small>
        </div>
      </div>
    </div>

    <!-- Test Run Details -->
    <div class="row mb-4" v-if="testRun">
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Status</h5>
            <span class="status-badge" :class="getStatusBadgeClass(testRun.status)">
              {{ testRun.status }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Result</h5>
            <p>{{ testRun.result || 'N/A' }}</p>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Executed By</h5>
            <p>{{ testRun.executed_by || 'N/A' }}</p>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Duration</h5>
            <p>{{ calculateDuration(testRun.started_at, testRun.completed_at) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Run Information -->
    <div class="card" v-if="testRun">
      <div class="card-header">
        <h5><i class="fas fa-info-circle"></i> Test Run Information</h5>
      </div>
      <div class="card-body">
        <div class="row">
          <div class="col-md-6">
            <p><strong>Test Case:</strong> 
              <router-link 
                v-if="testRun.test_case" 
                :to="`/project/${testRun.test_case.test_suite?.project_id}/test-suite/${testRun.test_case.test_suite_id}/test-case/${testRun.test_case.id}`"
                class="text-decoration-none"
              >
                {{ testRun.test_case.title }}
              </router-link>
              <span v-else>N/A</span>
            </p>
            <p><strong>Started At:</strong> {{ formatDate(testRun.started_at) }}</p>
            <p><strong>Completed At:</strong> {{ formatDate(testRun.completed_at) }}</p>
          </div>
          <div class="col-md-6">
            <p><strong>Notes:</strong></p>
            <div class="border rounded p-3 bg-light">
              {{ testRun.notes || 'No notes available' }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" v-html="showLoading()"></div>
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass } from '../utils/helpers.js'

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
      loading: true
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
  methods: {
    formatDate,
    showLoading,
    truncateText,
    getStatusBadgeClass,
    
    async loadData() {
      this.loading = true
      try {
        this.testRun = await api.getTestRun(this.id)
      } catch (error) {
        showAlert('Error loading test run data: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
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