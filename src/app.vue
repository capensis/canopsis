<template lang="pug">
  v-app#app
    v-layout
      template(v-if="$route.name !== 'login'")
        side-bar
        top-bar
      v-content(v-if="!currentUserPending")
        router-view
    modals
    popups
</template>


<script>
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import Modals from '@/components/modals/index.vue';
import Popups from '@/components/popups/index.vue';
import authMixin from '@/mixins/auth';

export default {
  components: {
    TopBar,
    SideBar,
    Modals,
    Popups,
  },
  mixins: [authMixin],
  mounted() {
    this.fetchCurrentUser();
  },
};
</script>
