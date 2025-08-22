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
        description: ''
      },
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
        if (this.project) {
          this.form.name = this.project.name
          this.form.description = this.project.description
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
        description: ''
      }
    },

    async save() {
      this.saving = true
      try {
        if (this.isEditing) {
          await api.updateProject(this.project.id, this.form)
        } else {
          await api.createProject(this.form)
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