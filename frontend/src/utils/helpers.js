// Utility functions for the frontend

export const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const formatDateShort = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric'
  })
}

export const getStatusBadgeClass = (status) => {
  const statusMap = {
    'Pass': 'pass',
    'Fail': 'fail', 
    'Blocked': 'blocked',
    'Skip': 'skip',
    'Not Executed': 'not-executed',
    'Active': 'active',
    'Inactive': 'inactive',
    'In Progress': 'in-progress',
    'Completed': 'completed'
  }
  return statusMap[status] || 'secondary'
}

export const getPriorityBadgeClass = (priority) => {
  const priorityMap = {
    'HIGH': 'high',
    'MEDIUM': 'medium',
    'LOW': 'low'
  }
  return priorityMap[priority?.toUpperCase()] || 'medium'
}

export const truncateText = (text, maxLength = 30) => {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

export const showAlert = (message, type = 'success') => {
  const alertContainer = document.getElementById('alert-container')
  if (!alertContainer) return

  const alertId = 'alert-' + Date.now()
  const alertHtml = `
    <div id="${alertId}" class="alert alert-${type} alert-dismissible fade show" role="alert">
      ${message}
      <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    </div>
  `
  
  alertContainer.insertAdjacentHTML('beforeend', alertHtml)
  
  // Auto-remove after 5 seconds
  setTimeout(() => {
    const alertElement = document.getElementById(alertId)
    if (alertElement) {
      alertElement.remove()
    }
  }, 5000)
}

export const showLoading = () => {
  return `
    <div class="loading">
      <div class="spinner"></div>
    </div>
  `
}

export const calculateStats = (items, statusField = 'status') => {
  const stats = {
    total: items.length,
    pass: 0,
    fail: 0,
    blocked: 0,
    skip: 0,
    notExecuted: 0
  }

  items.forEach(item => {
    const status = item[statusField]?.toLowerCase()
    if (status === 'pass') stats.pass++
    else if (status === 'fail') stats.fail++
    else if (status === 'blocked') stats.blocked++
    else if (status === 'skip') stats.skip++
    else stats.notExecuted++
  })

  return stats
}