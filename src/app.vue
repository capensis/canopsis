<template lang="pug">
  v-app#app
    v-layout
      navigation(v-if="$route.name !== 'login'")
      v-content#main-content(v-if="!currentUserPending")
        router-view(:key="$route.fullPath")
    side-bars
    modals
    popups
</template>


<script>
import Navigation from '@/components/layout/navigation/index.vue';
import SideBars from '@/components/side-bars/index.vue';
import Modals from '@/components/modals/index.vue';
import Popups from '@/components/popups/index.vue';

import authMixin from '@/mixins/auth';

import '@/assets/styles/main.scss';

export default {
  components: {
    Navigation,
    SideBars,
    Modals,
    Popups,
  },
  mixins: [authMixin],
  mounted() {
    this.fetchCurrentUser();
  },
};
</script>

<style lang="scss">
  #main-content {
    &.-fullscreen {
      background: #fafafa;
      padding: 0!important;
    }
  }
</style>
