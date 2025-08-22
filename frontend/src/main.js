// Use global CDN versions
const { createApp } = Vue
const { createRouter, createWebHistory } = VueRouter

// Import components
import Dashboard from './components/Dashboard.vue'
import ProjectDetail from './components/ProjectDetail.vue'
import TestSuiteDetail from './components/TestSuiteDetail.vue'
import TestCaseDetail from './components/TestCaseDetail.vue'
import TestRunDetail from './components/TestRunDetail.vue'
import Reports from './components/Reports.vue'

// Import styles
import './style.css'

// Router configuration
const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/project/:id', name: 'ProjectDetail', component: ProjectDetail, props: true },
  { path: '/project/:pid/test-suite/:sid', name: 'TestSuiteDetail', component: TestSuiteDetail, props: true },
  { path: '/project/:pid/test-suite/:sid/test-case/:cid', name: 'TestCaseDetail', component: TestCaseDetail, props: true },
  { path: '/test-run/:id', name: 'TestRunDetail', component: TestRunDetail, props: true },
  { path: '/reports', name: 'Reports', component: Reports }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Load App component dynamically since we can't import it
fetch('./src/App.vue')
  .then(response => response.text())
  .then(template => {
    const app = createApp({
      template: `<div id="app">
        <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
          <div class="container-fluid">
            <router-link to="/" class="navbar-brand">
              <i class="fas fa-clipboard-check"></i> Test Management Platform
            </router-link>
          </div>
        </nav>
        <div class="container-fluid">
          <router-view></router-view>
        </div>
      </div>`
    })
    app.use(router)
    app.mount('#app')
  })