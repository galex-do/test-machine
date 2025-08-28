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
        <div class="btn-group me-2">
          <button type="button" class="btn btn-outline-secondary" @click="showEditTestSuiteModal">
            <i class="fas fa-edit"></i> Edit Test Suite
          </button>
        </div>
        <button type="button" class="btn btn-primary" @click="showCreateTestCaseModal">
          <i class="fas fa-plus"></i> New Test Case
        </button>
      </div>
    </div>

    <!-- Test Cases List -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5><i class="fas fa-list-check"></i> Test Cases</h5>
        <SortBy 
          :sortOptions="testCaseSortOptions"
          :defaultSort="sortBy"
          componentId="test-cases"
          @sort-changed="handleSortChange"
        />
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
                <th>Created</th>
                <th>Updated</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="testCase in sortedTestCases" :key="testCase.id">
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
                <td>{{ formatDate(testCase.created_at) }}</td>
                <td>{{ formatDate(testCase.updated_at) }}</td>
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
      
      <!-- Pagination -->
      <div class="card-footer" v-if="!loading && pagination.total > 0">
        <Pagination 
          :pagination="pagination"
          @page-changed="changePage"
          @page-size-changed="changePageSize"
        />
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

    <!-- Test Suite Modal -->
    <TestSuiteModal 
      :show="showTestSuiteModal" 
      :testSuite="testSuite"
      :projectId="parseInt(pid)"
      @close="closeTestSuiteModal"
      @saved="handleTestSuiteSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass } from '../utils/helpers.js'
import { applySorting, SORT_OPTION_SETS } from '../utils/sortUtils.js'
import TestCaseModal from './modals/TestCaseModal.vue'
import TestSuiteModal from './modals/TestSuiteModal.vue'
import Pagination from './Pagination.vue'
import SortBy from './SortBy.vue'

export default {
  name: 'TestSuiteDetail',
  components: {
    TestCaseModal,
    TestSuiteModal,
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
    }
  },
  data() {
    return {
      testSuite: null,
      testCases: [],
      loading: true,
      showModal: false,
      selectedTestCase: null,
      showTestSuiteModal: false,
      sortBy: 'created_desc',
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      },
      allTestCases: [], // Store all test cases for sorting and pagination
      testCaseSortOptions: SORT_OPTION_SETS.TEST_CASES
    }
  },
  computed: {
    sortedTestCases() {
      // For paginated results, we return the current page data already sorted
      return this.testCases || []
    }
  },
  watch: {
    sid() {
      this.loadData()
    }
  },
  mounted() {
    this.loadData()
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
        this.allTestCases = Array.isArray(testCasesData) ? testCasesData : []
        this.applyCurrentSorting()
      } catch (error) {
        showAlert('Error loading test suite data: ' + error.message, 'danger')
        this.allTestCases = []
        this.testCases = []
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
      if (!this.allTestCases || !Array.isArray(this.allTestCases)) {
        this.testCases = []
        return
      }

      // Use the utility function to sort the data
      const sortedCases = applySorting(this.allTestCases, this.sortBy)

      // Apply pagination
      const startIndex = (this.pagination.page - 1) * this.pagination.page_size
      const endIndex = startIndex + this.pagination.page_size
      
      // Update pagination info
      this.pagination.total = sortedCases.length
      this.pagination.total_pages = Math.ceil(sortedCases.length / this.pagination.page_size)
      this.pagination.has_next = this.pagination.page < this.pagination.total_pages
      this.pagination.has_prev = this.pagination.page > 1
      
      // Get current page data
      this.testCases = sortedCases.slice(startIndex, endIndex)
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
    },

    showEditTestSuiteModal() {
      this.showTestSuiteModal = true
    },

    closeTestSuiteModal() {
      this.showTestSuiteModal = false
    },

    handleTestSuiteSaved() {
      this.closeTestSuiteModal()
      this.loadData()
      showAlert('Test suite updated successfully!', 'success')
    }
  }
}
</script>