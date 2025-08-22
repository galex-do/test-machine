<template>
  <div class="modal" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Edit Test Case' : 'Create New Test Case' }}</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>
        <form @submit.prevent="save">
          <div class="modal-body">
            <div class="mb-3">
              <label for="testCaseTitle" class="form-label">Test Case Title *</label>
              <input type="text" class="form-control" id="testCaseTitle" v-model="form.title" required>
            </div>
            <div class="mb-3">
              <label for="testCaseDescription" class="form-label">Description</label>
              <textarea class="form-control" id="testCaseDescription" v-model="form.description" rows="3"></textarea>
            </div>
            <div class="row">
              <div class="col-md-6 mb-3">
                <label for="testCasePriority" class="form-label">Priority *</label>
                <select class="form-select" id="testCasePriority" v-model="form.priority" required>
                  <option value="">Select Priority</option>
                  <option value="Critical">Critical</option>
                  <option value="High">High</option>
                  <option value="Medium">Medium</option>
                </select>
              </div>
              <div class="col-md-6 mb-3">
                <label for="testCaseStatus" class="form-label">Status *</label>
                <select class="form-select" id="testCaseStatus" v-model="form.status" required>
                  <option value="">Select Status</option>
                  <option value="Active">Active</option>
                  <option value="Disabled">Disabled</option>
                </select>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update Test Case' : 'Create Test Case' }}
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
  name: 'TestCaseModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    testCase: {
      type: Object,
      default: null
    },
    testSuiteId: {
      type: Number,
      required: true
    }
  },
  emits: ['close', 'saved'],
  data() {
    return {
      form: {
        title: '',
        description: '',
        priority: 'Medium',
        status: 'Active',
        test_suite_id: this.testSuiteId
      },
      saving: false
    }
  },
  computed: {
    isEditing() {
      return this.testCase !== null
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm()
        if (this.testCase) {
          this.form.title = this.testCase.title
          this.form.description = this.testCase.description
          this.form.priority = this.testCase.priority
          this.form.status = this.testCase.status
        }
        this.form.test_suite_id = this.testSuiteId
      }
    }
  },
  methods: {
    close() {
      this.$emit('close')
    },

    resetForm() {
      this.form = {
        title: '',
        description: '',
        priority: 'Medium',
        status: 'Active',
        test_suite_id: this.testSuiteId
      }
    },

    async save() {
      this.saving = true
      try {
        if (this.isEditing) {
          await api.updateTestCase(this.testCase.id, this.form)
        } else {
          await api.createTestCase(this.form)
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