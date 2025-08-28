<template>
  <div class="repositories-container">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1>Repositories</h1>
      <button 
        type="button" 
        class="btn btn-primary"
        @click="showCreateModal"
      >
        <i class="fas fa-plus"></i> Add Repository
      </button>
    </div>

    <div v-if="loading" class="text-center">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-if="!loading && (!repositories || repositories.length === 0)" class="text-center py-5">
      <div class="text-muted">
        <i class="fas fa-code-branch fa-3x mb-3"></i>
        <p>No repositories found. Add your first repository to get started!</p>
      </div>
    </div>

    <div v-if="!loading && repositories && repositories.length > 0" class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5><i class="fas fa-code-branch"></i> Repositories</h5>
        <SortBy 
          :sortOptions="repositorySortOptions"
          :defaultSort="sortBy"
          componentId="repositories"
          @sort-changed="handleSortChange"
        />
      </div>
      <div class="card-body">
        <div class="table-responsive">
        <table class="table table-hover">
          <thead class="table-light">
            <tr>
              <th>Repository Name</th>
              <th>Remote URL</th>
              <th>Authentication Key</th>
              <th>Last Sync</th>
              <th class="text-end">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="repository in repositories" :key="repository.id">
              <td>
                <div>
                  <router-link 
                    :to="`/repository/${repository.id}`" 
                    class="text-decoration-none"
                  >
                    <strong>{{ repository.name }}</strong>
                  </router-link>
                  <div v-if="repository.description" class="text-muted small">{{ repository.description }}</div>
                </div>
              </td>
              <td>
                <a :href="repository.remote_url" target="_blank" class="text-decoration-none">
                  {{ repository.remote_url }}
                  <i class="fas fa-external-link-alt ms-1 text-muted small"></i>
                </a>
              </td>
              <td>
                <span v-if="repository.key" class="badge bg-secondary">
                  <i class="fas fa-key me-1"></i>{{ repository.key.name }}
                </span>
                <span v-else class="text-muted">
                  <i class="fas fa-globe me-1"></i>Public
                </span>
              </td>
              <td>
                <span v-if="repository.synced_at" class="text-success">
                  <i class="fas fa-check-circle me-1"></i>{{ formatDate(repository.synced_at) }}
                </span>
                <span v-else class="text-muted">
                  <i class="fas fa-clock me-1"></i>Not synced
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group" role="group">
                  <button 
                    type="button" 
                    class="btn btn-outline-primary btn-sm"
                    @click="syncRepository(repository)"
                    :disabled="syncing === repository.id"
                    :title="syncing === repository.id ? 'Syncing...' : 'Sync repository'"
                  >
                    <i class="fas fa-sync-alt" :class="{ 'fa-spin': syncing === repository.id }"></i>
                  </button>
                  <button 
                    type="button" 
                    class="btn btn-outline-secondary btn-sm"
                    @click="editRepository(repository)"
                    title="Edit repository"
                  >
                    <i class="fas fa-edit"></i>
                  </button>
                  <button 
                    type="button" 
                    class="btn btn-outline-danger btn-sm"
                    @click="deleteRepository(repository)"
                    title="Delete repository"
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
      <div class="card-footer" v-if="!loading && pagination.total > 0">
        <Pagination 
          :pagination="pagination"
          @page-changed="changePage"
          @page-size-changed="changePageSize"
        />
      </div>
    </div>

    <!-- Repository Modal -->
    <RepositoryModal 
      :show="showRepositoryModal" 
      :repository="selectedRepository"
      @close="closeRepositoryModal"
      @saved="handleRepositorySaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert } from '../utils/helpers.js'
import { applySorting, SORT_OPTION_SETS } from '../utils/sortUtils.js'
import RepositoryModal from './modals/RepositoryModal.vue'
import Pagination from './Pagination.vue'
import SortBy from './SortBy.vue'

