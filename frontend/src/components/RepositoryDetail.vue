<template>
  <div class="repository-detail-container">
    <nav aria-label="breadcrumb" class="mb-4">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/repositories">Repositories</router-link>
        </li>
        <li class="breadcrumb-item active" aria-current="page">
          {{ repository?.name || 'Repository Details' }}
        </li>
      </ol>
    </nav>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-else-if="error" class="alert alert-danger">
      <i class="fas fa-exclamation-triangle"></i> {{ error }}
    </div>

    <div v-else-if="repository">
      <!-- Repository Header -->
      <div class="row mb-4">
        <div class="col-md-8">
          <div class="d-flex align-items-center mb-3">
            <h1 class="me-3">{{ repository.name }}</h1>
            <div class="btn-group">
              <button 
                type="button" 
                class="btn btn-outline-primary"
                @click="syncRepository"
                :disabled="syncing"
              >
                <i class="fas fa-sync-alt" :class="{ 'fa-spin': syncing }"></i>
                {{ syncing ? 'Syncing...' : 'Sync' }}
              </button>
              <button 
                type="button" 
                class="btn btn-outline-secondary"
                @click="editRepository"
              >
                <i class="fas fa-edit"></i> Edit
              </button>
            </div>
          </div>
          
          <p v-if="repository.description" class="text-muted">
            {{ repository.description }}
          </p>
        </div>
        
        <div class="col-md-4">
          <div class="card">
            <div class="card-body">
              <h6 class="card-title">Repository Info</h6>
              <div class="mb-2">
                <small class="text-muted">Remote URL:</small><br>
                <a :href="repository.remote_url" target="_blank" class="text-break">
                  {{ repository.remote_url }} <i class="fas fa-external-link-alt"></i>
                </a>
              </div>
              
              <div class="mb-2">
                <small class="text-muted">Authentication:</small><br>
                <span v-if="repository.key" class="badge bg-secondary">
                  <i class="fas fa-key me-1"></i>{{ repository.key.name }}
                </span>
                <span v-else class="badge bg-success">
                  <i class="fas fa-globe me-1"></i>Public
                </span>
              </div>
              
              <div class="mb-2">
                <small class="text-muted">Default Branch:</small><br>
                <span v-if="repository.default_branch" class="badge bg-primary">
                  {{ repository.default_branch }}
                </span>
                <span v-else class="text-muted">Not set</span>
              </div>
              
              <div>
                <small class="text-muted">Last Sync:</small><br>
                <span v-if="repository.synced_at" class="text-success">
                  <i class="fas fa-check-circle me-1"></i>{{ formatDate(repository.synced_at) }}
                </span>
                <span v-else class="text-warning">
                  <i class="fas fa-clock me-1"></i>Not synced
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Branches and Tags -->
      <div class="row">
        <!-- Branches -->
        <div class="col-md-6 mb-4">
          <div class="card h-100">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h5 class="mb-0">
                <i class="fas fa-code-branch text-primary"></i> Branches
              </h5>
              <span class="badge bg-primary">{{ repository.branches?.length || 0 }}</span>
            </div>
            <div class="card-body">
              <div v-if="!repository.branches || repository.branches.length === 0" class="text-center py-4 text-muted">
                <i class="fas fa-info-circle fa-2x mb-2"></i>
                <p class="mb-0">No branches found</p>
                <small>Try syncing the repository to fetch branches</small>
              </div>
              
              <div v-else>
                <div class="table-responsive">
                  <table class="table table-sm">
                    <thead>
                      <tr>
                        <th>Branch Name</th>
                        <th>Latest Commit</th>
                        <th>Date</th>
                        <th>Message</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="branch in repository.branches" :key="branch.id">
                        <td>
                          <span class="d-flex align-items-center">
                            {{ branch.name }}
                            <span v-if="branch.is_default" class="badge bg-success ms-2 small">
                              DEFAULT
                            </span>
                          </span>
                        </td>
                        <td>
                          <code v-if="branch.commit_hash" class="small">
                            {{ truncateHash(branch.commit_hash) }}
                          </code>
                          <span v-else class="text-muted">-</span>
                        </td>
                        <td class="small text-muted">
                          <span v-if="branch.commit_date">
                            {{ formatDate(branch.commit_date) }}
                          </span>
                          <span v-else class="text-muted">-</span>
                        </td>
                        <td class="small">
                          <span v-if="branch.commit_message" class="text-truncate" :title="branch.commit_message" style="max-width: 200px; display: inline-block;">
                            {{ truncateMessage(branch.commit_message) }}
                          </span>
                          <span v-else class="text-muted">-</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Tags -->
        <div class="col-md-6 mb-4">
          <div class="card h-100">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h5 class="mb-0">
                <i class="fas fa-tag text-warning"></i> Tags
              </h5>
              <span class="badge bg-warning">{{ repository.tags?.length || 0 }}</span>
            </div>
            <div class="card-body">
              <div v-if="!repository.tags || repository.tags.length === 0" class="text-center py-4 text-muted">
                <i class="fas fa-info-circle fa-2x mb-2"></i>
                <p class="mb-0">No tags found</p>
                <small>Try syncing the repository to fetch tags</small>
              </div>
              
              <div v-else>
                <div class="table-responsive">
                  <table class="table table-sm">
                    <thead>
                      <tr>
                        <th>Tag Name</th>
                        <th>Commit Hash</th>
                        <th>Date</th>
                        <th>Message</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="tag in repository.tags" :key="tag.id">
                        <td>
                          <span class="badge bg-outline-warning">{{ tag.name }}</span>
                        </td>
                        <td>
                          <code v-if="tag.commit_hash" class="small">
                            {{ truncateHash(tag.commit_hash) }}
                          </code>
                          <span v-else class="text-muted">-</span>
                        </td>
                        <td class="small text-muted">
                          <span v-if="tag.commit_date">
                            {{ formatDate(tag.commit_date) }}
                          </span>
                          <span v-else class="text-muted">-</span>
                        </td>
                        <td class="small">
                          <span v-if="tag.commit_message" class="text-truncate" :title="tag.commit_message" style="max-width: 200px; display: inline-block;">
                            {{ truncateMessage(tag.commit_message) }}
                          </span>
                          <span v-else class="text-muted">-</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Repository Edit Modal -->
    <RepositoryModal 
      :show="showEditModal" 
      :repository="repository"
      @close="closeEditModal"
      @saved="handleRepositorySaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert } from '../utils/helpers.js'
