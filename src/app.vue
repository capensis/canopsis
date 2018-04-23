<template lang="pug">
  v-app#app
    #nav
      router-link(to="/") Home
      router-link(to="/about") About
    router-view
    modals
</template>

<script>
import Modals from './components/modals/index.vue';

export default {
  components: {
    Modals,
  },
};
</script>
