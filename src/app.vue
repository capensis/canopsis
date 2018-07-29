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
import { createNamespacedHelpers } from 'vuex';

import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import Modals from '@/components/modals/index.vue';
import Popups from '@/components/popups/index.vue';

const { mapActions } = createNamespacedHelpers('auth');

export default {
  components: {
    TopBar,
    SideBar,
    Modals,
    Popups,
  },
  mounted() {
    this.fetchCurrentUser();
  },
  methods: {
    ...mapActions(['fetchCurrentUser']),
  },
};
</script>
