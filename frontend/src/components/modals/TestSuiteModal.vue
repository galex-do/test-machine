<template>
  <div class="modal" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Edit Test Suite' : 'Create New Test Suite' }}</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>
        <form @submit.prevent="save">
          <div class="modal-body">
            <div class="mb-3">
              <label for="testSuiteName" class="form-label">Test Suite Name *</label>
              <input type="text" class="form-control" id="testSuiteName" v-model="form.name" required>
            </div>
            <div class="mb-3">
              <label for="testSuiteDescription" class="form-label">Description</label>
              <textarea class="form-control" id="testSuiteDescription" v-model="form.description" rows="3"></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update Test Suite' : 'Create Test Suite' }}
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
  name: 'TestSuiteModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    testSuite: {
      type: Object,
      default: null
    },
    projectId: {
      type: Number,
      required: true
    }
  },
  emits: ['close', 'saved'],
  data() {
    return {
      form: {
        name: '',
        description: '',
        project_id: this.projectId
      },
      saving: false
    }
  },
  computed: {
    isEditing() {
      return this.testSuite !== null
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm()
        if (this.testSuite) {
          this.form.name = this.testSuite.name
          this.form.description = this.testSuite.description
        }
        this.form.project_id = this.projectId
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
        project_id: this.projectId
      }
    },

    async save() {
      this.saving = true
      try {
        if (this.isEditing) {
          await api.updateTestSuite(this.testSuite.id, this.form)
        } else {
          await api.createTestSuite(this.form)
        }
        this.$emit('saved')
      } catch (error) {
        throw error
      } finally {
        this.saving = false
      }
    }
  }
}
</script>