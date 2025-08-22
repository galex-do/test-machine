<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item active">Keys</li>
      </ol>
    </nav>

    <!-- Page Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2"><i class="fas fa-key"></i> Authentication Keys</h1>
        <p class="text-muted">Manage SSH keys and login credentials for Git repository integration</p>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <button type="button" class="btn btn-primary" @click="showCreateKeyModal">
          <i class="fas fa-plus"></i> New Key
        </button>
      </div>
    </div>

    <!-- Keys List -->
    <div class="card">
      <div class="card-header">
        <h5><i class="fas fa-key"></i> Authentication Keys</h5>
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="keys.length === 0" class="empty-state">
          <i class="fas fa-key"></i>
          <h5>No Keys Found</h5>
          <p>Create authentication keys to connect with Git repositories.</p>
          <button class="btn btn-primary" @click="showCreateKeyModal">
            <i class="fas fa-plus"></i> Create Key
          </button>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Username</th>
                <th>Description</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="key in keys" :key="key.id">
                <td>
                  <strong>{{ key.name }}</strong>
                </td>
                <td>
                  <span class="badge" :class="key.key_type === 'RSA' ? 'bg-info' : 'bg-success'">
                    <i :class="key.key_type === 'RSA' ? 'fas fa-lock' : 'fas fa-user'"></i>
                    {{ key.key_type }}
                  </span>
                </td>
                <td>{{ key.username || '-' }}</td>
                <td>{{ key.description || 'No description' }}</td>
                <td>{{ formatDate(key.created_at) }}</td>
                <td>
                  <div class="action-buttons">
                    <button class="btn btn-outline-primary btn-sm" @click="editKey(key)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                    <button class="btn btn-outline-info btn-sm" @click="viewKeyData(key)" :disabled="loadingKeyData">
                      <i class="fas fa-eye"></i> View
                    </button>
                    <button class="btn btn-outline-danger btn-sm" @click="deleteKey(key.id)">
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

    <!-- Key Modal -->
    <KeyModal 
      :show="showModal" 
      :keyItem="selectedKey"
      @close="closeModal"
      @saved="handleKeySaved"
    />

    <!-- Key Data Modal -->
    <div class="modal" :class="{ show: showDataModal }" :style="{ display: showDataModal ? 'block' : 'none' }" tabindex="-1">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="fas fa-eye"></i> Key Data: {{ selectedKey?.name }}
            </h5>
            <button type="button" class="btn-close" @click="closeDataModal"></button>
          </div>
          <div class="modal-body">
            <div v-if="loadingKeyData" class="text-center">
              <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            <div v-else-if="keyDataError" class="alert alert-danger">
              {{ keyDataError }}
            </div>
            <div v-else>
              <div class="mb-3">
                <label class="form-label"><strong>Type:</strong></label>
                <span class="badge ms-2" :class="selectedKey?.key_type === 'RSA' ? 'bg-info' : 'bg-success'">
                  {{ selectedKey?.key_type }}
                </span>
              </div>
              <div v-if="selectedKey?.username" class="mb-3">
                <label class="form-label"><strong>Username:</strong></label>
                <div>{{ selectedKey.username }}</div>
              </div>
              <div class="mb-3">
                <label class="form-label"><strong>Secret Data:</strong></label>
                <textarea class="form-control" :value="keyData" readonly rows="10" style="font-family: monospace; font-size: 0.9em;"></textarea>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeDataModal">Close</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showDataModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading } from '../utils/helpers.js'
import KeyModal from './modals/KeyModal.vue'

export default {
  name: 'Keys',
  components: {
    KeyModal
  },
  data() {
    return {
      keys: [],
      loading: true,
      showModal: false,
      showDataModal: false,
      selectedKey: null,
      keyData: '',
      loadingKeyData: false,
      keyDataError: ''
    }
  },
  mounted() {
    this.loadKeys()
  },
  methods: {
    formatDate,
    showLoading,
    
    async loadKeys() {
      this.loading = true
      try {
        const data = await api.getKeys()
        this.keys = Array.isArray(data) ? data : []
      } catch (error) {
        showAlert('Error loading keys: ' + error.message, 'danger')
        this.keys = []
      } finally {
        this.loading = false
      }
    },

    showCreateKeyModal() {
      this.selectedKey = null
      this.showModal = true
    },

    editKey(key) {
      this.selectedKey = key
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.selectedKey = null
    },

    handleKeySaved() {
      this.closeModal()
      this.loadKeys()
      showAlert('Key saved successfully!', 'success')
    },

    async viewKeyData(key) {
      this.selectedKey = key
      this.keyData = ''
      this.keyDataError = ''
      this.loadingKeyData = true
      this.showDataModal = true

      try {
        const response = await api.getKeyData(key.id)
        this.keyData = response.data
      } catch (error) {
        this.keyDataError = 'Error loading key data: ' + error.message
      } finally {
        this.loadingKeyData = false
      }
    },

    closeDataModal() {
      this.showDataModal = false
      this.selectedKey = null
      this.keyData = ''
      this.keyDataError = ''
    },

    async deleteKey(keyId) {
      if (!confirm('Are you sure you want to delete this key? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteKey(keyId)
        showAlert('Key deleted successfully!', 'success')
        this.loadKeys()
      } catch (error) {
        showAlert('Error deleting key: ' + error.message, 'danger')
      }
    }
  }
}
</script>