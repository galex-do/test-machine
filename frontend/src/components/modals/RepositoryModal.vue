<template>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            {{ repository ? 'Edit Repository' : 'Add Repository' }}
          </h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        <form @submit.prevent="handleSubmit">
          <div class="modal-body">
            <div class="mb-3">
              <label for="repositoryName" class="form-label">Repository Name *</label>
              <input
                type="text"
                class="form-control"
                id="repositoryName"
                v-model="formData.name"
                :class="{ 'is-invalid': errors.name }"
                required
                placeholder="My Project Repository"
              >
              <div v-if="errors.name" class="invalid-feedback">{{ errors.name }}</div>
            </div>

            <div class="mb-3">
              <label for="repositoryDescription" class="form-label">Description</label>
              <textarea
                class="form-control"
                id="repositoryDescription"
                v-model="formData.description"
                rows="3"
                placeholder="Optional description for this repository"
              ></textarea>
            </div>

            <div class="mb-3">
              <label for="repositoryUrl" class="form-label">
                Repository URL *
                <span v-if="repository" class="text-muted small">(immutable after creation)</span>
              </label>
              <input
                type="url"
                class="form-control"
                id="repositoryUrl"
                v-model="formData.remote_url"
                :class="{ 'is-invalid': errors.remote_url }"
                :disabled="!!repository"
                required
                placeholder="https://github.com/username/repository.git"
              >
              <div v-if="errors.remote_url" class="invalid-feedback">{{ errors.remote_url }}</div>
              <div class="form-text">
                <i class="fas fa-info-circle"></i>
                {{ repository ? 'Repository URL cannot be changed after creation for security reasons.' : 'Supports both HTTPS and SSH Git URLs' }}
              </div>
            </div>

            <div class="mb-3">
              <label for="repositoryKey" class="form-label">Authentication Key</label>
              <select
                class="form-select"
                id="repositoryKey"
                v-model="formData.key_id"
                :class="{ 'is-invalid': errors.key_id }"
              >
                <option :value="null">No authentication (public repository)</option>
                <option 
                  v-for="key in keys" 
                  :key="key.id" 
                  :value="key.id"
                >
                  {{ key.name }} ({{ key.key_type }})
                </option>
              </select>
              <div v-if="errors.key_id" class="invalid-feedback">{{ errors.key_id }}</div>
              <div class="form-text">
                <i class="fas fa-info-circle"></i>
                Select an authentication key for private repositories. 
                <router-link to="/keys" class="text-decoration-none">Manage keys</router-link>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="loading">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              {{ repository ? 'Update Repository' : 'Create Repository' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { api } from '../../services/api.js'
import { showAlert } from '../../utils/helpers.js'

export default {
  name: 'RepositoryModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    repository: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'saved'],
  data() {
    return {
      formData: {
        name: '',
        description: '',
        remote_url: '',
        key_id: null
      },
      keys: [],
      errors: {},
      loading: false
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.loadKeys()
        this.resetForm()
      }
    }
  },
  methods: {
    async loadKeys() {
      try {
        this.keys = await api.getKeys()
      } catch (error) {
        console.error('Error loading keys:', error)
        this.keys = []
      }
    },

    resetForm() {
      if (this.repository) {
        // Editing existing repository
        this.formData = {
          name: this.repository.name || '',
          description: this.repository.description || '',
          remote_url: this.repository.remote_url || '',
          key_id: this.repository.key_id || null
        }
      } else {
        // Creating new repository
        this.formData = {
          name: '',
          description: '',
          remote_url: '',
          key_id: null
        }
      }
      this.errors = {}
    },

    validateForm() {
      this.errors = {}

      if (!this.formData.name?.trim()) {
        this.errors.name = 'Repository name is required'
      }

      if (!this.repository && !this.formData.remote_url?.trim()) {
        this.errors.remote_url = 'Repository URL is required'
      }

      if (this.formData.remote_url && !this.isValidUrl(this.formData.remote_url)) {
        this.errors.remote_url = 'Please enter a valid Git repository URL'
      }

      return Object.keys(this.errors).length === 0
    },

    isValidUrl(string) {
      try {
        new URL(string)
        return true
      } catch (_) {
        // Also accept SSH URLs
        return /^git@[\w.-]+:[\w.-]+\/[\w.-]+\.git$/.test(string)
      }
    },

    async handleSubmit() {
      if (!this.validateForm()) {
        return
      }

      this.loading = true
      try {
        if (this.repository) {
          // Update existing repository (only name, description, key_id)
          await api.updateRepository(this.repository.id, {
            name: this.formData.name,
            description: this.formData.description,
            key_id: this.formData.key_id
          })
        } else {
          // Create new repository
          await api.createRepository(this.formData)
        }
        
        this.$emit('saved')
      } catch (error) {
        showAlert('Error saving repository: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.modal.show {
  display: block;
}

.form-text {
  font-size: 0.875rem;
}

.form-text i {
  margin-right: 4px;
}
</style>