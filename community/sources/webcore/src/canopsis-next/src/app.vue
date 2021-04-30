<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      the-navigation#main-navigation(v-if="$route.name !== 'login'")
      v-content#main-content
        v-btn(@click="approve") APPROVE
        active-broadcast-message
        router-view(:key="routeViewKey")
    the-side-bars
    the-modals
    the-popups
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isEmpty } from 'lodash';

import { MAX_LIMIT, MODALS } from '@/constants';

import TheNavigation from '@/components/layout/navigation/the-navigation.vue';
import TheSideBars from '@/components/side-bars/the-sidebars.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import { authMixin } from '@/mixins/auth';
import systemMixin from '@/mixins/system';
import entitiesInfoMixin from '@/mixins/entities/info';
import entitiesUserMixin from '@/mixins/entities/user';
import keepaliveMixin from '@/mixins/entities/keepalive';

import '@/assets/styles/main.scss';

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

    approve() {
      this.$modals.show({
        name: MODALS.remediationInstructionApproval,
        config: {
          remediationInstruction: {
            _id: '98665a84-cede-4778-afe8-07b299b00892',
            type: 0,
            status: 0,
            alarm_patterns: [
              {
                _id: '1fe974af-489a-4fdd-a37f-d5b3fb7d160a',
              },
            ],
            entity_patterns: null,
            name: 'Instruction',
            description: 'Description',
            author: 'root',
            enabled: true,
            timeout_after_execution: {
              seconds: 0,
              unit: '',
            },
            steps: [
              {
                name: 'Servers',
                operations: [
                  {
                    name: 'Images',
                    time_to_complete: {
                      seconds: 600,
                      unit: 'm',
                    },
                    description: '\u003cp\u003eCould you a check all servers?\u003c/p\u003e',
                    jobs: [
                      {
                        _id: '78bbdc94-3c00-475d-9974-bc96a336437d',
                        name: 'Test job',
                        author: 'root',
                        config: {
                          _id: '5a55233b-ac6e-43b0-938f-43a2be07e328',
                          name: 'Config name 4',
                          type: 'rundeck',
                          host: 'http://host2.host',
                          author: 'root',
                          auth_token: 'AuthToken',
                        },
                        job_id: 'string',
                        payload: '{}',
                      },
                      {
                        _id: '5fa8eaa4-2223-4a59-b450-40908a39ff4d',
                        name: 'Run diagnostic',
                        author: 'root',
                        config: {
                          _id: 'ef7e65ba-0e43-4b41-b44b-00f19e6c51fd',
                          name: 'Run diagnostic',
                          type: 'rundeck',
                          host: 'http://host.ru',
                          author: 'root',
                          auth_token: 'token',
                        },
                        job_id: '123123123',
                        payload: '{}',
                      },
                    ],
                  },
                  {
                    name: 'Network',
                    time_to_complete: {
                      seconds: 300,
                      unit: 'm',
                    },
                    description: '\u003cp\u003eCould you a check all network?\u003c/p\u003e\u003cp\u003e\u003cbr\u003e\u003c/p\u003e\u003cp\u003e\u003ca href="/api/api/v4/cat/file/0e4c9bc3-da69-4269-a93c-24c74aaa72a1" title="number.jpg.jpeg" target="_blank"\u003e\u003cimg src="/backend/api/v4/cat/file/cbaf7b3e-5d57-45f0-ba86-407db34d14b6" style="width: 300px;"\u003e\u003cbr\u003e\u003c/a\u003e\u003c/p\u003e',
                    jobs: [
                      {
                        _id: '651cc21b-02fc-41be-8da5-712c9aa92921',
                        name: 'Job name22',
                        author: 'root',
                        config: {
                          _id: '67e16d23-ea17-4f29-b6ec-629d9f65102c',
                          name: 'Name',
                          type: 'awx',
                          host: 'http://localhost.ru',
                          author: 'root',
                          auth_token: '124125125125',
                        },
                        job_id: '12512512512512512',
                        payload: '{"alarm":"{{ .Alarm }}"}',
                      },
                    ],
                  },
                ],
                stop_on_fail: false,
                endpoint: 'Are the servers running?',
              },
              {
                name: 'System',
                operations: [
                  {
                    name: 'Kubernetes',
                    time_to_complete: {
                      seconds: 120,
                      unit: 'm',
                    },
                    description: 'Kubernetes operation desc',
                  },
                ],
                stop_on_fail: true,
                endpoint: 'Are the systems running?',
              },
            ],
            created: 1603787616,
            last_modified: 1611133139,
            running: true,
            deletable: false,
            last_executed_on: 1608266407,
          },
        },
      });
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
