<template lang="pug">
  v-app#app
    login
</template>


<script>
import Login from '@/views/login.vue';
import { createNamespacedHelpers } from 'vuex';

const { mapState } = createNamespacedHelpers('app');

export default {
  name: 'App',
  components: {
    Login,
  },
  data() {
    return {
      windowSize: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    ...mapState({
      isSideBarOpen: state => state.app.isSideBarOpen,
    }),
  },
  methods: {
    onResize() {
      this.windowSize = { x: window.innerWidth, y: window.innerHeight };
    },
  },
};
</script>
