<template lang="pug">
  v-app#app
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
