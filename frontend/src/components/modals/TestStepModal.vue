<template>
  <div class="modal" :class="{ show: show }" :style="{ display: show ? 'block' : 'none' }" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Edit Test Step' : 'Add New Test Step' }}</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>
        <form @submit.prevent="save">
          <div class="modal-body">
            <div class="mb-3">
              <label for="testStepNumber" class="form-label">Step Number *</label>
              <input 
                type="number" 
                class="form-control" 
                id="testStepNumber" 
                v-model.number="form.step_number" 
                required 
                min="1"
                :readonly="!isEditing"
              >
            </div>
            <div class="mb-3">
              <label for="testStepDescription" class="form-label">Step Description *</label>
              <textarea 
                class="form-control" 
                id="testStepDescription" 
                v-model="form.description" 
                rows="3" 
                required 
                placeholder="Describe what action should be performed"
              ></textarea>
            </div>
            <div class="mb-3">
              <label for="testStepExpectedResult" class="form-label">Expected Result *</label>
              <textarea 
                class="form-control" 
                id="testStepExpectedResult" 
                v-model="form.expected_result" 
                rows="3" 
                required 
                placeholder="Describe the expected outcome of this step"
              ></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="close">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update Test Step' : 'Add Test Step' }}
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
  name: 'TestStepModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    testStep: {
      type: Object,
      default: null
    },
    testCaseId: {
      type: Number,
      required: true
    }
  },
  emits: ['close', 'saved'],
  data() {
    return {
      form: {
        step_number: 1,
        description: '',
        expected_result: ''
      },
      saving: false
    }
  },
  computed: {
    isEditing() {
      return this.testStep !== null
    }
  },
  watch: {
    async show(newVal) {
      if (newVal) {
        this.resetForm()
        if (this.testStep) {
          this.form.step_number = this.testStep.step_number
          this.form.description = this.testStep.description
          this.form.expected_result = this.testStep.expected_result
        } else {
          // Auto-calculate next step number for new steps
          await this.calculateNextStepNumber()
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
        step_number: 1,
        description: '',
        expected_result: ''
      }
    },

    async calculateNextStepNumber() {
      try {
        const testSteps = await api.getTestSteps(this.testCaseId)
        let nextStepNumber = 1
        if (testSteps && testSteps.length > 0) {
          const maxStepNumber = Math.max(...testSteps.map(step => step.step_number))
          nextStepNumber = maxStepNumber + 1
        }
        this.form.step_number = nextStepNumber
      } catch (error) {
        console.error('Error calculating next step number:', error)
        this.form.step_number = 1
      }
    },

    async save() {
      this.saving = true
      try {
        if (this.isEditing) {
          await api.updateTestStep(this.testStep.id, this.form)
        } else {
          await api.createTestStep(this.testCaseId, this.form)
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