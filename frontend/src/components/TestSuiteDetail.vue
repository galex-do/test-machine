<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link :to="`/project/${pid}`">{{ truncateText(testSuite?.project?.name, 30) }}</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(testSuite?.name, 30) }}</li>
      </ol>
    </nav>

    <!-- Test Suite Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">{{ testSuite?.name }}</h1>
        <p class="text-muted">{{ testSuite?.description || 'No description available' }}</p>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <button type="button" class="btn btn-primary" @click="showCreateTestCaseModal">
          <i class="fas fa-plus"></i> New Test Case
        </button>
      </div>
    </div>

    <!-- Test Cases List -->
    <div class="card">
      <div class="card-header">
        <h5><i class="fas fa-list-check"></i> Test Cases</h5>
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="!testCases || testCases.length === 0" class="empty-state">
          <i class="fas fa-list-check"></i>
          <h5>No Test Cases Found</h5>
          <p>Create test cases to define your testing procedures.</p>
          <button class="btn btn-primary" @click="showCreateTestCaseModal">
            <i class="fas fa-plus"></i> Create Test Case
          </button>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>Title</th>
                <th>Priority</th>
                <th>Status</th>
                <th>Steps</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="testCase in testCases" :key="testCase.id">
                <td>
                  <router-link :to="`/project/${pid}/test-suite/${sid}/test-case/${testCase.id}`" class="text-decoration-none">
                    <strong>{{ testCase.title }}</strong>
                  </router-link>
                  <br>
                  <small class="text-muted">{{ testCase.description || 'No description' }}</small>
                </td>
                <td>
                  <span class="priority-badge" :class="getPriorityBadgeClass(testCase.priority)">
                    {{ testCase.priority }}
                  </span>
                </td>
                <td>
                  <span class="status-badge" :class="getStatusBadgeClass(testCase.status)">
                    {{ testCase.status }}
                  </span>
                </td>
                <td>{{ testCase.test_steps_count || 0 }}</td>
                <td>
                  <div class="action-buttons">
                    <button class="btn btn-outline-primary btn-sm" @click="editTestCase(testCase)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                    <button class="btn btn-outline-danger btn-sm" @click="deleteTestCase(testCase.id)">
                      <i class="fas fa-trash"></i> Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Test Case Modal -->
    <TestCaseModal 
      :show="showModal" 
      :testCase="selectedTestCase"
      :testSuiteId="parseInt(sid)"
      @close="closeModal"
      @saved="handleTestCaseSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass } from '../utils/helpers.js'
import TestCaseModal from './modals/TestCaseModal.vue'

export default {
  name: 'TestSuiteDetail',
  components: {
    TestCaseModal
  },
  props: {
    pid: {
      type: String,
      required: true
    },
    sid: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      testSuite: null,
      testCases: [],
      loading: true,
      showModal: false,
      selectedTestCase: null
    }
  },
  mounted() {
    this.loadData()
  },
  watch: {
    sid() {
      this.loadData()
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
        const [testSuiteData, testCasesData] = await Promise.all([
          api.getTestSuite(this.sid),
          api.getTestCases(this.sid)
        ])
        this.testSuite = testSuiteData
        this.testCases = Array.isArray(testCasesData) ? testCasesData : []
      } catch (error) {
        showAlert('Error loading test suite data: ' + error.message, 'danger')
        this.testCases = [] // Ensure testCases is always an array
      } finally {
        this.loading = false
      }
    },

    showCreateTestCaseModal() {
      this.selectedTestCase = null
      this.showModal = true
    },

    editTestCase(testCase) {
      this.selectedTestCase = testCase
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.selectedTestCase = null
    },

    handleTestCaseSaved() {
      this.closeModal()
      this.loadData()
      showAlert('Test case saved successfully!', 'success')
    },

    async deleteTestCase(testCaseId) {
      if (!confirm('Are you sure you want to delete this test case? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteTestCase(testCaseId)
        showAlert('Test case deleted successfully!', 'success')
        this.loadData()
      } catch (error) {
        showAlert('Error deleting test case: ' + error.message, 'danger')
      }
    }
  }
}
</script>