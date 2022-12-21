<script lang="ts" setup>
import { sidebarRoutes } from '@/router'
import { useMe } from '@/store/me'

const me = useMe()
</script>

<template>
  <div class="sidebar-container">
    <div class="routes-container">
      <div v-for="route in sidebarRoutes" :key="route.name" class="route">
        <router-link
          v-if="route.meta && $route.meta.title"
          :to="{ name: route.name }"
          class="link"
          :class="{ 'active-link': $route.name === route.name }"
        >
          {{ route.name }}
        </router-link>
      </div>
    </div>
    <div class="logout">
      <router-link :to="{ name: 'Login' }" class="link">
        <button class="button" :on-click="me.logout">Logout</button>
      </router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: $bgcolor-primary;
  width: 20%;
  min-width: fit-content;
  max-width: 10rem;
  padding: 5rem 1.75rem;
}
.routes-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}
.route {
  margin-bottom: 1.5rem;
  width: 60%;
  min-width: 5rem;
}
.link {
  text-decoration: none;
  font-size: larger;
  font-weight: 500;
  color: black;
  transition: opacity 0.2s ease-in-out;
  &:hover {
    opacity: 0.6;
  }
}
.active-link {
  color: $color-primary;
}

.logout {
  margin-top: auto;
}
.button {
  background-color: $color-primary;
  border: none;
  border-radius: 0.5rem;
  font-size: medium;
  padding: 0.5rem 1rem;
  color: white;
  transition: opacity 0.2s ease-in-out;
  cursor: pointer;
  &:hover {
    opacity: 0.8;
  }
}
</style>
