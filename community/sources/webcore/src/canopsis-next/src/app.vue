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

import { SOCKET_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';
import { MAX_LIMIT } from '@/constants';

import TheNavigation from '@/components/layout/navigation/the-navigation.vue';
import TheSideBars from '@/components/side-bars/the-sidebars.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import { authMixin } from '@/mixins/auth';
import systemMixin from '@/mixins/system';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesViewStatsMixin } from '@/mixins/entities/view-stats';
import entitiesUserMixin from '@/mixins/entities/user';

import '@/assets/styles/main.scss';

import localStorageService from '@/services/local-storage';

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
    entitiesViewStatsMixin,
    entitiesUserMixin,
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

    this.pending = false;
  },
  beforeDestroy() {
    this.stopViewStats();
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedExecutionsWithoutStore',
    }),

    registerCurrentUserOnceWatcher() {
      const unwatch = this.$watch('currentUser', async (currentUser) => {
        if (!isEmpty(currentUser)) {
          this.$socket.connect(`${SOCKET_URL}?token=${localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY)}`);
          await this.fetchAppInfos();

          this.setSystemData({
            timezone: this.timezone,
            jobExecutorFetchTimeoutSeconds: this.jobExecutorFetchTimeoutSeconds,
          });

          this.setTitle();
          this.startViewStats();
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
