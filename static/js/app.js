// Test Management Platform JavaScript

// Global variables
let currentProject = null;
let currentTestSuite = null;
let currentTestCase = null;
let currentTestRun = null;

// Utility functions
function showLoading(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.innerHTML = '<div class="loading"><div class="spinner"></div></div>';
    }
}

function hideLoading() {
    $('.loading').remove();
}

function showAlert(message, type = 'success') {
    const alertHtml = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
    $('#alerts').html(alertHtml);
    setTimeout(() => {
        $('.alert').fadeOut();
    }, 5000);
}

function formatDate(dateString) {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString() + ' ' + new Date(dateString).toLocaleTimeString();
}

function getStatusBadgeClass(status) {
    const statusMap = {
        'Pass': 'status-pass',
        'Fail': 'status-fail',
        'Blocked': 'status-blocked',
        'Skip': 'status-skip',
        'Not Executed': 'status-not-executed',
        'In Progress': 'status-in-progress',
        'Completed': 'status-completed',
        'Active': 'status-active',
        'Inactive': 'status-inactive'
    };
    return statusMap[status] || 'status-not-executed';
}

function getPriorityBadgeClass(priority) {
    const priorityMap = {
        'High': 'priority-high',
        'Medium': 'priority-medium',
        'Low': 'priority-low'
    };
    return priorityMap[priority] || 'priority-medium';
}

