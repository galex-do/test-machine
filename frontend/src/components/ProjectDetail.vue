<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(project?.name, 30) }}</li>
      </ol>
    </nav>

    <!-- Project Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">
          {{ project?.name }}
          <a v-if="project?.repository" 
             :href="project.repository.remote_url" 
             target="_blank" 
             rel="noopener noreferrer" 
             class="ms-2 text-muted" 
             :title="`Open Git repository: ${project.repository.remote_url}`">
            <i class="fab fa-git-alt"></i>
          </a>
        </h1>
        <p class="text-muted">{{ project?.description || 'No description available' }}</p>
        <p v-if="project?.repository" class="text-muted">
          <i class="fab fa-git-alt"></i> <strong>Git Repository:</strong> 
          <a :href="project.repository.remote_url" target="_blank" rel="noopener noreferrer" class="text-decoration-none">
            {{ project.repository.remote_url }} <i class="fas fa-external-link-alt ms-1"></i>
          </a>
          <span v-if="project.repository.key" class="ms-3">
            <i class="fas fa-key"></i> <strong>Auth:</strong> {{ project.repository.key.name }} ({{ project.repository.key.key_type }})
          </span>
        </p>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-outline-secondary" @click="editProject">
            <i class="fas fa-edit"></i> Edit Project
          </button>
          <button 
            v-if="project?.repository" 
            type="button" 
            class="btn btn-outline-info" 
            @click="syncProject"
            :disabled="syncing">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': syncing }"></i> 
            {{ syncing ? 'Syncing...' : 'Sync Git' }}
          </button>
        </div>
        <button type="button" class="btn btn-primary" @click="showCreateTestSuiteModal">
          <i class="fas fa-plus"></i> New Test Suite
        </button>
      </div>
    </div>

    <!-- Project Stats -->
    <div class="row mb-4" v-if="project">
      <div class="col-md-4">
        <div class="card stats-card">
          <h3>{{ project.test_suites_count || 0 }}</h3>
          <p>Test Suites</p>
        </div>
      </div>
      <div class="col-md-4">
        <div class="card stats-card" style="background: linear-gradient(135deg, #28a745, #20c997);">
          <h3>{{ totalTestCases }}</h3>
          <p>Total Test Cases</p>
        </div>
      </div>
      <div class="col-md-4">
        <div class="card stats-card" style="background: linear-gradient(135deg, #ffc107, #fd7e14);">
          <h3>{{ formatDate(project.created_at) }}</h3>
          <p>Created</p>
        </div>
      </div>
    </div>

    <!-- Test Suites List -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5><i class="fas fa-layer-group"></i> Test Suites</h5>
        <SortBy 
          :sortOptions="testSuiteSortOptions"
          :defaultSort="sortBy"
          componentId="test-suites"
          @sort-changed="handleSortChange"
        />
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="testSuites.length === 0" class="empty-state">
          <i class="fas fa-layer-group"></i>
          <h5>No Test Suites Found</h5>
          <p>Create test suites to organize your test cases.</p>
          <button class="btn btn-primary" @click="showCreateTestSuiteModal">
            <i class="fas fa-plus"></i> Create Test Suite
          </button>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Test Cases</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="testSuite in testSuites" :key="testSuite.id">
                <td>
                  <router-link :to="`/project/${id}/test-suite/${testSuite.id}`" class="text-decoration-none">
                    <strong>{{ testSuite.name }}</strong>
                  </router-link>
                </td>
                <td>{{ testSuite.description || 'No description' }}</td>
                <td>
                  <span class="badge bg-primary">{{ testSuite.test_cases_count || 0 }}</span>
                </td>
                <td>{{ formatDate(testSuite.created_at) }}</td>
                <td>
                  <div class="action-buttons">
                    <button class="btn btn-outline-primary btn-sm" @click="editTestSuite(testSuite)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                    <button class="btn btn-outline-danger btn-sm" @click="deleteTestSuite(testSuite.id)">
                      <i class="fas fa-trash"></i> Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
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

    <!-- Test Suite Modal -->
    <TestSuiteModal 
      :show="showTestSuiteModal" 
      :testSuite="selectedTestSuite"
      :projectId="parseInt(id)"
      @close="closeTestSuiteModal"
      @saved="handleTestSuiteSaved"
    />

    <!-- Project Modal -->
    <ProjectModal 
      :show="showProjectModal" 
      :project="project"
      @close="closeProjectModal"
      @saved="handleProjectSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText } from '../utils/helpers.js'