export default {
  name: 'Repositories',
  components: {
    RepositoryModal,
    Pagination,
    SortBy
  },
  data() {
    return {
      repositories: [],
      loading: true,
      showRepositoryModal: false,
      selectedRepository: null,
      syncing: null,
      sortBy: 'created_desc',
      allRepositories: [], // Store all repositories for sorting and pagination
      repositorySortOptions: SORT_OPTION_SETS.REPOSITORIES,
      pagination: {
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      }
    }
  },
  mounted() {
    this.loadRepositories()
  },
  methods: {
    formatDate,
    
    async loadRepositories() {
      this.loading = true
      try {
        const allRepositories = await api.getRepositories()
        
        // Ensure repositories is always an array
        this.allRepositories = Array.isArray(allRepositories) ? allRepositories : []
        this.applyCurrentSorting()
        
      } catch (error) {
        showAlert('Error loading repositories: ' + error.message, 'danger')
        this.repositories = []
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
      if (!this.allRepositories || !Array.isArray(this.allRepositories)) {
        this.repositories = []
        return
      }

      // Use the utility function to sort the data
      const sortedRepositories = applySorting(this.allRepositories, this.sortBy)

      // Apply pagination
      const startIndex = (this.pagination.page - 1) * this.pagination.page_size
      const endIndex = startIndex + this.pagination.page_size
      
      // Update pagination info
      this.pagination.total = sortedRepositories.length
      this.pagination.total_pages = Math.ceil(sortedRepositories.length / this.pagination.page_size)
      this.pagination.has_next = this.pagination.page < this.pagination.total_pages
      this.pagination.has_prev = this.pagination.page > 1
      
      // Get current page data
      this.repositories = sortedRepositories.slice(startIndex, endIndex)
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

    showCreateModal() {
      this.selectedRepository = null
      this.showRepositoryModal = true
    },

    editRepository(repository) {
      this.selectedRepository = repository
      this.showRepositoryModal = true
    },

    closeRepositoryModal() {
      this.showRepositoryModal = false
      this.selectedRepository = null
    },

    handleRepositorySaved() {
      this.closeRepositoryModal()
      this.loadRepositories()
      showAlert('Repository saved successfully!', 'success')
    },

    async syncRepository(repository) {
      this.syncing = repository.id
      try {
        const response = await api.syncRepository(repository.id)
        if (response.success) {
          showAlert(`Repository synced successfully! Found ${response.branch_count} branches and ${response.tag_count} tags.`, 'success')
          this.loadRepositories() // Reload to show updated sync time
        } else {
          showAlert(`Sync failed: ${response.message}`, 'danger')
        }
      } catch (error) {
        showAlert('Error syncing repository: ' + error.message, 'danger')
      } finally {
        this.syncing = null
      }
    },

    async deleteRepository(repository) {
      // Check if repository is being used by any projects
      try {
        const projects = await api.getProjects()
        const linkedProjects = projects.filter(p => p.repository_id === repository.id)
        
        if (linkedProjects.length > 0) {
          const projectNames = linkedProjects.map(p => p.name).join(', ')
          showAlert(
            `Cannot delete repository "${repository.name}" because it is linked to the following project(s): ${projectNames}. Please unlink the repository from these projects first.`,
            'warning'
          )
          return
        }
        
        if (!confirm(`Are you sure you want to delete repository "${repository.name}"? This action cannot be undone.`)) {
          return
        }

        await api.deleteRepository(repository.id)
        showAlert('Repository deleted successfully!', 'success')
        this.loadRepositories()
      } catch (error) {
        showAlert('Error deleting repository: ' + error.message, 'danger')
      }
    }
  }
}
</script>

<style scoped>
.repositories-container {
  padding: 20px;
}

.table th {
  border-top: none;
  font-weight: 600;
}

.btn-group .btn {
  border-radius: 0.25rem;
  margin-left: 2px;
}

.btn-group .btn:first-child {
  margin-left: 0;
}

.fa-spin {
  animation: fa-spin 1s infinite linear;
}

@keyframes fa-spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(359deg);
  }
}

.table-responsive {
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
}

.table td {
  vertical-align: middle;
}

.badge {
  font-size: 0.75rem;
}
</style>