<template>
  <div>
    <!-- Navigation -->
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
      <div class="container-fluid">
        <router-link to="/" class="navbar-brand">
          <i class="fas fa-clipboard-check"></i> Test Management Platform
        </router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto">
            <li class="nav-item">
              <router-link to="/" class="nav-link" :class="{ active: $route.name === 'Dashboard' }">
                <i class="fas fa-home"></i> Dashboard
              </router-link>
            </li>
            <li class="nav-item">
              <router-link to="/reports" class="nav-link" :class="{ active: $route.name === 'Reports' }">
                <i class="fas fa-chart-bar"></i> Reports
              </router-link>
            </li>
          </ul>
          <div class="d-flex">
            <div class="search-box me-3">
              <input type="text" class="form-control" v-model="searchQuery" placeholder="Search test cases..." @keyup.enter="searchTestCases">
              <i class="fas fa-search search-icon" @click="searchTestCases"></i>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <div class="container-fluid">
      <div class="row">
        <!-- Sidebar -->
        <nav class="col-md-2 d-md-block sidebar">
          <div class="position-sticky pt-3">
            <ul class="nav flex-column">
              <li class="nav-item">
                <router-link to="/" class="nav-link" :class="{ active: $route.name === 'Dashboard' }">
                  <i class="fas fa-tachometer-alt"></i>
                  Dashboard
                </router-link>
              </li>
              <li class="nav-item">
                <router-link to="/reports" class="nav-link" :class="{ active: $route.name === 'Reports' }">
                  <i class="fas fa-chart-bar"></i>
                  Reports
                </router-link>
              </li>
            </ul>
          </div>
        </nav>

        <!-- Main content -->
        <main class="col-md-10 ms-sm-auto col-lg-10 px-md-4">
          <router-view></router-view>
        </main>
      </div>
    </div>

    <!-- Alert container -->
    <div id="alert-container" class="position-fixed top-0 end-0 p-3" style="z-index: 1055;"></div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      searchQuery: ''
    }
  },
  methods: {
    searchTestCases() {
      if (this.searchQuery.trim()) {
        // Implement search functionality
        console.log('Searching for:', this.searchQuery)
      }
    }
  }
}
</script>

<style>
/* Import custom styles */
</style>