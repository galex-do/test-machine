<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link :to="`/project/${pid}`">{{ truncateText(testCase?.test_suite?.project?.name, 30) }}</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link :to="`/project/${pid}/test-suite/${sid}`">{{ truncateText(testCase?.test_suite?.name, 30) }}</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(testCase?.title, 30) }}</li>
      </ol>
    </nav>

    <!-- Test Case Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">{{ testCase?.title }}</h1>
        <p class="text-muted">{{ testCase?.description || 'No description available' }}</p>
        <div class="mt-2">
          <span class="priority-badge me-2" :class="getPriorityBadgeClass(testCase?.priority)">
            {{ testCase?.priority }}
          </span>
          <span class="status-badge" :class="getStatusBadgeClass(testCase?.status)">
            {{ testCase?.status }}
          </span>
        </div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-outline-secondary" @click="showEditTestCaseModal">
            <i class="fas fa-edit"></i> Edit Test Case
          </button>
        </div>
        <button type="button" class="btn btn-primary" @click="showCreateTestStepModal">
          <i class="fas fa-plus"></i> Add Test Step
        </button>
      </div>
    </div>

    <!-- Test Case Info -->
    <div class="row mb-4" v-if="testCase">
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Priority</h5>
            <span class="priority-badge" :class="getPriorityBadgeClass(testCase.priority)">
              {{ testCase.priority }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Status</h5>
            <span class="status-badge" :class="getStatusBadgeClass(testCase.status)">
              {{ testCase.status }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Test Steps</h5>
            <h3>{{ testSteps?.length || 0 }}</h3>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Created</h5>
            <small>{{ formatDate(testCase.created_at) }}</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Steps -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5><i class="fas fa-list-ol"></i> Test Steps</h5>
        <SortBy 
          :sortOptions="testStepSortOptions"
          :defaultSort="sortBy"
          componentId="test-steps"
          @sort-changed="handleSortChange"
        />
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="!testSteps || testSteps.length === 0" class="empty-state">
          <i class="fas fa-list-ol"></i>
          <h5>No Test Steps Found</h5>
          <p>Add test steps to define the testing procedure.</p>
          <button class="btn btn-primary" @click="showCreateTestStepModal">
            <i class="fas fa-plus"></i> Add Test Step
          </button>
        </div>
        <div v-else>
          <div v-for="step in sortedTestSteps" :key="step.id" class="test-step">
            <div class="d-flex align-items-start">
              <div class="step-number">{{ step.step_number }}</div>
              <div class="flex-grow-1">
                <div class="mb-2">
                  <strong>Description:</strong><br>
                  {{ step.description }}
                </div>
                <div class="mb-2">
                  <strong>Expected Result:</strong><br>
                  {{ step.expected_result }}
                </div>
                <div class="action-buttons">
                  <button class="btn btn-outline-secondary btn-sm" @click="editTestStep(step)">
                    <i class="fas fa-edit"></i> Edit
                  </button>
                  <button class="btn btn-outline-danger btn-sm" @click="deleteTestStep(step.id)">
                    <i class="fas fa-trash"></i> Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Pagination -->
      <div class="card-footer" v-if="!loading && pagination.total > 0">
        <Pagination 
          :pagination="pagination"
          @page-changed="changePage"
          @page-size-changed="changePageSize"
        />
      </div>
    </div>

    <!-- Test Step Modal -->
    <TestStepModal 
      :show="showModal" 
      :testStep="selectedTestStep"
      :testCaseId="parseInt(cid)"
      @close="closeModal"
      @saved="handleTestStepSaved"
    />

    <!-- Test Case Modal -->
    <TestCaseModal 
      :show="showTestCaseModal" 
      :testCase="testCase"
      :testSuiteId="parseInt(sid)"
      @close="closeTestCaseModal"
      @saved="handleTestCaseSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass } from '../utils/helpers.js'
import { applySorting } from '../utils/sortUtils.js'
import TestStepModal from './modals/TestStepModal.vue'
import TestCaseModal from './modals/TestCaseModal.vue'
import Pagination from './Pagination.vue'
import SortBy from './SortBy.vue'

export default {
  name: 'TestCaseDetail',
  components: {
    TestStepModal,
    TestCaseModal,
    Pagination,
    SortBy
  },
  props: {
    pid: {
      type: String,
      required: true
    },
    sid: {
      type: String,
      required: true
    },
    cid: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      testCase: null,
      testSteps: [],
      loading: true,
      showModal: false,
      selectedTestStep: null,
      showTestCaseModal: false,
      sortBy: 'step_number_asc',
      testStepSortOptions: [
        { value: 'step_number_asc', label: 'Step Number (1-10)' },
        { value: 'step_number_desc', label: 'Step Number (10-1)' },
        { value: 'created_desc', label: 'Created Date (Newest First)' },
        { value: 'created_asc', label: 'Created Date (Oldest First)' }
      ],
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      },
      allTestSteps: [] // Store all test steps for pagination
    }
  },
  computed: {
    sortedTestSteps() {
      return this.testSteps ? [...this.testSteps].sort((a, b) => a.step_number - b.step_number) : []
    }
  },
  mounted() {
    this.loadData()
  },
  watch: {
    cid() {
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
        const [testCaseData, testStepsData] = await Promise.all([
          api.getTestCase(this.cid),
          api.getTestSteps(this.cid)
        ])
        this.testCase = testCaseData
        
        const stepsArray = Array.isArray(testStepsData) ? testStepsData : []
        this.allTestSteps = stepsArray
        this.applyCurrentSorting()
      } catch (error) {
        showAlert('Error loading test case data: ' + error.message, 'danger')
        this.allTestSteps = []
        this.testSteps = []
      } finally {
        this.loading = false
      }
    },

    handleSortChange(newSortBy) {
      this.sortBy = newSortBy
      this.pagination.page = 1 // Reset to first page when sorting changes
      this.applyCurrentSorting()
    },

    applyCurrentSorting() {
      if (!this.allTestSteps || !Array.isArray(this.allTestSteps)) {
        this.testSteps = []
        return
      }

      let sortedSteps = [...this.allTestSteps]
      
      // Custom sorting for test steps
      switch (this.sortBy) {
        case 'step_number_asc':
          sortedSteps = sortedSteps.sort((a, b) => a.step_number - b.step_number)
          break
        case 'step_number_desc':
          sortedSteps = sortedSteps.sort((a, b) => b.step_number - a.step_number)
          break
        default:
          // Use utility function for other sorts
          sortedSteps = applySorting(this.allTestSteps, this.sortBy)
      }

      // Apply pagination
      const startIndex = (this.pagination.page - 1) * this.pagination.page_size
      const endIndex = startIndex + this.pagination.page_size
      
      // Update pagination info
      this.pagination.total = sortedSteps.length
      this.pagination.total_pages = Math.ceil(sortedSteps.length / this.pagination.page_size)
      this.pagination.has_next = this.pagination.page < this.pagination.total_pages
      this.pagination.has_prev = this.pagination.page > 1
      
      // Get current page data
      this.testSteps = sortedSteps.slice(startIndex, endIndex)
    },

    applyPagination() {
      // Deprecated - use applyCurrentSorting instead
      this.applyCurrentSorting()
    },

    changePage(page) {
      this.pagination.page = page
      this.applyCurrentSorting()
    },

    changePageSize(pageSize) {
      this.pagination.page_size = pageSize
      this.pagination.page = 1 // Reset to first page
      this.applyCurrentSorting()
    },

    showCreateTestStepModal() {
      this.selectedTestStep = null
      this.showModal = true
    },

    editTestStep(testStep) {
      this.selectedTestStep = testStep
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.selectedTestStep = null
    },

    handleTestStepSaved() {
      this.closeModal()
      this.loadData()
      showAlert('Test step saved successfully!', 'success')
    },

    async deleteTestStep(testStepId) {
      if (!confirm('Are you sure you want to delete this test step? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteTestStep(testStepId)
        showAlert('Test step deleted successfully!', 'success')
        this.loadData()
      } catch (error) {
        showAlert('Error deleting test step: ' + error.message, 'danger')
      }
    },

    showEditTestCaseModal() {
      this.showTestCaseModal = true
    },

    closeTestCaseModal() {
      this.showTestCaseModal = false
    },

    handleTestCaseSaved() {
      this.closeTestCaseModal()
      this.loadData()
      showAlert('Test case updated successfully!', 'success')
    }
  }
}
</script>