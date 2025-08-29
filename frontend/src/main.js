import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'

// Import components
import Dashboard from './components/Dashboard.vue'
import ProjectDetail from './components/ProjectDetail.vue'
import TestSuiteDetail from './components/TestSuiteDetail.vue'
import TestCaseDetail from './components/TestCaseDetail.vue'
import TestRuns from './components/TestRuns.vue'
import TestRunDetail from './components/TestRunDetail.vue'
import TestRunForm from './components/TestRunForm.vue'
import Reports from './components/Reports.vue'
import Repositories from './components/Repositories.vue'
import RepositoryDetail from './components/RepositoryDetail.vue'
import Keys from './components/Keys.vue'

// Import styles
import './style.css'

// Router configuration
const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/project/:id', name: 'ProjectDetail', component: ProjectDetail, props: true },
  { path: '/project/:pid/test-suite/:sid', name: 'TestSuiteDetail', component: TestSuiteDetail, props: true },
  { path: '/project/:pid/test-suite/:sid/test-case/:cid', name: 'TestCaseDetail', component: TestCaseDetail, props: true },
  { path: '/test-runs', name: 'TestRuns', component: TestRuns },
  { path: '/test-runs/new', name: 'TestRunForm', component: TestRunForm },
  { path: '/test-runs/:id/edit', name: 'TestRunEdit', component: TestRunForm, props: true },
  { path: '/test-runs/:id', name: 'TestRunDetail', component: TestRunDetail, props: true },
  { path: '/reports', name: 'Reports', component: Reports },
  { path: '/repositories', name: 'Repositories', component: Repositories },
  { path: '/repository/:id', name: 'RepositoryDetail', component: RepositoryDetail, props: true },
  { path: '/keys', name: 'Keys', component: Keys }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)
app.use(router)
app.mount('#app')