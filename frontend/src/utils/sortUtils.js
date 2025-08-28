/**
 * Utility functions for sorting data
 */

/**
 * Get priority value for sorting
 */
export function getPriorityValue(priority) {
  const priorityMap = {
    'Critical': 4,
    'High': 3,
    'Medium': 2,
    'Low': 1
  }
  return priorityMap[priority] || 0
}

/**
 * Get status value for sorting
 */
export function getStatusValue(status) {
  const statusMap = {
    'Completed': 4,
    'In Progress': 3,
    'Not Started': 2,
    'Cancelled': 1,
    'Pass': 4,
    'Fail': 3,
    'Blocked': 2,
    'Skip': 1,
    'Not Executed': 0
  }
  return statusMap[status] || 0
}

/**
 * Common sort options for different data types
 */
export const SORT_OPTIONS = {
  // Date sorting
  DATE_NEWEST_FIRST: { value: 'created_desc', label: 'Created Date (Newest First)' },
  DATE_OLDEST_FIRST: { value: 'created_asc', label: 'Created Date (Oldest First)' },
  DATE_UPDATED_DESC: { value: 'updated_desc', label: 'Updated Date (Newest First)' },
  DATE_UPDATED_ASC: { value: 'updated_asc', label: 'Updated Date (Oldest First)' },
  
  // Name/Title sorting
  NAME_A_TO_Z: { value: 'name_asc', label: 'Name (A-Z)' },
  NAME_Z_TO_A: { value: 'name_desc', label: 'Name (Z-A)' },
  TITLE_A_TO_Z: { value: 'title_asc', label: 'Title (A-Z)' },
  TITLE_Z_TO_A: { value: 'title_desc', label: 'Title (Z-A)' },
  
  // Priority sorting
  PRIORITY_HIGH_TO_LOW: { value: 'priority_desc', label: 'Priority (Critical to Low)' },
  PRIORITY_LOW_TO_HIGH: { value: 'priority_asc', label: 'Priority (Low to Critical)' },
  
  // Status sorting
  STATUS_DESC: { value: 'status_desc', label: 'Status (High to Low)' },
  STATUS_ASC: { value: 'status_asc', label: 'Status (Low to High)' }
}

/**
 * Generic sort function that handles different data types
 */
export function applySorting(data, sortBy, customSortFunctions = {}) {
  if (!data || !Array.isArray(data)) {
    return []
  }

  const sortedData = [...data]

  // Check if custom sort function exists
  if (customSortFunctions[sortBy]) {
    return sortedData.sort(customSortFunctions[sortBy])
  }

  // Default sorting logic
  switch (sortBy) {
    // Date sorting
    case 'created_asc':
      return sortedData.sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
    case 'created_desc':
      return sortedData.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    case 'updated_asc':
      return sortedData.sort((a, b) => new Date(a.updated_at) - new Date(b.updated_at))
    case 'updated_desc':
      return sortedData.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
      
    // Name/Title sorting
    case 'name_asc':
      return sortedData.sort((a, b) => (a.name || '').localeCompare(b.name || ''))
    case 'name_desc':
      return sortedData.sort((a, b) => (b.name || '').localeCompare(a.name || ''))
    case 'title_asc':
      return sortedData.sort((a, b) => (a.title || '').localeCompare(b.title || ''))
    case 'title_desc':
      return sortedData.sort((a, b) => (b.title || '').localeCompare(a.title || ''))
      
    // Priority sorting
    case 'priority_desc':
      return sortedData.sort((a, b) => getPriorityValue(b.priority) - getPriorityValue(a.priority))
    case 'priority_asc':
      return sortedData.sort((a, b) => getPriorityValue(a.priority) - getPriorityValue(b.priority))
      
    // Status sorting
    case 'status_desc':
      return sortedData.sort((a, b) => getStatusValue(b.status) - getStatusValue(a.status))
    case 'status_asc':
      return sortedData.sort((a, b) => getStatusValue(a.status) - getStatusValue(b.status))
      
    default:
      return sortedData
  }
}

/**
 * Predefined sort option sets for common use cases
 */
export const SORT_OPTION_SETS = {
  // For test cases
  TEST_CASES: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.TITLE_A_TO_Z,
    SORT_OPTIONS.TITLE_Z_TO_A,
    SORT_OPTIONS.PRIORITY_HIGH_TO_LOW,
    SORT_OPTIONS.PRIORITY_LOW_TO_HIGH,
    SORT_OPTIONS.STATUS_DESC,
    SORT_OPTIONS.STATUS_ASC
  ],
  
  // For projects
  PROJECTS: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.NAME_A_TO_Z,
    SORT_OPTIONS.NAME_Z_TO_A
  ],
  
  // For test suites
  TEST_SUITES: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.NAME_A_TO_Z,
    SORT_OPTIONS.NAME_Z_TO_A
  ],
  
  // For test runs
  TEST_RUNS: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.NAME_A_TO_Z,
    SORT_OPTIONS.NAME_Z_TO_A,
    SORT_OPTIONS.STATUS_DESC,
    SORT_OPTIONS.STATUS_ASC
  ],
  
  // For repositories
  REPOSITORIES: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.NAME_A_TO_Z,
    SORT_OPTIONS.NAME_Z_TO_A
  ],
  
  // For authentication keys
  KEYS: [
    SORT_OPTIONS.DATE_NEWEST_FIRST,
    SORT_OPTIONS.DATE_OLDEST_FIRST,
    SORT_OPTIONS.NAME_A_TO_Z,
    SORT_OPTIONS.NAME_Z_TO_A
  ]
}