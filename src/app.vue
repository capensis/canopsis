<template lang="pug">
  v-app#app
    #nav
      router-link(to="/") Home
      router-link(to="/about") About
    router-view
    popups
</template>

<script>
import Popups from '@/components/popups/index.vue';

export default {
  components: { Popups },
};
</script>
