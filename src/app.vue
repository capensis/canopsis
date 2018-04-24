<template lang="pug">
  v-app#app
    //#nav
      router-link(to="/") Home
      router-link(to="/about") About
    //router-view
    login
</template>

<script>
import Login from '@/views/login.vue';

export default {
  name: 'app',
  components: {
    Login,
  },
};
</script>
