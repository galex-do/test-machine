<template>
  <div class="modal" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" tabindex="-1">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Edit Key' : 'Create New Key' }}</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>
        <form @submit.prevent="save">
          <div class="modal-body">
            <div class="row">
              <div class="col-md-6">
                <div class="mb-3">
                  <label for="keyName" class="form-label">Key Name *</label>
                  <input type="text" class="form-control" id="keyName" v-model="form.name" required>
                </div>
              </div>
              <div class="col-md-6">
                <div class="mb-3">
                  <label for="keyType" class="form-label">Key Type *</label>
                  <select class="form-select" id="keyType" v-model="form.key_type" required :disabled="isEditing">
                    <option value="">Select Type</option>
                    <option value="RSA">RSA (SSH Key)</option>
                    <option value="Login">Login (Username/Password)</option>
                  </select>
                  <div class="form-text">
                    <i class="fas fa-info-circle"></i>
                    RSA for SSH connections, Login for HTTPS with credentials
                  </div>
                </div>
              </div>
            </div>

            <div class="mb-3">
              <label for="keyDescription" class="form-label">Description</label>
              <textarea class="form-control" id="keyDescription" v-model="form.description" rows="2" placeholder="Describe this key's purpose"></textarea>
            </div>

            <!-- Username field for Login type -->
            <div v-if="form.key_type === 'Login'" class="mb-3">
              <label for="keyUsername" class="form-label">Username *</label>
              <input type="text" class="form-control" id="keyUsername" v-model="form.username" required>
              <div class="form-text">Git repository username</div>
            </div>

            <!-- Secret data field -->
            <div class="mb-3">
              <label for="keySecretData" class="form-label">
                {{ form.key_type === 'RSA' ? 'Private Key *' : form.key_type === 'Login' ? 'Password *' : 'Secret Data *' }}
              </label>
              <textarea 
                class="form-control" 
                id="keySecretData" 
                v-model="form.secret_data" 
                :rows="form.key_type === 'RSA' ? 10 : 3"
                :placeholder="getSecretDataPlaceholder()"
                style="font-family: monospace; font-size: 0.9em;"
                :required="!isEditing">
              </textarea>
              <div class="form-text">
                <i class="fas fa-lock text-warning"></i>
                {{ form.key_type === 'RSA' ? 'Paste your private SSH key (begins with -----BEGIN OPENSSH PRIVATE KEY-----)' : 
                   form.key_type === 'Login' ? 'Enter the password for Git repository access' : 
                   'This data will be encrypted before storage' }}
              </div>
              <div v-if="isEditing" class="form-text text-muted">
                <i class="fas fa-info-circle"></i>
                Leave empty to keep existing secret data unchanged
              </div>
            </div>

            <!-- Security notice -->
            <div class="alert alert-info">
              <i class="fas fa-shield-alt"></i>
              <strong>Security Notice:</strong> All sensitive data is encrypted before storage in the database using AES-256 encryption.
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update Key' : 'Create Key' }}
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
  name: 'KeyModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    keyItem: {
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
        key_type: '',
        username: '',
        secret_data: ''
      },
      saving: false
    }
  },
  computed: {
    isEditing() {
      return this.keyItem !== null
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm()
        if (this.keyItem) {
          this.form.name = this.keyItem.name
          this.form.description = this.keyItem.description
          this.form.key_type = this.keyItem.key_type
          this.form.username = this.keyItem.username || ''
          // Don't populate secret_data when editing for security
          this.form.secret_data = ''
        }
      }
    },
    'form.key_type'(newType) {
      // Clear username when switching away from Login type
      if (newType !== 'Login') {
        this.form.username = ''
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
        key_type: '',
        username: '',
        secret_data: ''
      }
    },

    getSecretDataPlaceholder() {
      if (this.form.key_type === 'RSA') {
        return '-----BEGIN OPENSSH PRIVATE KEY-----\n...\n-----END OPENSSH PRIVATE KEY-----'
      } else if (this.form.key_type === 'Login') {
        return 'Enter password'
      }
      return 'Enter secret data'
    },

    async save() {
      this.saving = true
      try {
        // Prepare data for API
        const data = {
          name: this.form.name,
          description: this.form.description,
          key_type: this.form.key_type,
          username: this.form.username || null,
          secret_data: this.form.secret_data
        }

        // For editing, only include secret_data if it's provided
        if (this.isEditing && !this.form.secret_data) {
          delete data.secret_data
        }

        if (this.isEditing) {
          await api.updateKey(this.keyItem.id, data)
        } else {
          await api.createKey(data)
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