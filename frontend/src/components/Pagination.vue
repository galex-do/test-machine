<template>
  <div class="pagination-container d-flex justify-content-between align-items-center mt-4">
    <!-- Page Size Selector -->
    <div class="d-flex align-items-center">
      <label for="page-size-select" class="form-label me-2 mb-0">Show:</label>
      <select 
        id="page-size-select"
        class="form-select form-select-sm"
        v-model="currentPageSize"
        @change="changePageSize"
        style="width: auto; min-width: 80px;"
      >
        <option value="10">10</option>
        <option value="25">25</option>
        <option value="100">100</option>
      </select>
      <span class="ms-2 text-muted">
        Showing {{ startItem }}-{{ endItem }} of {{ pagination.total }} results
      </span>
    </div>

    <!-- Pagination Navigation -->
    <div class="d-flex align-items-center" v-if="pagination.total_pages > 1">
      <!-- Previous Button -->
      <button 
        class="btn btn-sm btn-outline-secondary me-2"
        :disabled="!pagination.has_prev"
        @click="goToPage(pagination.page - 1)"
      >
        <i class="fas fa-chevron-left"></i> Previous
      </button>

      <!-- Page Numbers -->
      <div class="d-flex align-items-center">
        <button
          v-for="page in visiblePages"
          :key="page"
          class="btn btn-sm me-1"
          :class="page === pagination.page ? 'btn-primary' : 'btn-outline-secondary'"
          @click="goToPage(page)"
          :disabled="page === '...'"
        >
          {{ page }}
        </button>
      </div>

      <!-- Next Button -->
      <button 
        class="btn btn-sm btn-outline-secondary ms-2"
        :disabled="!pagination.has_next"
        @click="goToPage(pagination.page + 1)"
      >
        Next <i class="fas fa-chevron-right"></i>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Pagination',
  props: {
    pagination: {
      type: Object,
      required: true,
      default: () => ({
        page: 1,
        page_size: 25,
        total: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      })
    }
  },
  data() {
    return {
      currentPageSize: this.pagination.page_size || 25
    }
  },
  computed: {
    startItem() {
      if (this.pagination.total === 0) return 0
      return ((this.pagination.page - 1) * this.pagination.page_size) + 1
    },
    endItem() {
      const end = this.pagination.page * this.pagination.page_size
      return Math.min(end, this.pagination.total)
    },
    visiblePages() {
      const current = this.pagination.page
      const total = this.pagination.total_pages
      const visible = []

      if (total <= 7) {
        // Show all pages if total is small
        for (let i = 1; i <= total; i++) {
          visible.push(i)
        }
      } else {
        // Show smart pagination with ellipsis
        visible.push(1)
        
        if (current > 3) {
          visible.push('...')
        }
        
        for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
          visible.push(i)
        }
        
        if (current < total - 2) {
          visible.push('...')
        }
        
        if (total > 1) {
          visible.push(total)
        }
      }

      return visible
    }
  },
  watch: {
    'pagination.page_size'(newValue) {
      this.currentPageSize = newValue || 25
    }
  },
  methods: {
    goToPage(page) {
      if (page >= 1 && page <= this.pagination.total_pages && page !== this.pagination.page) {
        this.$emit('page-changed', page)
      }
    },
    changePageSize() {
      this.$emit('page-size-changed', parseInt(this.currentPageSize))
    }
  }
}
</script>

<style scoped>
.pagination-container {
  border-top: 1px solid #e5e5e5;
  padding-top: 1rem;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.form-select-sm {
  padding: 0.25rem 1.5rem 0.25rem 0.5rem;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .pagination-container {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch !important;
  }
  
  .d-flex:first-child {
    justify-content: center;
  }
  
  .d-flex:last-child {
    justify-content: center;
  }
}
</style>