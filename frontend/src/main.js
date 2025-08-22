import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'

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

const app = createApp(App)
app.use(router)
app.mount('#app')