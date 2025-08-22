<template>
  <div>
    <!-- Breadcrumb -->
    <nav aria-label="breadcrumb">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <router-link to="/">Dashboard</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link :to="`/project/${pid}`">{{ truncateText(testCase?.test_suite?.project?.name, 30) }}</router-link>
        </li>
        <li class="breadcrumb-item">
          <router-link :to="`/project/${pid}/test-suite/${sid}`">{{ truncateText(testCase?.test_suite?.name, 30) }}</router-link>
        </li>
        <li class="breadcrumb-item active">{{ truncateText(testCase?.title, 30) }}</li>
      </ol>
    </nav>

    <!-- Test Case Header -->
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div>
        <h1 class="h2">{{ testCase?.title }}</h1>
        <p class="text-muted">{{ testCase?.description || 'No description available' }}</p>
        <div class="mt-2">
          <span class="priority-badge me-2" :class="getPriorityBadgeClass(testCase?.priority)">
            {{ testCase?.priority }}
          </span>
          <span class="status-badge" :class="getStatusBadgeClass(testCase?.status)">
            {{ testCase?.status }}
          </span>
        </div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <button type="button" class="btn btn-primary" @click="showCreateTestStepModal">
          <i class="fas fa-plus"></i> Add Test Step
        </button>
      </div>
    </div>

    <!-- Test Case Info -->
    <div class="row mb-4" v-if="testCase">
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Priority</h5>
            <span class="priority-badge" :class="getPriorityBadgeClass(testCase.priority)">
              {{ testCase.priority }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Status</h5>
            <span class="status-badge" :class="getStatusBadgeClass(testCase.status)">
              {{ testCase.status }}
            </span>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Test Steps</h5>
            <h3>{{ testSteps.length }}</h3>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h5>Created</h5>
            <small>{{ formatDate(testCase.created_at) }}</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Steps -->
    <div class="card">
      <div class="card-header">
        <h5><i class="fas fa-list-ol"></i> Test Steps</h5>
      </div>
      <div class="card-body">
        <div v-if="loading" v-html="showLoading()"></div>
        <div v-else-if="testSteps.length === 0" class="empty-state">
          <i class="fas fa-list-ol"></i>
          <h5>No Test Steps Found</h5>
          <p>Add test steps to define the testing procedure.</p>
          <button class="btn btn-primary" @click="showCreateTestStepModal">
            <i class="fas fa-plus"></i> Add Test Step
          </button>
        </div>
        <div v-else>
          <div v-for="step in sortedTestSteps" :key="step.id" class="test-step">
            <div class="d-flex align-items-start">
              <div class="step-number">{{ step.step_number }}</div>
              <div class="flex-grow-1">
                <div class="mb-2">
                  <strong>Description:</strong><br>
                  {{ step.description }}
                </div>
                <div class="mb-2">
                  <strong>Expected Result:</strong><br>
                  {{ step.expected_result }}
                </div>
                <div class="action-buttons">
                  <button class="btn btn-outline-secondary btn-sm" @click="editTestStep(step)">
                    <i class="fas fa-edit"></i> Edit
                  </button>
                  <button class="btn btn-outline-danger btn-sm" @click="deleteTestStep(step.id)">
                    <i class="fas fa-trash"></i> Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Step Modal -->
    <TestStepModal 
      :show="showModal" 
      :testStep="selectedTestStep"
      :testCaseId="parseInt(cid)"
      @close="closeModal"
      @saved="handleTestStepSaved"
    />
  </div>
</template>

<script>
import { api } from '../services/api.js'
import { formatDate, showAlert, showLoading, truncateText, getStatusBadgeClass, getPriorityBadgeClass } from '../utils/helpers.js'
import TestStepModal from './modals/TestStepModal.vue'

export default {
  name: 'TestCaseDetail',
  components: {
    TestStepModal
  },
  props: {
    pid: {
      type: String,
      required: true
    },
    sid: {
      type: String,
      required: true
    },
    cid: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      testCase: null,
      testSteps: [],
      loading: true,
      showModal: false,
      selectedTestStep: null
    }
  },
  computed: {
    sortedTestSteps() {
      return [...this.testSteps].sort((a, b) => a.step_number - b.step_number)
    }
  },
  mounted() {
    this.loadData()
  },
  watch: {
    cid() {
      this.loadData()
    }
  },
  methods: {
    formatDate,
    showLoading,
    truncateText,
    getStatusBadgeClass,
    getPriorityBadgeClass,
    
    async loadData() {
      this.loading = true
      try {
        const [testCaseData, testStepsData] = await Promise.all([
          api.getTestCase(this.cid),
          api.getTestSteps(this.cid)
        ])
        this.testCase = testCaseData
        this.testSteps = testStepsData
      } catch (error) {
        showAlert('Error loading test case data: ' + error.message, 'danger')
      } finally {
        this.loading = false
      }
    },

    showCreateTestStepModal() {
      this.selectedTestStep = null
      this.showModal = true
    },

    editTestStep(testStep) {
      this.selectedTestStep = testStep
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.selectedTestStep = null
    },

    handleTestStepSaved() {
      this.closeModal()
      this.loadData()
      showAlert('Test step saved successfully!', 'success')
    },

    async deleteTestStep(testStepId) {
      if (!confirm('Are you sure you want to delete this test step? This action cannot be undone.')) {
        return
      }

      try {
        await api.deleteTestStep(testStepId)
        showAlert('Test step deleted successfully!', 'success')
        this.loadData()
      } catch (error) {
        showAlert('Error deleting test step: ' + error.message, 'danger')
      }
    }
  }
}
</script>