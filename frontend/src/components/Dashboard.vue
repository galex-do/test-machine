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
      <div class="card-header">
        <h5><i class="fas fa-folder"></i> Projects</h5>
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
              <tr v-for="project in projects" :key="project.id">
                <td>
                  <router-link :to="`/project/${project.id}`" class="text-decoration-none">
                    <strong>{{ project.name }}</strong>
                  </router-link>
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
import ProjectModal from './modals/ProjectModal.vue'

export default {
  name: 'Dashboard',
  components: {
    ProjectModal
  },
  data() {
    return {
      projects: [],
      stats: null,
      loading: true,
      showModal: false,
      selectedProject: null
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
        this.projects = projectsData
        this.stats = statsData
      } catch (error) {
        showAlert('Error loading dashboard data: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
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