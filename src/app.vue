<template lang="pug">
  v-app#app
    v-layout(v-resize='onResize')
      template(v-if="$route.name !== 'login'")
        side-bar(:windowSize='windowSize')
        top-bar
      v-content
        router-view
    modals
</template>


<script>
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import { createNamespacedHelpers } from 'vuex';
import Modals from '@/components/modals/index.vue';

const { mapState } = createNamespacedHelpers('app');

export default {
  name: 'App',
  components: {
    TopBar,
    SideBar,
    Modals,
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
