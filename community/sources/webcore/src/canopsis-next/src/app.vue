<template lang="pug">
  v-app#app(:dark="system.dark")
    c-progress-overlay(
      :pending="wholePending",
      :transition="false",
      color="gray"
    )
    v-layout(v-if="!wholePending")
      the-navigation#main-navigation(v-if="shownNavigation")
      v-content#main-content
        active-broadcast-message
        router-view(:key="routeViewKey")
    the-sidebar
    the-modals
    the-popups
</template>

<script>
import { isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';
import { EXCLUDED_SERVER_ERROR_STATUSES, MAX_LIMIT, ROUTES_NAMES } from '@/constants';

import { reloadPageWithTrailingSlashes } from '@/helpers/url';
import { convertDateToString } from '@/helpers/date/date';

import localStorageService from '@/services/local-storage';

import Socket from '@/plugins/socket/services/socket';

import { authMixin } from '@/mixins/auth';
import { systemMixin } from '@/mixins/system';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesUserMixin } from '@/mixins/entities/user';
import { entitiesTemplateVarsMixin } from '@/mixins/entities/template-vars';

import TheNavigation from '@/components/layout/navigation/the-navigation.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import '@/assets/styles/main.scss';

const { mapActions } = createNamespacedHelpers('remediationInstructionExecution');

export default {
  components: {
    TheNavigation,
    ActiveBroadcastMessage,
  },
  mixins: [
    authMixin,
    systemMixin,
    entitiesInfoMixin,
    entitiesUserMixin,
    entitiesTemplateVarsMixin,
  ],
  data() {
    return {
      currentUserLocalPending: true,
      appInfoLocalPending: true,
    };
  },
  computed: {
    wholePending() {
      return this.currentUserLocalPending || this.appInfoLocalPending || this.templateVarsPending;
    },

    routeViewKey() {
      if (this.$route.name === ROUTES_NAMES.view) {
        return this.$route.path;
      }

      return this.$route.fullPath;
    },

    shownNavigation() {
      return this.currentUser && !this.$route.meta.hideNavigation;
    },
  },
  beforeCreate() {
    reloadPageWithTrailingSlashes();
  },
  created() {
    this.registerCurrentUserOnceWatcher();
  },
  mounted() {
    this.fetchAppInfoWithErrorHandling();
    this.socketConnectWithErrorHandling();
    this.fetchCurrentUserWithErrorHandling();
    this.showLocalStorageWarningPopupMessage();
    this.fetchTemplateVars();
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
    }),

    showLocalStorageWarningPopupMessage() {
      const text = localStorageService.pop('warningPopup');

      if (text) {
        this.$popups.warning({ text, autoClose: false });
      }
    },

    registerCurrentUserOnceWatcher() {
      const unwatch = this.$watch('currentUser', async (currentUser) => {
        if (!isEmpty(currentUser)) {
          this.$socket.authenticate(localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY));
          this.setTheme(currentUser.ui_theme);

          await this.filesAccess();

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
        text: this.$t('remediation.instructionExecute.popups.wasPaused', {
          instructionName: execution.instruction_name,
          alarmName: execution.alarm_name,
          date: convertDateToString(execution.paused),
        }),
      }));
    },

    socketConnectWithErrorHandling() {
      try {
        this.$socket
          .connect(SOCKET_URL)
          .on('error', this.socketErrorHandler);
      } catch (err) {
        this.$popups.error({
          text: this.$t('errors.socketConnectionProblem'),
          autoClose: false,
        });

        console.error(err);
      }
    },

    socketErrorHandler({ message } = {}) {
      if (message) {
        this.$popups.error({ text: message });

        if (message === Socket.ERROR_MESSAGES.authenticationFailed) {
          localStorageService.set('warningPopup', this.$t('warnings.authTokenExpired'));
          this.logout();
        }
      }
    },

    async fetchCurrentUserWithErrorHandling() {
      try {
        this.currentUserLocalPending = true;

        await this.fetchCurrentUser();
      } catch (err) {
        if (!EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
          this.$router.push({ name: ROUTES_NAMES.error });
        }

        console.error(err);
      } finally {
        this.currentUserLocalPending = false;
      }
    },

    async fetchAppInfoWithErrorHandling() {
      try {
        this.appInfoLocalPending = true;
        await this.fetchAppInfo();

        this.setSystemData({
          timezone: this.timezone,
        });

        this.setTitle();
      } catch (err) {
        if (!EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
          this.$router.push({
            name: ROUTES_NAMES.error,
          });
        }

        console.error(err);
      } finally {
        this.appInfoLocalPending = false;
      }
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
