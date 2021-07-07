<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      the-navigation#main-navigation(v-if="$route.name !== 'login'")
      v-content#main-content
        active-broadcast-message
        router-view(:key="routeViewKey")
    the-side-bars
    the-modals
    the-popups
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isEmpty } from 'lodash';

import TheNavigation from '@/components/layout/navigation/the-navigation.vue';
import TheSideBars from '@/components/side-bars/the-sidebars.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import { authMixin } from '@/mixins/auth';
import systemMixin from '@/mixins/system';
import entitiesInfoMixin from '@/mixins/entities/info';
import entitiesUserMixin from '@/mixins/entities/user';
import keepaliveMixin from '@/mixins/entities/keepalive';

import '@/assets/styles/main.scss';
import { MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('remediationInstructionExecution');

export default {
  components: {
    TheNavigation,
    TheSideBars,
    ActiveBroadcastMessage,
  },
  mixins: [
    authMixin,
    systemMixin,
    entitiesInfoMixin,
    entitiesUserMixin,
    keepaliveMixin,
  ],
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
  created() {
    this.registerCurrentUserOnceWatcher();
  },
  async mounted() {
    await this.fetchCurrentUser();

    const { errorMessage } = this.$route.query;

    if (errorMessage) {
      this.$popups.error({ text: errorMessage, autoClose: 7000 });
    }

    this.pending = false;
  },
  beforeDestroy() {
    this.stopKeepalive();
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedExecutionsWithoutStore',
    }),

    registerCurrentUserOnceWatcher() {
      const unwatch = this.$watch('currentUser', async (currentUser) => {
        if (!isEmpty(currentUser)) {
          await this.fetchAppInfos();

          this.setSystemData({
            timezone: this.timezone,
            jobExecutorFetchTimeoutSeconds: this.jobExecutorFetchTimeoutSeconds,
          });

          this.setTitle();

          this.startKeepalive();
          this.showPausedExecutionsPopup();

          unwatch();
        }
      });
    },

    async showPausedExecutionsPopup() {
      const pausedExecutions = await this.fetchPausedExecutionsWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      if (!pausedExecutions || !pausedExecutions.length) {
        return;
      }

      pausedExecutions.forEach((execution = {}) => this.$popups.info({
        text: this.$t('remediationInstructionExecute.popups.wasPaused', {
          instructionName: execution.instruction_name,
          alarmName: execution.alarm_name,
          date: this.$options.filters.date(execution.paused, 'long', true),
        }),
      }));
    },
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
