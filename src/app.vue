<template lang="pug">
  v-app#app
    v-layout
      side-bar
      top-bar
      v-content
        router-view
    modals
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import Modals from '@/components/modals/index.vue';

const { mapState } = createNamespacedHelpers('app');

export default {
  name: 'App',
  components: {
    TopBar,
    SideBar,
    Modals,
  },
  computed: {
    ...mapState({
      isSideBarOpen: state => state.app.isSideBarOpen,
    }),
  },
};
</script>
