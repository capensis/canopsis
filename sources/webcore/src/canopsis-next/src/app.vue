<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      navigation#main-navigation(v-if="$route.name !== 'login'")
      v-content#main-content
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

    this.pending = false;
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
        padding: 0!important;
      }
    }
  }
</style>