function truncateText(text, maxLength = 30) {
    if (!text || text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// API functions
const API = {
    // Projects
    getProjects: () => $.get('/api/projects'),
    createProject: (data) => $.post('/api/projects', JSON.stringify(data), null, 'json').fail(handleApiError),
    getProject: (id) => $.get(`/api/projects/${id}`),
    updateProject: (id, data) => $.ajax({
        url: `/api/projects/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),
    deleteProject: (id) => $.ajax({
        url: `/api/projects/${id}`,
        method: 'DELETE'
    }).fail(handleApiError),

    // Test Suites
    getTestSuites: (projectId) => $.get('/api/test-suites' + (projectId ? `?project_id=${projectId}` : '')),
    createTestSuite: (data) => $.post('/api/test-suites', JSON.stringify(data), null, 'json').fail(handleApiError),
    getTestSuite: (id) => $.get(`/api/test-suites/${id}`),
    updateTestSuite: (id, data) => $.ajax({
        url: `/api/test-suites/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),
    deleteTestSuite: (id) => $.ajax({
        url: `/api/test-suites/${id}`,
        method: 'DELETE'
    }).fail(handleApiError),

    // Test Cases
    getTestCases: (testSuiteId) => $.get('/api/test-cases' + (testSuiteId ? `?test_suite_id=${testSuiteId}` : '')),
    createTestCase: (data) => $.post('/api/test-cases', JSON.stringify(data), null, 'json').fail(handleApiError),
    getTestCase: (id) => $.get(`/api/test-cases/${id}`),
    updateTestCase: (id, data) => $.ajax({
        url: `/api/test-cases/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),
    deleteTestCase: (id) => $.ajax({
        url: `/api/test-cases/${id}`,
        method: 'DELETE'
    }).fail(handleApiError),
    searchTestCases: (query) => $.get(`/api/test-cases/search?q=${encodeURIComponent(query)}`),

    // Test Steps
    getTestSteps: (testCaseId) => $.get(`/api/test-cases/${testCaseId}/steps`),
    createTestStep: (testCaseId, data) => $.post(`/api/test-cases/${testCaseId}/steps`, JSON.stringify(data), null, 'json').fail(handleApiError),
    updateTestStep: (id, data) => $.ajax({
        url: `/api/test-steps/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),
    deleteTestStep: (id) => $.ajax({
        url: `/api/test-steps/${id}`,
        method: 'DELETE'
    }).fail(handleApiError),

    // Test Runs
    getTestRuns: (testSuiteId) => $.get('/api/test-runs' + (testSuiteId ? `?test_suite_id=${testSuiteId}` : '')),
    createTestRun: (data) => $.post('/api/test-runs', JSON.stringify(data), null, 'json').fail(handleApiError),
    getTestRun: (id) => $.get(`/api/test-runs/${id}`),
    updateTestRun: (id, data) => $.ajax({
        url: `/api/test-runs/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),
    deleteTestRun: (id) => $.ajax({
        url: `/api/test-runs/${id}`,
        method: 'DELETE'
    }).fail(handleApiError),

    // Test Executions
    getTestExecutions: (testRunId) => $.get(`/api/test-runs/${testRunId}/executions`),
    updateTestExecution: (id, data) => $.ajax({
        url: `/api/test-executions/${id}`,
        method: 'PUT',
        data: JSON.stringify(data),
        contentType: 'application/json'
    }).fail(handleApiError),

    // Reports
    getReportSummary: () => $.get('/api/reports/summary'),
    exportReport: (format) => window.open(`/api/reports/export/${format}`)
};

function handleApiError(xhr) {
    let message = 'An error occurred';
    if (xhr.responseJSON && xhr.responseJSON.error) {
        message = xhr.responseJSON.error;
    }
    showAlert(message, 'danger');
}

// Main application functions
function loadProjects() {
    showLoading('projects-container');
    API.getProjects().then(projects => {
        hideLoading();
        renderProjects(projects);
    });
}

function renderProjects(projects) {
    const container = document.getElementById('projects-container');
    if (!container) {
        console.error('projects-container not found');
        return;
    }
    if (!projects || projects.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-folder-open"></i>
                <h5>No Projects Found</h5>
                <p>Create your first project to get started with test management.</p>
                <button class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#projectModal">
                    <i class="fas fa-plus"></i> Create Project
                </button>
            </div>
        `;
        return;
    }

    let html = '<div class="row">';
    projects.forEach(project => {
        const testSuitesCount = project.test_suites_count || 0;
        html += `
            <div class="col-md-6 col-lg-4 mb-4">
                <div class="card h-100">
                    <div class="card-body">
                        <h5 class="card-title">${project.name}</h5>
                        <p class="card-text">${project.description || 'No description'}</p>
                        <div class="d-flex justify-content-between align-items-center">
                            <small class="text-muted">${testSuitesCount} test suite(s)</small>
                            <div class="action-buttons">
                                <a href="/project/${project.id}" class="btn btn-outline-primary btn-sm">
                                    <i class="fas fa-eye"></i> View
                                </a>
                                <button class="btn btn-outline-secondary btn-sm" onclick="editProject(${project.id})">
                                    <i class="fas fa-edit"></i> Edit
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="card-footer text-muted">
                        <small>Created ${formatDate(project.created_at)}</small>
                    </div>
                </div>
            </div>
        `;
    });
    html += '</div>';
    container.innerHTML = html;
}

function loadProjectDetails(projectId) {
    currentProject = projectId;
    
    // Ensure DOM is ready before proceeding
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => loadProjectDetails(projectId));
        return;
    }
    
    showLoading('project-details');
    
    Promise.all([
        API.getProject(projectId),
        API.getTestSuites(projectId)
    ]).then(([project, testSuites]) => {
        hideLoading();
        console.log('Loaded project:', project);
        console.log('Loaded test suites:', testSuites);
        
        // Use setTimeout to ensure DOM elements are available
        setTimeout(() => {
            renderProjectDetails(project, testSuites);
        }, 100);
    }).catch(error => {
        hideLoading();
        console.error('Error loading project details:', error);
        showAlert('Error loading project details', 'danger');
    });
}

function renderProjectDetails(project, testSuites) {
    console.log('Rendering project details:', project, testSuites);
    
    const projectNameEl = document.getElementById('project-name');
    const projectDescEl = document.getElementById('project-description');
    const projectBreadcrumb = document.getElementById('project-breadcrumb');
    
    if (projectNameEl) projectNameEl.textContent = project.name;
    if (projectDescEl) projectDescEl.textContent = project.description || 'No description';
    if (projectBreadcrumb) projectBreadcrumb.textContent = truncateText(project.name);
    
    // Try multiple ways to find the container
    let container = document.getElementById('test-suites-container');
    
    if (!container) {
        // Wait and try again
        setTimeout(() => {
            container = document.getElementById('test-suites-container');
            if (!container) {
                console.error('test-suites-container still not found after delay');
                // Try creating the element if it doesn't exist
                const projectDetails = document.getElementById('project-details');
                if (projectDetails) {
                    container = document.createElement('div');
                    container.id = 'test-suites-container';
                    container.innerHTML = '<h3><i class="fas fa-tasks"></i> Test Suites</h3>';
                    projectDetails.appendChild(container);
                    console.log('Created test-suites-container element');
                }
            }
            if (container) {
                renderTestSuitesInContainer(container, testSuites);
            }
        }, 200);
        return;
    }
    
    renderTestSuitesInContainer(container, testSuites);
}

function renderTestSuitesInContainer(container, testSuites) {
    if (!testSuites || testSuites.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-tasks"></i>
                <h5>No Test Suites Found</h5>
                <p>Create your first test suite to organize test cases.</p>
            </div>
        `;
        return;
    }

    let html = '<div class="row">';
    testSuites.forEach(testSuite => {
        const testCasesCount = testSuite.test_cases_count || 0;
        html += `
            <div class="col-md-6 col-lg-4 mb-4">
                <div class="card h-100">
                    <div class="card-body">
                        <h5 class="card-title">${testSuite.name}</h5>
                        <p class="card-text">${testSuite.description || 'No description'}</p>
                        <div class="d-flex justify-content-between align-items-center">
                            <small class="text-muted">${testCasesCount} test case(s)</small>
                            <div class="action-buttons">
                                <a href="/project/${currentProject}/test-suite/${testSuite.id}" class="btn btn-outline-primary btn-sm">
                                    <i class="fas fa-eye"></i> View
                                </a>
                                <button class="btn btn-outline-success btn-sm" onclick="createTestRun(${testSuite.id})">
                                    <i class="fas fa-play"></i> Run Tests
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="card-footer text-muted">
                        <small>Created ${formatDate(testSuite.created_at)}</small>
                    </div>
                </div>
            </div>
        `;
    });
    html += '</div>';
    container.innerHTML = html;
}

function loadTestSuiteDetails(testSuiteId) {
    currentTestSuite = testSuiteId;
    
    // Ensure DOM is ready before proceeding
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => loadTestSuiteDetails(testSuiteId));
        return;
    }
    
    showLoading('test-suite-details');
    
    Promise.all([
        API.getTestSuite(testSuiteId),
        API.getTestCases(testSuiteId),
        API.getTestRuns(testSuiteId)
    ]).then(([testSuite, testCases, testRuns]) => {
        hideLoading();
        
        // Use setTimeout to ensure DOM elements are available
        setTimeout(() => {
            renderTestSuiteDetails(testSuite, testCases, testRuns);
        }, 100);
    }).catch(error => {
        hideLoading();
        console.error('Error loading test suite details:', error);
        showAlert('Error loading test suite details', 'danger');
    });
}

function renderTestSuiteDetails(testSuite, testCases, testRuns) {
    console.log('Rendering test suite details:', testSuite, testCases, testRuns);
    
    // Update basic information with null checks
    const testSuiteNameEl = document.getElementById('test-suite-name');
    const testSuiteDescEl = document.getElementById('test-suite-description');
    const projectNameEl = document.getElementById('project-name');
    const projectBreadcrumb = document.getElementById('project-breadcrumb');
    const testSuiteBreadcrumb = document.getElementById('test-suite-breadcrumb');
    
    if (testSuiteNameEl) testSuiteNameEl.textContent = testSuite.name;
    if (testSuiteDescEl) testSuiteDescEl.textContent = testSuite.description || 'No description';
    if (projectNameEl) projectNameEl.textContent = testSuite.project.name;
    if (projectBreadcrumb) projectBreadcrumb.textContent = truncateText(testSuite.project.name);
    if (testSuiteBreadcrumb) testSuiteBreadcrumb.textContent = truncateText(testSuite.name);
    
    // Render test cases with fallback creation
    let casesContainer = document.getElementById('test-cases-container');
    if (!casesContainer) {
        console.log('test-cases-container not found, creating fallback');
        const mainContent = document.querySelector('main');
        if (mainContent) {
            // Create section with proper Bootstrap styling
            const section = document.createElement('div');
            section.className = 'mb-4';
            section.innerHTML = `
                <div class="mb-3">
                    <h3>Test Cases</h3>
                </div>
                <div id="test-cases-container"></div>
            `;
            mainContent.appendChild(section);
            casesContainer = document.getElementById('test-cases-container');
            console.log('Created test-cases-container element');
        }
    }
    
    if (casesContainer && testCases && testCases.length > 0) {
        console.log('Calling renderTestCasesTable with container:', casesContainer);
        renderTestCasesTable(testCases, casesContainer);
    } else if (casesContainer) {
        casesContainer.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-clipboard-list"></i>
                <h5>No Test Cases Found</h5>
                <p>Create your first test case to start testing.</p>
            </div>
        `;
    }
    
    // Render test runs with fallback creation
    let runsContainer = document.getElementById('test-runs-container');
    if (!runsContainer) {
        console.log('test-runs-container not found, creating fallback');
        const mainContent = document.querySelector('main');
        if (mainContent) {
            // Create section with proper Bootstrap styling
            const section = document.createElement('div');
            section.className = 'mb-4';
            section.innerHTML = `
                <div class="mb-3">
                    <h3>Test Runs</h3>
                </div>
                <div id="test-runs-container"></div>
            `;
            mainContent.appendChild(section);
            runsContainer = document.getElementById('test-runs-container');
            console.log('Created test-runs-container element');
        }
    }
    
    if (runsContainer && testRuns && testRuns.length > 0) {
        renderTestRunsTable(testRuns, runsContainer);
    } else if (runsContainer) {
        runsContainer.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-play-circle"></i>
                <h5>No Test Runs Found</h5>
                <p>Execute your first test run to track progress.</p>
            </div>
        `;
    }
}

function renderTestCasesTable(testCases, container) {
    if (!container) {
        console.error('Container not provided for test cases table');
        return;
    }
    
    let html = `
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Title</th>
                        <th>Priority</th>
                        <th>Status</th>
                        <th>Steps</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    testCases.forEach(testCase => {
        const stepsCount = testCase.test_steps_count || 0;
        html += `
            <tr>
                <td>
                    <strong>${testCase.title}</strong><br>
                    <small class="text-muted">${testCase.description || 'No description'}</small>
                </td>
                <td>
                    <span class="status-badge ${getPriorityBadgeClass(testCase.priority)}">${testCase.priority}</span>
                </td>
                <td>
                    <span class="status-badge ${getStatusBadgeClass(testCase.status)}">${testCase.status}</span>
                </td>
                <td>${stepsCount}</td>
                <td>
                    <div class="action-buttons">
                        <a href="/project/${currentProject}/test-suite/${currentTestSuite}/test-case/${testCase.id}" class="btn btn-outline-primary btn-sm">
                            <i class="fas fa-eye"></i> View
                        </a>
                        <button class="btn btn-outline-secondary btn-sm" onclick="editTestCase(${testCase.id})">
                            <i class="fas fa-edit"></i> Edit
                        </button>
                        <button class="btn btn-outline-danger btn-sm" onclick="deleteTestCase(${testCase.id}, '${testCase.title}')">
                            <i class="fas fa-trash"></i> Delete
                        </button>
                    </div>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    container.innerHTML = html;
}

function renderTestRunsTable(testRuns, container) {
    let html = `
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Status</th>
                        <th>Started</th>
                        <th>Completed</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    testRuns.forEach(testRun => {
        html += `
            <tr>
                <td>
                    <strong>${testRun.name}</strong><br>
                    <small class="text-muted">${testRun.description || 'No description'}</small>
                </td>
                <td>
                    <span class="status-badge ${getStatusBadgeClass(testRun.status)}">${testRun.status}</span>
                </td>
                <td>${formatDate(testRun.started_at)}</td>
                <td>${formatDate(testRun.completed_at)}</td>
                <td>
                    <div class="action-buttons">
                        <a href="/test-run/${testRun.id}" class="btn btn-outline-primary btn-sm">
                            <i class="fas fa-eye"></i> View
                        </a>
                    </div>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    container.innerHTML = html;
}

function loadTestCaseDetails(testCaseId) {
    currentTestCase = testCaseId;
    
    // Ensure DOM is ready before proceeding
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => loadTestCaseDetails(testCaseId));
        return;
    }
    
    showLoading('test-case-details');
    
    Promise.all([
        API.getTestCase(testCaseId),
        API.getTestSteps(testCaseId)
    ]).then(([testCase, testSteps]) => {
        hideLoading();
        
        // Use setTimeout to ensure DOM elements are available
        setTimeout(() => {
            renderTestCaseDetails(testCase, testSteps);
        }, 100);
    }).catch(error => {
        hideLoading();
        console.error('Error loading test case details:', error);
        showAlert('Error loading test case details', 'danger');
    });
}

function renderTestCaseDetails(testCase, testSteps) {
    console.log('Starting renderTestCaseDetails:', testCase, testSteps);
    
    try {
        // Update basic information with null checks
        console.log('Getting DOM elements...');
        const titleEl = document.getElementById('test-case-title');
        const descEl = document.getElementById('test-case-description');
        const priorityEl = document.getElementById('test-case-priority');
        const statusEl = document.getElementById('test-case-status');
        const testSuiteNameEl = document.getElementById('test-suite-name');
        const projectNameEl = document.getElementById('project-name');
        
        console.log('Found elements:', {
            titleEl: !!titleEl,
            descEl: !!descEl,
            priorityEl: !!priorityEl,
            statusEl: !!statusEl,
            testSuiteNameEl: !!testSuiteNameEl,
            projectNameEl: !!projectNameEl
        });
        
        console.log('Setting element values...');
        if (titleEl) {
            titleEl.textContent = testCase.title;
            console.log('Set title:', testCase.title);
        } else {
            console.log('test-case-title element not found');
        }
        
        if (descEl) {
            descEl.textContent = testCase.description || 'No description';
            console.log('Set description');
        } else {
            console.log('test-case-description element not found');
        }
        
        if (priorityEl) {
            priorityEl.innerHTML = `<span class="status-badge ${getPriorityBadgeClass(testCase.priority)}">${testCase.priority}</span>`;
            console.log('Set priority');
        } else {
            console.log('test-case-priority element not found');
        }
        
        if (statusEl) {
            statusEl.innerHTML = `<span class="status-badge ${getStatusBadgeClass(testCase.status)}">${testCase.status}</span>`;
            console.log('Set status');
        } else {
            console.log('test-case-status element not found');
        }
        
        if (testSuiteNameEl) {
            testSuiteNameEl.textContent = testCase.test_suite.name;
            console.log('Set test suite name');
        } else {
            console.log('test-suite-name element not found');
        }
        
        if (projectNameEl) {
            projectNameEl.textContent = testCase.test_suite.project.name;
            console.log('Set project name');
        } else {
            console.log('project-name element not found');
        }
        
        // Update breadcrumbs with actual names
        const projectBreadcrumb = document.getElementById('project-breadcrumb');
        const testSuiteBreadcrumb = document.getElementById('test-suite-breadcrumb');
        if (projectBreadcrumb) projectBreadcrumb.textContent = truncateText(testCase.test_suite.project.name);
        if (testSuiteBreadcrumb) testSuiteBreadcrumb.textContent = truncateText(testCase.test_suite.name);
        
        console.log('Finished setting basic info');
    } catch (error) {
        console.error('Error in renderTestCaseDetails:', error);
        console.error('Error stack:', error.stack);
    }
    
    // Render test steps with fallback creation
    let stepsContainer = document.getElementById('test-steps-container');
    if (!stepsContainer) {
        console.log('test-steps-container not found, creating fallback');
        const mainContent = document.querySelector('main');
        if (mainContent) {
            const section = document.createElement('div');
            section.className = 'mb-4';
            section.innerHTML = `
                <div class="mb-3">
                    <h3>Test Steps</h3>
                </div>
                <div id="test-steps-container"></div>
            `;
            mainContent.appendChild(section);
            stepsContainer = document.getElementById('test-steps-container');
            console.log('Created test-steps-container element');
        }
    }
    
    if (stepsContainer) {
        if (!testSteps || testSteps.length === 0) {
            stepsContainer.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-list-ol"></i>
                    <h5>No Test Steps Found</h5>
                    <p>Add test steps to define the testing procedure.</p>
                </div>
            `;
        } else {
            renderTestSteps(testSteps, stepsContainer);
        }
    }
}

function renderTestSteps(testSteps, container) {
    if (!container) {
        console.error('Container not provided for test steps');
        return;
    }
    
    console.log('Rendering test steps:', testSteps, 'in container:', container);
    
    let html = '';
    testSteps.sort((a, b) => a.step_number - b.step_number);
    
    testSteps.forEach(step => {
        html += `
            <div class="test-step">
                <div class="d-flex align-items-start">
                    <div class="test-step-number">${step.step_number}</div>
                    <div class="flex-grow-1">
                        <h6 class="mb-2">Step ${step.step_number}</h6>
                        <div class="mb-2">
                            <strong>Description:</strong><br>
                            ${step.description}
                        </div>
                        <div class="mb-2">
                            <strong>Expected Result:</strong><br>
                            ${step.expected_result}
                        </div>
                        <div class="action-buttons">
                            <button class="btn btn-outline-secondary btn-sm" onclick="editTestStep(${step.id})">
                                <i class="fas fa-edit"></i> Edit
                            </button>
                            <button class="btn btn-outline-danger btn-sm" onclick="deleteTestStep(${step.id})">
                                <i class="fas fa-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
    });
    
    try {
        container.innerHTML = html;
        console.log('Successfully rendered test steps');
    } catch (error) {
        console.error('Error setting innerHTML for test steps:', error);
    }
}

function loadTestRunDetails(testRunId) {
    currentTestRun = testRunId;
    showLoading('test-run-details');
    
    Promise.all([
        API.getTestRun(testRunId),
        API.getTestExecutions(testRunId)
    ]).then(([testRun, testExecutions]) => {
        hideLoading();
        renderTestRunDetails(testRun, testExecutions);
    });
}

function renderTestRunDetails(testRun, testExecutions) {
    document.getElementById('test-run-name').textContent = testRun.name;
    document.getElementById('test-run-description').textContent = testRun.description || 'No description';
    document.getElementById('test-run-status').innerHTML = `<span class="status-badge ${getStatusBadgeClass(testRun.status)}">${testRun.status}</span>`;
    document.getElementById('test-suite-name').textContent = testRun.test_suite.name;
    document.getElementById('project-name').textContent = testRun.test_suite.project.name;
    
    // Update breadcrumbs with actual names
    const projectBreadcrumb = document.getElementById('project-breadcrumb');
    const testSuiteBreadcrumb = document.getElementById('test-suite-breadcrumb');
    if (projectBreadcrumb) projectBreadcrumb.textContent = truncateText(testRun.test_suite.project.name);
    if (testSuiteBreadcrumb) testSuiteBreadcrumb.textContent = truncateText(testRun.test_suite.name);
    document.getElementById('started-at').textContent = formatDate(testRun.started_at);
    document.getElementById('completed-at').textContent = formatDate(testRun.completed_at);
    
    // Calculate statistics
    const stats = {
        total: testExecutions.length,
        passed: testExecutions.filter(e => e.status === 'Pass').length,
        failed: testExecutions.filter(e => e.status === 'Fail').length,
        blocked: testExecutions.filter(e => e.status === 'Blocked').length,
        skipped: testExecutions.filter(e => e.status === 'Skip').length,
        notExecuted: testExecutions.filter(e => e.status === 'Not Executed').length
    };
    
    const passRate = stats.total > 0 ? Math.round((stats.passed / stats.total) * 100) : 0;
    
    document.getElementById('execution-stats').innerHTML = `
        <div class="row">
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-primary">${stats.total}</h4>
                    <small>Total</small>
                </div>
            </div>
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-success">${stats.passed}</h4>
                    <small>Passed</small>
                </div>
            </div>
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-danger">${stats.failed}</h4>
                    <small>Failed</small>
                </div>
            </div>
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-warning">${stats.blocked}</h4>
                    <small>Blocked</small>
                </div>
            </div>
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-info">${stats.skipped}</h4>
                    <small>Skipped</small>
                </div>
            </div>
            <div class="col-md-2">
                <div class="text-center">
                    <h4 class="mb-0 text-secondary">${stats.notExecuted}</h4>
                    <small>Not Executed</small>
                </div>
            </div>
        </div>
        <div class="mt-3">
            <div class="d-flex justify-content-between align-items-center mb-2">
                <span>Pass Rate</span>
                <span class="fw-bold">${passRate}%</span>
            </div>
            <div class="progress" style="height: 10px;">
                <div class="progress-bar bg-success" style="width: ${passRate}%"></div>
            </div>
        </div>
    `;
    
    renderTestExecutions(testExecutions);
}

function renderTestExecutions(testExecutions) {
    const container = document.getElementById('test-executions-container');
    
    if (!testExecutions || testExecutions.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-clipboard-check"></i>
                <h5>No Test Executions Found</h5>
                <p>Test executions will appear here when tests are run.</p>
            </div>
        `;
        return;
    }
    
    let html = `
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Test Case</th>
                        <th>Status</th>
                        <th>Executed By</th>
                        <th>Executed At</th>
                        <th>Notes</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    testExecutions.forEach(execution => {
        html += `
            <tr>
                <td>
                    <strong>${execution.test_case.title}</strong><br>
                    <small class="text-muted">${execution.test_case.description || ''}</small>
                </td>
                <td>
                    <span class="status-badge ${getStatusBadgeClass(execution.status)}">${execution.status}</span>
                </td>
                <td>${execution.executed_by || 'N/A'}</td>
                <td>${formatDate(execution.executed_at)}</td>
                <td>${execution.notes || 'N/A'}</td>
                <td>
                    <button class="btn btn-outline-primary btn-sm" onclick="showExecutionModal(${execution.id})">
                        <i class="fas fa-edit"></i> Update
                    </button>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    container.innerHTML = html;
}

function loadReports() {
    showLoading('reports-container');
    API.getReportSummary().then(summary => {
        hideLoading();
        renderReports(summary);
    });
}

function renderReports(summary) {
    // Overview statistics
    document.getElementById('overview-stats').innerHTML = `
        <div class="row">
            <div class="col-md-3">
                <div class="stats-card">
                    <div class="stats-number">${summary.total_projects}</div>
                    <div class="stats-label">Projects</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="stats-card">
                    <div class="stats-number">${summary.total_test_suites}</div>
                    <div class="stats-label">Test Suites</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="stats-card">
                    <div class="stats-number">${summary.total_test_cases}</div>
                    <div class="stats-label">Test Cases</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="stats-card">
                    <div class="stats-number">${summary.total_test_runs}</div>
                    <div class="stats-label">Test Runs</div>
                </div>
            </div>
        </div>
    `;
    
    // Execution summary
    const execSummary = summary.execution_summary;
    document.getElementById('execution-summary').innerHTML = `
        <div class="row">
            <div class="col-md-8">
                <div class="row text-center">
                    <div class="col">
                        <h4 class="text-success">${execSummary.passed}</h4>
                        <small>Passed</small>
                    </div>
                    <div class="col">
                        <h4 class="text-danger">${execSummary.failed}</h4>
                        <small>Failed</small>
                    </div>
                    <div class="col">
                        <h4 class="text-warning">${execSummary.blocked}</h4>
                        <small>Blocked</small>
                    </div>
                    <div class="col">
                        <h4 class="text-info">${execSummary.skipped}</h4>
                        <small>Skipped</small>
                    </div>
                    <div class="col">
                        <h4 class="text-secondary">${execSummary.not_executed}</h4>
                        <small>Not Executed</small>
                    </div>
                </div>
            </div>
            <div class="col-md-4 text-center">
                <h2 class="text-primary">${Math.round(execSummary.pass_rate)}%</h2>
                <small>Pass Rate</small>
                <div class="progress mt-2">
                    <div class="progress-bar bg-success" style="width: ${execSummary.pass_rate}%"></div>
                </div>
            </div>
        </div>
    `;
    
    // Recent runs
    const recentContainer = document.getElementById('recent-runs');
    if (summary.recent_runs && summary.recent_runs.length > 0) {
        let html = `
            <div class="table-responsive">
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th>Test Run</th>
                            <th>Project</th>
                            <th>Status</th>
                            <th>Started</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
        `;
        
        summary.recent_runs.forEach(run => {
            html += `
                <tr>
                    <td>
                        <strong>${run.name}</strong><br>
                        <small class="text-muted">${run.description || ''}</small>
                    </td>
                    <td>${run.test_suite.project.name}</td>
                    <td>
                        <span class="status-badge ${getStatusBadgeClass(run.status)}">${run.status}</span>
                    </td>
                    <td>${formatDate(run.started_at)}</td>
                    <td>
                        <a href="/test-run/${run.id}" class="btn btn-outline-primary btn-sm">
                            <i class="fas fa-eye"></i> View
                        </a>
                    </td>
                </tr>
            `;
        });
        
        html += `
                    </tbody>
                </table>
            </div>
        `;
        recentContainer.innerHTML = html;
    } else {
        recentContainer.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-chart-bar"></i>
                <h5>No Recent Test Runs</h5>
                <p>Recent test runs will appear here once you start executing tests.</p>
            </div>
        `;
    }
}

// Search functionality
function setupSearch() {
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        let searchTimeout;
        searchInput.addEventListener('input', function() {
            clearTimeout(searchTimeout);
            const query = this.value.trim();
            
            if (query.length < 2) {
                document.getElementById('search-results').innerHTML = '';
                return;
            }
            
            searchTimeout = setTimeout(() => {
                performSearch(query);
            }, 300);
        });
    }
}

function performSearch(query) {
    showLoading('search-results');
    API.searchTestCases(query).then(results => {
        hideLoading();
        renderSearchResults(results, query);
    });
}

function renderSearchResults(results, query) {
    const container = document.getElementById('search-results');
    
    if (!results || results.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-search"></i>
                <h5>No Results Found</h5>
                <p>No test cases found matching "${query}"</p>
            </div>
        `;
        return;
    }
    
    let html = `
        <h5>Search Results for "${query}" (${results.length} found)</h5>
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Test Case</th>
                        <th>Project</th>
                        <th>Test Suite</th>
                        <th>Priority</th>
                        <th>Status</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    results.forEach(testCase => {
        html += `
            <tr>
                <td>
                    <strong>${testCase.title}</strong><br>
                    <small class="text-muted">${testCase.description || 'No description'}</small>
                </td>
                <td>${testCase.test_suite.project.name}</td>
                <td>${testCase.test_suite.name}</td>
                <td>
                    <span class="status-badge ${getPriorityBadgeClass(testCase.priority)}">${testCase.priority}</span>
                </td>
                <td>
                    <span class="status-badge ${getStatusBadgeClass(testCase.status)}">${testCase.status}</span>
                </td>
                <td>
                    <a href="/test-case/${testCase.id}" class="btn btn-outline-primary btn-sm">
                        <i class="fas fa-eye"></i> View
                    </a>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    container.innerHTML = html;
}

// Modal functions
function showCreateProjectModal() {
    $('#projectModal').modal('show');
    document.getElementById('projectForm').reset();
    document.getElementById('projectModalTitle').textContent = 'Create New Project';
    document.getElementById('projectId').value = '';
}

function showCreateTestSuiteModal() {
    $('#testSuiteModal').modal('show');
    document.getElementById('testSuiteForm').reset();
    document.getElementById('testSuiteModalTitle').textContent = 'Create New Test Suite';
    document.getElementById('testSuiteId').value = '';
    document.getElementById('projectId').value = currentProject || '';
}

function showCreateTestCaseModal() {
    $('#testCaseModal').modal('show');
    document.getElementById('testCaseForm').reset();
    document.getElementById('testCaseModalTitle').textContent = 'Create New Test Case';
    document.getElementById('testCaseId').value = '';
    document.getElementById('testSuiteId').value = currentTestSuite || '';
}

function showCreateTestStepModal() {
    $('#testStepModal').modal('show');
    document.getElementById('testStepForm').reset();
    document.getElementById('testStepModalTitle').textContent = 'Add New Test Step';
    document.getElementById('testStepId').value = '';
    
    // Auto-calculate next step number
    API.getTestSteps(currentTestCase).then(testSteps => {
        let nextStepNumber = 1;
        if (testSteps && testSteps.length > 0) {
            const maxStepNumber = Math.max(...testSteps.map(step => step.step_number));
            nextStepNumber = maxStepNumber + 1;
        }
        document.getElementById('testStepNumber').value = nextStepNumber;
    }).catch(error => {
        console.error('Error getting test steps for auto-numbering:', error);
        // Fallback to step number 1 if there's an error
        document.getElementById('testStepNumber').value = 1;
    });
}

function showCreateTestRunModal() {
    $('#testRunModal').modal('show');
    document.getElementById('testRunForm').reset();
    document.getElementById('testRunModalTitle').textContent = 'Create New Test Run';
    document.getElementById('testRunId').value = '';
    document.getElementById('testSuiteId').value = currentTestSuite || '';
}

function showExecutionModal(executionId) {
    // Load execution details and show modal
    $('#executionModal').modal('show');
    document.getElementById('executionId').value = executionId;
}

function editProject(projectId) {
    API.getProject(projectId).then(project => {
        $('#projectModal').modal('show');
        
        // Safely set element content if elements exist
        const titleEl = document.getElementById('projectModalTitle');
        const idEl = document.getElementById('projectId');
        const nameEl = document.getElementById('projectName');
        const descEl = document.getElementById('projectDescription');
        
        if (titleEl) titleEl.textContent = 'Edit Project';
        if (idEl) idEl.value = project.id;
        if (nameEl) nameEl.value = project.name;
        if (descEl) descEl.value = project.description || '';
    }).catch(error => {
        console.error('Error loading project:', error);
        showAlert('Error loading project details', 'danger');
    });
}

function editTestCase(testCaseId) {
    API.getTestCase(testCaseId).then(testCase => {
        $('#testCaseModal').modal('show');
        document.getElementById('testCaseModalTitle').textContent = 'Edit Test Case';
        document.getElementById('testCaseId').value = testCase.id;
        document.getElementById('testCaseTitle').value = testCase.title;
        document.getElementById('testCaseDescription').value = testCase.description || '';
        document.getElementById('testCasePriority').value = testCase.priority;
        document.getElementById('testSuiteId').value = testCase.test_suite_id;
    });
}

function createTestRun(testSuiteId) {
    currentTestSuite = testSuiteId;
    showCreateTestRunModal();
}

// Form submission handlers
function handleProjectForm() {
    const form = document.getElementById('projectForm');
    const formData = new FormData(form);
    const data = {
        name: formData.get('name'),
        description: formData.get('description')
    };
    
    const projectId = document.getElementById('projectId').value;
    
    const apiCall = projectId ? 
        API.updateProject(projectId, data) : 
        API.createProject(data);
    
    apiCall.then(() => {
        $('#projectModal').modal('hide');
        showAlert(projectId ? 'Project updated successfully' : 'Project created successfully');
        if (typeof loadProjects === 'function') loadProjects();
        if (typeof loadProjectDetails === 'function' && currentProject) loadProjectDetails(currentProject);
    });
    
    return false;
}

function handleTestSuiteForm() {
    const form = document.getElementById('testSuiteForm');
    const formData = new FormData(form);
    const data = {
        name: formData.get('name'),
        description: formData.get('description'),
        project_id: parseInt(formData.get('project_id'))
    };
    
    const testSuiteId = document.getElementById('testSuiteId').value;
    
    const apiCall = testSuiteId ? 
        API.updateTestSuite(testSuiteId, data) : 
        API.createTestSuite(data);
    
    apiCall.then(() => {
        $('#testSuiteModal').modal('hide');
        showAlert(testSuiteId ? 'Test suite updated successfully' : 'Test suite created successfully');
        if (typeof loadProjectDetails === 'function' && currentProject) loadProjectDetails(currentProject);
    });
    
    return false;
}

function handleTestCaseForm() {
    const form = document.getElementById('testCaseForm');
    const formData = new FormData(form);
    const data = {
        title: formData.get('title'),
        description: formData.get('description'),
        priority: formData.get('priority'),
        test_suite_id: parseInt(formData.get('test_suite_id'))
    };
    
    const testCaseId = document.getElementById('testCaseId').value;
    
    const apiCall = testCaseId ? 
        API.updateTestCase(testCaseId, data) : 
        API.createTestCase(data);
    
    apiCall.then(() => {
        $('#testCaseModal').modal('hide');
        showAlert(testCaseId ? 'Test case updated successfully' : 'Test case created successfully');
        if (typeof loadTestSuiteDetails === 'function' && currentTestSuite) loadTestSuiteDetails(currentTestSuite);
    });
    
    return false;
}

function handleTestStepForm() {
    event.preventDefault();
    
    const form = document.getElementById('testStepForm');
    const formData = new FormData(form);
    const data = {
        step_number: parseInt(formData.get('step_number')),
        description: formData.get('description'),
        expected_result: formData.get('expected_result')
    };
    
    const testStepId = document.getElementById('testStepId').value;
    
    const apiCall = testStepId ? 
        API.updateTestStep(testStepId, data) : 
        API.createTestStep(currentTestCase, data);
    
    apiCall.then(() => {
        $('#testStepModal').modal('hide');
        showAlert(testStepId ? 'Test step updated successfully' : 'Test step created successfully', 'success');
        if (typeof loadTestCaseDetails === 'function' && currentTestCase) {
            loadTestCaseDetails(currentTestCase);
        }
    }).catch(error => {
        console.error('Error saving test step:', error);
        showAlert('Error saving test step', 'danger');
    });
    
    return false;
}

function editTestStep(testStepId) {
    // Find the test step data in currently loaded steps or fetch it
    API.getTestSteps(currentTestCase).then(testSteps => {
        const testStep = testSteps.find(step => step.id === testStepId);
        if (testStep) {
            document.getElementById('testStepModalTitle').textContent = 'Edit Test Step';
            document.getElementById('testStepId').value = testStep.id;
            document.getElementById('testStepNumber').value = testStep.step_number;
            document.getElementById('testStepDescription').value = testStep.description;
            document.getElementById('testStepExpectedResult').value = testStep.expected_result;
            $('#testStepModal').modal('show');
        }
    }).catch(error => {
        console.error('Error loading test step for editing:', error);
        showAlert('Error loading test step details', 'danger');
    });
}

function deleteTestStep(testStepId) {
    if (confirm('Are you sure you want to delete this test step?')) {
        API.deleteTestStep(testStepId).then(() => {
            showAlert('Test step deleted successfully', 'success');
            if (currentTestCase) {
                loadTestCaseDetails(currentTestCase);
            }
        }).catch(error => {
            console.error('Error deleting test step:', error);
            showAlert('Error deleting test step', 'danger');
        });
    }
}

function handleTestRunForm() {
    const form = document.getElementById('testRunForm');
    const formData = new FormData(form);
    const data = {
        name: formData.get('name'),
        description: formData.get('description'),
        test_suite_id: parseInt(formData.get('test_suite_id'))
    };
    
    API.createTestRun(data).then(testRun => {
        $('#testRunModal').modal('hide');
        showAlert('Test run created successfully');
        window.location.href = `/test-run/${testRun.id}`;
    });
    
    return false;
}

function handleExecutionForm() {
    const form = document.getElementById('executionForm');
    const formData = new FormData(form);
    const data = {
        status: formData.get('status'),
        notes: formData.get('notes'),
        executed_by: formData.get('executed_by')
    };
    
    const executionId = document.getElementById('executionId').value;
    
    API.updateTestExecution(executionId, data).then(() => {
        $('#executionModal').modal('hide');
        showAlert('Test execution updated successfully');
        if (typeof loadTestRunDetails === 'function' && currentTestRun) loadTestRunDetails(currentTestRun);
    });
    
    return false;
}

// Test Case Deletion
function deleteTestCase(testCaseId, testCaseTitle) {
    if (confirm(`Are you sure you want to delete the test case "${testCaseTitle}"?\n\nThis will also delete all related test steps and cannot be undone.`)) {
        API.deleteTestCase(testCaseId).then(() => {
            showAlert(`Test case "${testCaseTitle}" deleted successfully`, 'success');
            // Reload the test suite details to update the display
            if (typeof loadTestSuiteDetails === 'function' && currentTestSuite) {
                loadTestSuiteDetails(currentTestSuite);
            }
        }).catch(error => {
            console.error('Error deleting test case:', error);
            showAlert('Error deleting test case', 'danger');
        });
    }
}

// Export functions
function exportReports(format) {
    API.exportReport(format);
}

// Initialize application
$(document).ready(function() {
    // Setup search functionality
    setupSearch();
    
    // Set active navigation
    const path = window.location.pathname;
    $('.nav-link').removeClass('active');
    if (path === '/') {
        $('.nav-link[href="/"]').addClass('active');
    } else if (path.includes('/reports')) {
        $('.nav-link[href="/reports"]').addClass('active');
    }
});