import { applySorting, SORT_OPTION_SETS } from '../utils/sortUtils.js'
import TestSuiteModal from './modals/TestSuiteModal.vue'
import ProjectModal from './modals/ProjectModal.vue'
import Pagination from './Pagination.vue'
import SortBy from './SortBy.vue'

export default {
  name: 'ProjectDetail',
  components: {
    TestSuiteModal,
    ProjectModal,
    Pagination,
    SortBy
  },
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      project: null,
      testSuites: [],
      loading: true,
      showTestSuiteModal: false,
      showProjectModal: false,
      selectedTestSuite: null,
      syncing: false,
      sortBy: 'created_desc',
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      },
      allTestSuites: [], // Store all test suites for sorting and pagination
      testSuiteSortOptions: SORT_OPTION_SETS.TEST_SUITES
    }
  },
  computed: {
    totalTestCases() {
      if (!this.testSuites || !Array.isArray(this.testSuites)) {
        return 0
      }
      return this.testSuites.reduce((total, suite) => total + (suite.test_cases_count || 0), 0)
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
    
    async loadData() {
      this.loading = true
      try {
        const [projectData, allTestSuites] = await Promise.all([
          api.getProject(this.id),
          api.getTestSuites(this.id)
        ])
        this.project = projectData
        
        // Ensure testSuites is always an array
        this.allTestSuites = Array.isArray(allTestSuites) ? allTestSuites : []
        this.applyCurrentSorting()
      } catch (error) {
        showAlert('Error loading project data: ' + error.message, 'danger')
        this.allTestSuites = []
        this.testSuites = [] // Ensure testSuites is always an array
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
      if (!this.allTestSuites || !Array.isArray(this.allTestSuites)) {
        this.testSuites = []
        return
      }

      // Use the utility function to sort the data
      const sortedSuites = applySorting(this.allTestSuites, this.sortBy)

      // Apply pagination
      const startIndex = (this.pagination.page - 1) * this.pagination.page_size
      const endIndex = startIndex + this.pagination.page_size
      
      // Update pagination info
      this.pagination.total = sortedSuites.length
      this.pagination.total_pages = Math.ceil(sortedSuites.length / this.pagination.page_size)
      this.pagination.has_next = this.pagination.page < this.pagination.total_pages
      this.pagination.has_prev = this.pagination.page > 1
      
      // Get current page data
      this.testSuites = sortedSuites.slice(startIndex, endIndex)
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

    editProject() {
      this.showProjectModal = true
    },

    closeProjectModal() {
      this.showProjectModal = false
    },

    handleProjectSaved() {
      this.closeProjectModal()
      this.loadData()
      showAlert('Project updated successfully!', 'success')
    },

    showCreateTestSuiteModal() {
      this.selectedTestSuite = null
      this.showTestSuiteModal = true
    },

    editTestSuite(testSuite) {
      this.selectedTestSuite = testSuite
      this.showTestSuiteModal = true
    },

    closeTestSuiteModal() {
      this.showTestSuiteModal = false
      this.selectedTestSuite = null
    },

    handleTestSuiteSaved() {
      this.closeTestSuiteModal()
      this.loadData()
      showAlert('Test suite saved successfully!', 'success')
    },

    async deleteTestSuite(testSuiteId) {
      if (!confirm('Are you sure you want to delete this test suite? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteTestSuite(testSuiteId)
        showAlert('Test suite deleted successfully!', 'success')
        this.loadData()
      } catch (error) {
        showAlert('Error deleting test suite: ' + error.message, 'danger')
      }
    },

    async syncProject() {
      if (!this.project?.repository) {
        showAlert('No Git repository configured for this project', 'warning')
        return
      }

      this.syncing = true
      try {
        const response = await api.syncRepository(this.project.repository.id)
        if (response.success) {
          showAlert(`Git repository synced successfully! Found ${response.branch_count} branches and ${response.tag_count} tags.`, 'success')
          // Optionally reload project data to show updated sync status
          this.loadData()
        } else {
          showAlert(`Sync failed: ${response.message}`, 'warning')
        }
      } catch (error) {
        showAlert('Error syncing Git repository: ' + error.message, 'danger')
      } finally {
        this.syncing = false
      }
    }
  }
}
</script>