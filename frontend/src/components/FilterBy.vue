<template>
  <div class="filter-by-container d-flex gap-2 align-items-center">
    <div v-for="filter in filterOptions" :key="filter.key" class="filter-item">
      <select 
        :id="`filter-${filter.key}-${componentId}`"
        class="form-select form-select-sm"
        :value="currentFilters[filter.key] || ''"
        @change="handleFilterChange(filter.key, $event.target.value)"
        :style="{ minWidth: '140px' }"
      >
        <option value="">{{ filter.allLabel || `All ${filter.label}` }}</option>
        <option 
          v-for="option in filter.options" 
          :key="option.value" 
          :value="option.value"
        >
          {{ option.label }}
        </option>
      </select>
      <label 
        :for="`filter-${filter.key}-${componentId}`" 
        class="visually-hidden"
      >
        Filter by {{ filter.label }}
      </label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FilterBy',
  props: {
    filterOptions: {
      type: Array,
      required: true,
      default: () => []
    },
    defaultFilters: {
      type: Object,
      default: () => ({})
    },
    componentId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      currentFilters: { ...this.defaultFilters }
    }
  },
  watch: {
    defaultFilters: {
      handler(newFilters) {
        this.currentFilters = { ...newFilters }
      },
      deep: true,
      immediate: true
    }
  },
  methods: {
    handleFilterChange(filterKey, newValue) {
      // Update the filter value
      if (newValue === '') {
        // Remove filter if empty value selected
        const updatedFilters = { ...this.currentFilters }
        delete updatedFilters[filterKey]
        this.currentFilters = updatedFilters
      } else {
        this.currentFilters = {
          ...this.currentFilters,
          [filterKey]: newValue
        }
      }
      
      // Emit the change event with all current filters
      this.$emit('filters-changed', this.currentFilters)
    },
    
    clearAllFilters() {
      this.currentFilters = {}
      this.$emit('filters-changed', {})
    },
    
    getActiveFilterCount() {
      return Object.keys(this.currentFilters).length
    }
  }
}
</script>

<style scoped>
.filter-by-container {
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
}

.form-select-sm {
  font-size: 0.875rem;
  border-color: #dee2e6;
}

.form-select-sm:focus {
  border-color: #86b7fe;
  box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
}

@media (max-width: 768px) {
  .filter-by-container {
    justify-content: flex-start;
  }
  
  .filter-item .form-select-sm {
    min-width: 120px !important;
  }
}
</style>