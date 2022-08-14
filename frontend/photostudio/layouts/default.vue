<template>
  <a-layout class="layout">
    <a-layout-header>
      <div class="logo" />
      <a-menu
        theme="dark"
        mode="horizontal"
        :style="{ lineHeight: '64px' }"
      >
          <a-menu-item key="1"><NuxtLink to="/">Главная</NuxtLink></a-menu-item>
          <a-menu-item key="3"><NuxtLink to="/gallery">Галлерея</NuxtLink></a-menu-item>
          <a-menu-item key="2"><NuxtLink to="/album-create">Мой Альбом</NuxtLink></a-menu-item>
      </a-menu>
    </a-layout-header>
    <a-layout-content style="padding: 0 50px">
      <a-breadcrumb style="margin: 16px 0">
        <a-breadcrumb-item>Home</a-breadcrumb-item>
      </a-breadcrumb>
      <div :style="{ background: '#fff', padding: '24px', minHeight: '280px' }">
        <Nuxt />
      </div>
    </a-layout-content>
    <a-layout-footer style="text-align: center">
      wizardry &copy; 2022
    </a-layout-footer>
  </a-layout>
</template>

<script>
export default {
  name: 'DefaultLayout'
}
</script>

<style scoped>

</style>
