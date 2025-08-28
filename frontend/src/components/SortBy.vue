<template>
  <div class="d-flex align-items-center">
    <label :for="sortSelectId" class="form-label me-2 mb-0">Sort by:</label>
    <select 
      :id="sortSelectId" 
      class="form-select form-select-sm" 
      :value="currentSort" 
      @change="handleSortChange"
      style="width: auto; min-width: 200px;"
    >
      <option v-for="option in sortOptions" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
  </div>
</template>

<script>
export default {
  name: 'SortBy',
  props: {
    sortOptions: {
      type: Array,
      required: true,
      default: () => []
    },
    defaultSort: {
      type: String,
      default: ''
    },
    // Unique identifier for multiple sort components on same page
    componentId: {
      type: String,
      default: 'default'
    }
  },
  data() {
    return {
      currentSort: this.defaultSort || (this.sortOptions.length > 0 ? this.sortOptions[0].value : '')
    }
  },
  computed: {
    sortSelectId() {
      return `sort-by-${this.componentId}`
    }
  },
  watch: {
    defaultSort(newValue) {
      this.currentSort = newValue
    }
  },
  methods: {
    handleSortChange(event) {
      this.currentSort = event.target.value
      this.$emit('sort-changed', this.currentSort)
    }
  }
}
</script>

<style scoped>
.form-select-sm {
  padding: 0.25rem 1.5rem 0.25rem 0.5rem;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .d-flex {
    flex-direction: column;
    align-items: flex-start !important;
    gap: 0.5rem;
  }
  
  .form-select-sm {
    width: 100% !important;
    min-width: auto !important;
  }
}
</style>