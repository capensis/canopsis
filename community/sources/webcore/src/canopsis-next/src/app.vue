<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      navigation#main-navigation(v-if="$route.name !== 'login'")
      v-content#main-content
        active-broadcast-message
        router-view(:key="routeViewKey")
    side-bars
    the-modals
    the-popups
</template>

<script>
import { isEmpty } from 'lodash';

import { prepareUserByData } from '@/helpers/entities';

import Navigation from '@/components/layout/navigation/index.vue';
import SideBars from '@/components/side-bars/index.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import authMixin from '@/mixins/auth';
import systemMixin from '@/mixins/system';
import entitiesInfoMixin from '@/mixins/entities/info';
import entitiesUserMixin from '@/mixins/entities/user';
import keepaliveMixin from '@/mixins/entities/keepalive';

import '@/assets/styles/main.scss';

export default {
  components: {
    Navigation,
    SideBars,
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

    this.pending = false;
  },
  beforeDestroy() {
    this.stopKeepalive();
  },
  methods: {
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
      }, { immediate: true });
    },

    async showPausedExecutionsPopup() {
      const { paused_executions: pausedExecutions = [] } = this.currentUser;

      if (!pausedExecutions.length) {
        return;
      }

      pausedExecutions.forEach((execution = {}) => this.$popups.info({
        text: this.$t('remediationInstructionExecute.popups.wasPaused', {
          instructionName: execution.instruction_name,
          alarmName: execution.alarm_name,
          date: this.$options.filters.date(execution.paused, 'long', true),
        }),
      }));

      const data = prepareUserByData({}, this.currentUser);

      data.paused_executions = [];

      await this.createUser({ data });
      await this.fetchCurrentUser();
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