import RepositoryModal from './modals/RepositoryModal.vue'

export default {
  name: 'RepositoryDetail',
  components: {
    RepositoryModal
  },
  props: {
    id: {
      type: [String, Number],
      required: true
    }
  },
  data() {
    return {
      repository: null,
      loading: true,
      error: null,
      syncing: false,
      showEditModal: false
    }
  },
  mounted() {
    this.loadRepository()
  },
  watch: {
    id(newId) {
      this.loadRepository()
    }
  },
  methods: {
    formatDate,
    
    async loadRepository() {
      this.loading = true
      this.error = null
      try {
        this.repository = await api.getRepositoryDetails(this.id)
      } catch (error) {
        this.error = 'Error loading repository: ' + error.message
      } finally {
        this.loading = false
      }
    },

    truncateHash(hash) {
      if (!hash) return '-'
      return hash.substring(0, 8)
    },

    truncateMessage(message) {
      if (!message) return '-'
      const firstLine = message.split('\n')[0]
      return firstLine.length > 50 ? firstLine.substring(0, 50) + '...' : firstLine
    },

    async syncRepository() {
      this.syncing = true
      try {
        const response = await api.syncRepository(this.repository.id)
        if (response.success) {
          showAlert(`Repository synced successfully! Found ${response.branch_count} branches and ${response.tag_count} tags.`, 'success')
          this.loadRepository() // Reload to show updated data
        } else {
          showAlert(`Sync failed: ${response.message}`, 'warning')
        }
      } catch (error) {
        showAlert('Error syncing repository: ' + error.message, 'danger')
      } finally {
        this.syncing = false
      }
    },

    editRepository() {
      this.showEditModal = true
    },

    closeEditModal() {
      this.showEditModal = false
    },

    handleRepositorySaved() {
      this.closeEditModal()
      this.loadRepository()
      showAlert('Repository updated successfully!', 'success')
    }
  }
}
</script>

<style scoped>
.repository-detail-container {
  padding: 20px;
}

.card {
  border: 1px solid #dee2e6;
}

.card-header {
  background-color: #f8f9fa;
  font-weight: 600;
}

.badge.bg-outline-warning {
  background-color: transparent !important;
  border: 1px solid #ffc107;
  color: #ffc107;
}

code {
  background-color: #f8f9fa;
  padding: 2px 4px;
  border-radius: 3px;
}

.table th {
  font-weight: 600;
  font-size: 0.875rem;
}

.table td {
  vertical-align: middle;
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