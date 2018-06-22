<template lang="pug">
  v-app#app
    v-layout
      template(v-if="$route.name !== 'login'")
        side-bar
        top-bar
      v-content
        router-view
    modals
    popups
</template>


<script>
// LIBS
import { createNamespacedHelpers } from 'vuex';
// COMPONENTS
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
  computed: {
    ...mapState({
      isSideBarOpen: state => state.app.isSideBarOpen,
    }),
  },
};
</script>
