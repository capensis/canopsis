<template>
  <alarm-webhook-execution
    v-model="alarm"
    :execution-status="executionStatus"
    :alarms-patterns-params="alarmsPatternsParams"
    :pending="pending"
    :has-errors="hasErrors"
    @run:execution="runTestExecution"
    @clear:execution="clearWebhookStatus"
  >
    <template #webhooks="{ webhooks }">
      <alarm-webhook-execution-timeline :webhooks="webhooks">
        <template #card="{ step }">
          <declare-ticket-rule-execution-webhooks-timeline-card :step="step" />
        </template>
      </alarm-webhook-execution-timeline>
    </template>
  </alarm-webhook-execution>
</template>

<script>
import { SOCKET_ROOMS } from '@/config';

import Socket from '@/plugins/socket/services/socket';

import { formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';
import { isWebhookExecutionFinished } from '@/helpers/entities/webhook-execution/entity';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';

import { validationErrorsMixinCreator } from '@/mixins/form';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import AlarmWebhookExecution from '@/components/other/alarm/partials/alarm-webhook-execution.vue';
import AlarmWebhookExecutionTimeline from '@/components/other/alarm/partials/alarm-webhook-execution-timeline.vue';

import DeclareTicketRuleExecutionWebhooksTimelineCard from './declare-ticket-rule-execution-webhooks-timeline-card.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleExecutionWebhooksTimelineCard,
    AlarmWebhookExecutionTimeline,
    AlarmWebhookExecution,
  },
  mixins: [entitiesDeclareTicketRuleMixin, validationErrorsMixinCreator()],
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
      return `${SOCKET_ROOMS.declareticket}/${id}`;
    },

    async setExecutionStatus(executionStatus) {
      this.executionStatus = executionStatus;
    },

    /**
     * Socket customClose event handler (we need to use for connection checking)
     */
    socketCloseHandler() {
      if (!this.$socket.isConnectionOpen) {
        this.$modals.hide();
        this.$popups.error({
          text: this.$t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    },

    /**
     * Socket closeRoom event handler
     */
    socketCloseRoomHandler() {
      this.$modals.hide();
    },

    async handleSocketError() {
      this.executionStatus = await this.fetchDeclareTicketExecutionWithoutStore({ id: this.executionStatus._id });
    },

    /**
     * Join from execution room
     */
    joinToSocketRoom() {
      this.$socket
        .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .on(Socket.EVENTS_TYPES.error, this.handleSocketError)
        .join(this.getSocketRoomName(this.executionStatus._id))
        .addListener(this.setExecutionStatus);
    },

    /**
     * Leave from execution room
     */
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

        const declareTicket = formToDeclareTicketRule(this.form);

        try {
          this.executionStatus = await this.createTestDeclareTicketExecution({
            data: {
              alarms: [this.alarm],
              ...declareTicket,
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

    clearWebhookStatus() {
      this.executionStatus = null;
    },
  },
};
</script>
