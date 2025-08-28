// Filter utility functions and configurations

/**
 * Apply multiple filters to a dataset
 * @param {Array} data - The data array to filter
 * @param {Object} filters - Object with filter keys and values
 * @param {Object} filterConfig - Configuration for how to apply each filter
 * @returns {Array} - Filtered data array
 */
export function applyFilters(data, filters, filterConfig) {
  if (!data || !Array.isArray(data)) {
    return []
  }

  let filteredData = [...data]

  // Apply each active filter
  Object.entries(filters).forEach(([filterKey, filterValue]) => {
    if (!filterValue || filterValue === '') {
      return // Skip empty filters
    }

    const config = filterConfig[filterKey]
    if (!config) {
      console.warn(`No filter configuration found for key: ${filterKey}`)
      return
    }

    switch (config.type) {
      case 'exact':
        filteredData = filteredData.filter(item => {
          const itemValue = getNestedValue(item, config.field)
          return itemValue === filterValue
        })
        break

      case 'contains':
        filteredData = filteredData.filter(item => {
          const itemValue = getNestedValue(item, config.field)
          return itemValue && itemValue.toLowerCase().includes(filterValue.toLowerCase())
        })
        break

      case 'array_contains':
        filteredData = filteredData.filter(item => {
          const itemArray = getNestedValue(item, config.field)
          return Array.isArray(itemArray) && itemArray.includes(filterValue)
        })
        break

      case 'custom':
        if (typeof config.filterFunction === 'function') {
          filteredData = filteredData.filter(item => config.filterFunction(item, filterValue))
        }
        break

      default:
        console.warn(`Unknown filter type: ${config.type}`)
    }
  })

  return filteredData
}

/**
 * Get nested value from object using dot notation
 * @param {Object} obj - The object to get value from
 * @param {String} path - Dot notation path (e.g., 'user.profile.name')
 * @returns {*} - The value at the path
 */
function getNestedValue(obj, path) {
  if (!obj || !path) return undefined
  
  return path.split('.').reduce((current, key) => {
    return current && current[key] !== undefined ? current[key] : undefined
  }, obj)
}

/**
 * Predefined filter option sets for different components
 */
export const FILTER_OPTION_SETS = {
  TEST_RUNS: [
    {
      key: 'status',
      label: 'Status',
      allLabel: 'All Statuses',
      options: [
        { value: 'Not Started', label: 'Not Started' },
        { value: 'In Progress', label: 'In Progress' },
        { value: 'Completed', label: 'Completed' },
        { value: 'Cancelled', label: 'Cancelled' }
      ]
    }
  ],
  
  TEST_CASES: [
    {
      key: 'status',
      label: 'Status',
      allLabel: 'All Statuses',
      options: [
        { value: 'Pass', label: 'Pass' },
        { value: 'Fail', label: 'Fail' },
        { value: 'Blocked', label: 'Blocked' },
        { value: 'Skip', label: 'Skip' },
        { value: 'Not Executed', label: 'Not Executed' }
      ]
    },
    {
      key: 'priority',
      label: 'Priority',
      allLabel: 'All Priorities',
      options: [
        { value: 'High', label: 'High' },
        { value: 'Medium', label: 'Medium' },
        { value: 'Low', label: 'Low' }
      ]
    }
  ],

  PROJECTS: [
    {
      key: 'status',
      label: 'Status',
      allLabel: 'All Statuses',
      options: [
        { value: 'Active', label: 'Active' },
        { value: 'Inactive', label: 'Inactive' },
        { value: 'Completed', label: 'Completed' }
      ]
    }
  ],

  REPOSITORIES: [
    {
      key: 'has_key',
      label: 'Authentication',
      allLabel: 'All Types',
      options: [
        { value: 'true', label: 'Private (With Key)' },
        { value: 'false', label: 'Public (No Key)' }
      ]
    },
    {
      key: 'synced',
      label: 'Sync Status',
      allLabel: 'All Sync States',
      options: [
        { value: 'true', label: 'Synced' },
        { value: 'false', label: 'Not Synced' }
      ]
    }
  ],

  KEYS: [
    {
      key: 'key_type',
      label: 'Type',
      allLabel: 'All Types',
      options: [
        { value: 'RSA', label: 'RSA Key' },
        { value: 'Username', label: 'Username/Password' }
      ]
    }
  ]
}

/**
 * Predefined filter configurations for different components
 */
export const FILTER_CONFIGS = {
  TEST_RUNS: {
    status: {
      type: 'exact',
      field: 'status'
    }
  },

  TEST_CASES: {
    status: {
      type: 'exact',
      field: 'status'
    },
    priority: {
      type: 'exact',
      field: 'priority'
    }
  },

  PROJECTS: {
    status: {
      type: 'exact',
      field: 'status'
    }
  },

  REPOSITORIES: {
    has_key: {
      type: 'custom',
      filterFunction: (item, filterValue) => {
        const hasKey = item.key_id && item.key_id !== null
        return filterValue === 'true' ? hasKey : !hasKey
      }
    },
    synced: {
      type: 'custom',
      filterFunction: (item, filterValue) => {
        const isSynced = item.synced_at && item.synced_at !== null
        return filterValue === 'true' ? isSynced : !isSynced
      }
    }
  },

  KEYS: {
    key_type: {
      type: 'exact',
      field: 'key_type'
    }
  }
}