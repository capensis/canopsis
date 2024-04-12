<template>
  <alarm-test-query-execution
    v-model="alarm"
    :execution-status="executionStatus"
    :alarms-patterns-params="alarmsPatternsParams"
    :pending="pending"
    :has-errors="hasErrors"
    @run:execution="runTestExecution"
    @clear:status="clearWebhookStatus"
  >
    <template #webhooks="{ webhooks }">
      <alarm-test-query-execution-webhooks-timeline :webhooks="webhooks">
        <template #card="{ step }">
          <scenario-execution-webhooks-timeline-card :step="step" />
        </template>
      </alarm-test-query-execution-webhooks-timeline>
    </template>
  </alarm-test-query-execution>
</template>

<script>
import { SOCKET_ROOMS } from '@/config';

import Socket from '@/plugins/socket/services/socket';

import { isWebhookExecutionFinished } from '@/helpers/entities/webhook-execution/entity';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { formToScenario } from '@/helpers/entities/scenario/form';

import { validationErrorsMixinCreator } from '@/mixins/form';
import { entitiesScenarioMixin } from '@/mixins/entities/scenario';

import AlarmTestQueryExecution from '@/components/other/alarm/partials/alarm-test-query-execution.vue';
import AlarmTestQueryExecutionWebhooksTimeline from '@/components/other/alarm/partials/alarm-test-query-execution-webhooks-timeline.vue';

import ScenarioExecutionWebhooksTimelineCard from './scenario-execution-webhooks-timeline-card.vue';

export default {
  inject: ['$validator'],
  components: {
    ScenarioExecutionWebhooksTimelineCard,
    AlarmTestQueryExecutionWebhooksTimeline,
    AlarmTestQueryExecution,
  },
  mixins: [entitiesScenarioMixin, validationErrorsMixinCreator()],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      query: {
        search: null,
      },

      alarm: '',
      pending: false,
      executionStatus: undefined,
    };
  },
  computed: {
    hasErrors() {
      return this.errors.any();
    },

    alarmsPatternsParams() {
      return Object.entries(formFilterToPatterns(this.form.patterns))
        .reduce((acc, [key, value]) => {
          acc[key] = JSON.stringify(value);

          return acc;
        }, {});
    },
  },
  watch: {
    executionStatus(executionStatus) {
      if (isWebhookExecutionFinished(executionStatus)) {
        this.leaveFromSocketRoom();
      }
    },
  },
  beforeDestroy() {
    if (this.executionStatus) {
      this.leaveFromSocketRoom();
    }
  },
  methods: {
    getSocketRoomName(id) {
      return `${SOCKET_ROOMS.testscenario}/${id}`;
    },

    setExecutionStatus(executionStatus) {
      this.executionStatus = executionStatus;
    },

    clearWebhookStatus() {
      this.executionStatus = null;
    },

    socketCloseHandler() {
      if (!this.$socket.isConnectionOpen) {
        this.$modals.hide();
        this.$popups.error({
          text: this.$t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    },

    socketCloseRoomHandler() {
      this.$modals.hide();
    },

    async handleSocketError() {
      const executionStatus = await this.fetchTestScenarioExecutionWithoutStore({ id: this.executionStatus._id });

      this.setExecutionStatus(executionStatus);
    },

    joinToSocketRoom() {
      this.$socket
        .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .on(Socket.EVENTS_TYPES.error, this.handleSocketError)
        .join(this.getSocketRoomName(this.executionStatus._id))
        .addListener(this.setExecutionStatus);
    },

    leaveFromSocketRoom() {
      this.$socket
        .off(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .off(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .off(Socket.EVENTS_TYPES.error, this.handleSocketError)
        .leave(this.getSocketRoomName(this.executionStatus._id))
        .removeListener(this.setExecutionStatus);
    },

    async runTestExecution() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        this.pending = true;
        this.clearWebhookStatus();

        try {
          this.executionStatus = await this.createTestScenarioExecution({
            data: {
              alarm: this.alarm,
              ...formToScenario(this.form),
            },
          });

          this.joinToSocketRoom();
        } catch (err) {
          if (err.error) {
            this.$popups.error({ text: err.error });
          } else {
            this.setFormErrors(err);
          }

          this.executionStatus = undefined;
        } finally {
          this.pending = false;
        }
      }
    },
  },
};
</script>
