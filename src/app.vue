<template lang="pug">
  v-app#app
    v-layout(v-resize='onResize')
      side-bar(:windowSize='windowSize')
      top-bar
      v-content
        router-view
    modals
    popups
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import Modals from '@/components/modals/index.vue';
import Popups from '@/components/popups/index.vue';

const { mapState } = createNamespacedHelpers('app');

export default {
  name: 'App',
  components: {
    TopBar,
    SideBar,
    Modals,
    Popups,
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
