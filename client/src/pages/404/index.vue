<template>
  <div class="not-found-container">
    <div class="title">404 Not Found</div>
    <router-link :to="{ name: 'Home' }">back to top</router-link>
  </div>
</template>

<style lang="scss" scoped>
.not-found-container {
  width: 600px;
  margin: 0 auto;
  padding-top: 20vh;
  text-align: center;

  .title {
    font-size: 30px;
    margin-bottom: 20px;
  }
}
</style>
