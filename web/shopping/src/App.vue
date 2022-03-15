<template>
  <div id="app">
  <img alt="Vue logo" src="./assets/cuhk_logo.png">
  <ShoppingProducts></ShoppingProducts>
  </div>
</template>

<script>
import ShoppingProducts from '@/components/ShoppingProducts.vue'
export default {
  name: 'App',
  components: {
    ShoppingProducts
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
