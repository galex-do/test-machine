<template>
  <div>
    <!-- Page Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Dashboard</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <button type="button" class="btn btn-primary" @click="showCreateProjectModal">
          <i class="fas fa-plus"></i> New Project
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
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

    <!-- Projects List -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5><i class="fas fa-folder"></i> Projects</h5>
        <SortBy 
          :sortOptions="projectSortOptions"
          :defaultSort="sortBy"
          componentId="projects"
          @sort-changed="handleSortChange"
        />
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="projects.length === 0" class="empty-state">
          <i class="fas fa-folder-open"></i>
          <h5>No Projects Found</h5>
          <p>Create your first project to get started with test management.</p>
          <button class="btn btn-primary" @click="showCreateProjectModal">
            <i class="fas fa-plus"></i> Create Project
          </button>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Test Suites</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="project in sortedProjects" :key="project.id">
                <td>
                  <router-link :to="`/project/${project.id}`" class="text-decoration-none">
                    <strong>{{ project.name }}</strong>
                  </router-link>
                  <a v-if="project.repository" 
                     :href="project.repository.remote_url" 
                     target="_blank" 
                     rel="noopener noreferrer" 
                     class="ms-2 text-muted" 
                     :title="`Open Git repository: ${project.repository.remote_url}`">
                    <i class="fab fa-git-alt"></i>
                  </a>
                </td>
                <td>{{ project.description || 'No description' }}</td>
                <td>
                  <span class="badge bg-primary">{{ project.test_suites_count || 0 }}</span>
                </td>
                <td>{{ formatDate(project.created_at) }}</td>
                <td>
                  <div class="action-buttons">
                    <button class="btn btn-outline-primary btn-sm" @click="editProject(project)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                    <button class="btn btn-outline-danger btn-sm" @click="deleteProject(project.id)">
                      <i class="fas fa-trash"></i> Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="card-footer" v-if="!loading && sortedProjects.length > 0">
          <Pagination 
            :pagination="pagination"
            @page-changed="changePage"
            @page-size-changed="changePageSize"
          />
        </div>
      </div>
    </div>

    <!-- Create/Edit Project Modal -->
    <ProjectModal 
      :show="showModal" 
      :project="selectedProject"
      @close="closeModal"
      @saved="handleProjectSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading } from '../utils/helpers.js'
import { applySorting, SORT_OPTION_SETS } from '../utils/sortUtils.js'
import ProjectModal from './modals/ProjectModal.vue'
import SortBy from './SortBy.vue'
import Pagination from './Pagination.vue'

export default {
  name: 'Dashboard',
  components: {
    ProjectModal,
    SortBy,
    Pagination
  },
  data() {
    return {
      allProjects: [], // Store all projects for sorting and pagination
      projects: [], // Current page projects
      stats: null,
      loading: true,
      showModal: false,
      selectedProject: null,
      sortBy: 'created_desc',
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      },
      projectSortOptions: SORT_OPTION_SETS.PROJECTS
    }
  },
  computed: {
    sortedProjects() {
      // For paginated results, we return the current page data already sorted
      return this.projects || []
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    formatDate,
    showLoading,
    
    async loadData() {
      this.loading = true
      try {
        const [projectsData, statsData] = await Promise.all([
          api.getProjects(),
          api.getStats().catch(() => null) // Stats endpoint might not exist yet
        ])
        this.allProjects = Array.isArray(projectsData) ? projectsData : []
        this.stats = statsData
        this.applyCurrentSorting()
      } catch (error) {
        showAlert('Error loading dashboard data: ' + error.message, 'danger')
        this.allProjects = []
        this.projects = []
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
      if (!this.allProjects || !Array.isArray(this.allProjects)) {
        this.projects = []
        return
      }

      // Use the utility function to sort the data
      const sortedProjects = applySorting(this.allProjects, this.sortBy)

      // Apply pagination
      const startIndex = (this.pagination.page - 1) * this.pagination.page_size
      const endIndex = startIndex + this.pagination.page_size
      
      // Update pagination info
      this.pagination.total = sortedProjects.length
      this.pagination.total_pages = Math.ceil(sortedProjects.length / this.pagination.page_size)
      this.pagination.has_next = this.pagination.page < this.pagination.total_pages
      this.pagination.has_prev = this.pagination.page > 1
      
      // Get current page data
      this.projects = sortedProjects.slice(startIndex, endIndex)
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

    showCreateProjectModal() {
      this.selectedProject = null
      this.showModal = true
    },

    editProject(project) {
      this.selectedProject = project
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.selectedProject = null
    },

    handleProjectSaved() {
      this.closeModal()
      this.loadData()
      showAlert('Project saved successfully!', 'success')
    },

    async deleteProject(projectId) {
      if (!confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteProject(projectId)
        showAlert('Project deleted successfully!', 'success')
        this.loadData()
      } catch (error) {
        showAlert('Error deleting project: ' + error.message, 'danger')
      }
    }
  }
}
</script>