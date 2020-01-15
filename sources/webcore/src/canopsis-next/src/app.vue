<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      navigation#main-navigation(v-if="$route.name !== 'login'")
      v-content#main-content
        router-view(:key="routeViewKey")
    side-bars
    the-modals
    the-popups
</template>


<script>
import Navigation from '@/components/layout/navigation/index.vue';
import SideBars from '@/components/side-bars/index.vue';

import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';
import keepaliveMixin from '@/mixins/entities/keepalive';


import '@/assets/styles/main.scss';

export default {
  components: {
    Navigation,
    SideBars,
  },
  mixins: [authMixin, entitiesInfoMixin, keepaliveMixin],
  data() {
    return {
      pending: true,

    };
  },
  computed: {
    routeViewKey() {
      if (this.$route.name === 'view') {
        return this.$route.path;
      }

      return this.$route.fullPath;
    },
  },
  async mounted() {
    await this.fetchCurrentUser();

    if (this.isLoggedIn) {
      await this.fetchAppInfos();
    }
    this.pending = false;
  },
  beforeUpdate() {
    if (this.isLoggedIn) {
      this.startKeepAlive();
    }
  },
  updated() {
    if (this.isLoggedIn) {
      this.startSessionHide();
    }
  },
  beforeDestroy() {
    this.stopKeepAlive();
  },


};
</script>

<style lang="scss">
  #app {
    &.-fullscreen {
      width: 100%;

      #main-navigation {
        display: none;
      }

      #main-content {
        padding: 0 !important;
      }
    }
  }
</style>
