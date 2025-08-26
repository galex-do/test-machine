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

    <div v-if="!loading && repositories.length === 0" class="text-center py-5">
      <div class="text-muted">
        <i class="fas fa-code-branch fa-3x mb-3"></i>
        <p>No repositories found. Add your first repository to get started!</p>
      </div>
    </div>

    <div v-if="!loading && repositories.length > 0">
      <div class="row">
        <div 
          v-for="repository in repositories" 
          :key="repository.id" 
          class="col-md-6 col-lg-4 mb-4"
        >
          <div class="card h-100">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-start mb-2">
                <h5 class="card-title">{{ repository.name }}</h5>
                <div class="dropdown">
                  <button 
                    class="btn btn-outline-secondary btn-sm dropdown-toggle" 
                    type="button" 
                    :id="`repo-dropdown-${repository.id}`"
                    data-bs-toggle="dropdown" 
                    aria-expanded="false"
                  >
                    <i class="fas fa-ellipsis-v"></i>
                  </button>
                  <ul class="dropdown-menu" :aria-labelledby="`repo-dropdown-${repository.id}`">
                    <li>
                      <a 
                        class="dropdown-item" 
                        href="#" 
                        @click.prevent="editRepository(repository)"
                      >
                        <i class="fas fa-edit"></i> Edit
                      </a>
                    </li>
                    <li>
                      <a 
                        class="dropdown-item text-primary" 
                        href="#" 
                        @click.prevent="syncRepository(repository)"
                        :disabled="syncing === repository.id"
                      >
                        <i class="fas fa-sync-alt" :class="{ 'fa-spin': syncing === repository.id }"></i> 
                        {{ syncing === repository.id ? 'Syncing...' : 'Sync' }}
                      </a>
                    </li>
                    <li><hr class="dropdown-divider"></li>
                    <li>
                      <a 
                        class="dropdown-item text-danger" 
                        href="#" 
                        @click.prevent="deleteRepository(repository)"
                      >
                        <i class="fas fa-trash"></i> Delete
                      </a>
                    </li>
                  </ul>
                </div>
              </div>
              
              <p class="card-text text-muted">{{ repository.description }}</p>
              
              <div class="small mb-2">
                <div class="d-flex align-items-center mb-1">
                  <i class="fas fa-link text-muted me-2"></i>
                  <a :href="repository.remote_url" target="_blank" class="text-decoration-none">
                    {{ truncateUrl(repository.remote_url) }}
                  </a>
                </div>
                
                <div v-if="repository.key" class="d-flex align-items-center mb-1">
                  <i class="fas fa-key text-muted me-2"></i>
                  <span class="badge bg-secondary">{{ repository.key.name }}</span>
                </div>
                
                <div v-if="repository.synced_at" class="d-flex align-items-center">
                  <i class="fas fa-clock text-muted me-2"></i>
                  <span>Last synced: {{ formatDate(repository.synced_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
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
import { formatDate, showAlert, truncateText } from '../utils/helpers.js'
import RepositoryModal from './modals/RepositoryModal.vue'

export default {
  name: 'Repositories',
  components: {
    RepositoryModal
  },
  data() {
    return {
      repositories: [],
      loading: true,
      showRepositoryModal: false,
      selectedRepository: null,
      syncing: null
    }
  },
  mounted() {
    this.loadRepositories()
  },
  methods: {
    formatDate,
    
    truncateUrl(url) {
      if (url.length > 40) {
        return url.substring(0, 37) + '...'
      }
      return url
    },
    
    async loadRepositories() {
      this.loading = true
      try {
        this.repositories = await api.getRepositories()
      } catch (error) {
        showAlert('Error loading repositories: ' + error.message, 'danger')
        this.repositories = []
      } finally {
        this.loading = false
      }
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
      if (!confirm(`Are you sure you want to delete repository "${repository.name}"? This action cannot be undone and will affect any projects using this repository.`)) {
        return
      }

      try {
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

.card {
  transition: box-shadow 0.15s ease-in-out;
}

.card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.dropdown-toggle::after {
  display: none;
}

.card-title {
  color: #333;
  font-size: 1.1rem;
  font-weight: 600;
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
</style>