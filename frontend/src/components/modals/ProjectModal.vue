<template>
  <div class="modal" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Edit Project' : 'Create New Project' }}</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>
        <form @submit.prevent="save">
          <div class="modal-body">
            <div class="mb-3">
              <label for="projectName" class="form-label">Project Name *</label>
              <input type="text" class="form-control" id="projectName" v-model="form.name" required>
            </div>
            <div class="mb-3">
              <label for="projectDescription" class="form-label">Description</label>
              <textarea class="form-control" id="projectDescription" v-model="form.description" rows="3"></textarea>
            </div>
            <div class="mb-3">
              <label for="projectGitProject" class="form-label">Git Project <small class="text-muted">(Optional)</small></label>
              <input type="url" class="form-control" id="projectGitProject" v-model="form.git_project" placeholder="https://github.com/username/project">
              <div class="form-text">Link to your Git repository (GitHub, GitLab, etc.) for test run integration</div>
            </div>
            
            <!-- Key Selector - only show when Git project is filled -->
            <div v-if="form.git_project && form.git_project.trim()" class="mb-3">
              <label for="projectKey" class="form-label">Authentication Key <small class="text-muted">(Optional)</small></label>
              <select class="form-select" id="projectKey" v-model="form.key_id">
                <option value="">No authentication required</option>
                <option v-for="key in keys" :key="key.id" :value="key.id">
                  {{ key.name }} ({{ key.key_type }}){{ key.username ? ` - ${key.username}` : '' }}
                </option>
              </select>
              <div class="form-text">
                <i class="fas fa-key"></i>
                Choose authentication credentials for this Git repository
              </div>
              <div v-if="keys.length === 0" class="form-text text-warning">
                <i class="fas fa-exclamation-triangle"></i>
                No keys available. <router-link to="/keys" target="_blank">Create keys</router-link> for Git authentication.
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update Project' : 'Create Project' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
  <div v-if="show" class="modal-backdrop fade show"></div>
</template>

<script>
import { api } from '../../services/api.js'

export default {
  name: 'ProjectModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    project: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'saved'],
  data() {
    return {
      form: {
        name: '',
        description: '',
        git_project: '',
        key_id: ''
      },
      keys: [],
      saving: false
    }
  },
  computed: {
    isEditing() {
      return this.project !== null
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm()
        this.loadKeys()
        if (this.project) {
          this.form.name = this.project.name
          this.form.description = this.project.description
          this.form.git_project = this.project.git_project || ''
          this.form.key_id = this.project.key_id || ''
        }
      }
    }
  },
  methods: {
    close() {
      this.$emit('close')
    },

    resetForm() {
      this.form = {
        name: '',
        description: '',
        git_project: '',
        key_id: ''
      }
    },

    async loadKeys() {
      try {
        const data = await api.getKeys()
        this.keys = Array.isArray(data) ? data : []
      } catch (error) {
        this.keys = []
        console.error('Error loading keys:', error)
      }
    },

    async save() {
      this.saving = true
      try {
        // Prepare form data, ensuring key_id is either a number or null
        const formData = {
          ...this.form,
          key_id: this.form.key_id && this.form.key_id !== '' ? parseInt(this.form.key_id) : null
        }

        if (this.isEditing) {
          await api.updateProject(this.project.id, formData)
        } else {
          await api.createProject(formData)
        }
        this.$emit('saved')
      } catch (error) {
        // Error will be handled by parent component
        throw error
      } finally {
        this.saving = false
      }
    }
  }
}
</script>