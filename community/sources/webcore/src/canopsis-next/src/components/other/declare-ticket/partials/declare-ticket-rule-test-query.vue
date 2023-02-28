<template lang="pug">
  v-layout.declare-ticket-test-query(column)
    v-layout(row, align-center, justify-space-between)
      v-flex(xs10)
        c-alarm-field(
          v-model="alarm",
          :disabled="pending || isExecutionRunning",
          :params="alarmsParams",
          name="alarms"
        )
      v-btn.white--text(
        :disabled="hasErrors || !alarm",
        :loading="pending || isExecutionRunning",
        color="orange",
        @click="runTestExecution"
      ) {{ $t('declareTicket.runTest') }}
    v-expand-transition
      v-layout(v-if="executionStatus", column)
        v-layout.mb-4(row, align-center)
          span.subheading.mr-5 {{ $t('declareTicket.webhookStatus') }}:
          declare-ticket-rule-execution-status(
            :running="isExecutionRunning",
            :success="isExecutionSucceeded",
            :fail-reason="executionStatus.fail_reason"
          )
          c-action-btn(v-if="webhooksResponses.length", type="delete", @click="clearResponses")
        v-card.grey.lighten-4.mb-2(v-for="(webhooksResponse, index) in webhooksResponses", :key="index", light, flat)
          v-card-text
            c-request-text-information(:value="webhooksResponse")
</template>

<script>
import { SOCKET_ROOMS } from '@/config';

import Socket from '@/plugins/socket/services/socket';

import {
  formToDeclareTicketRule,
  isDeclareTicketExecutionFailed,
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/forms/declare-ticket-rule';
import { formFilterToPatterns } from '@/helpers/forms/filter';

import { validationErrorsMixinCreator } from '@/mixins/form';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleExecutionStatus },
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
      webhooksResponses: [],
    };
  },
  computed: {
    hasErrors() {
      return this.errors.any();
    },

    isExecutionRunning() {
      return this.executionStatus && isDeclareTicketExecutionRunning(this.executionStatus);
    },

    isExecutionSucceeded() {
      return isDeclareTicketExecutionSucceeded(this.executionStatus);
    },

    alarmsParams() {
      return Object.entries(formFilterToPatterns(this.form.patterns)).reduce((acc, [key, value]) => {
        acc[key] = JSON.stringify(value);

        return acc;
      }, {});
    },
  },
  watch: {
    executionStatus(executionStatus) {
      if (
        executionStatus
        && (isDeclareTicketExecutionSucceeded(executionStatus) || isDeclareTicketExecutionFailed(executionStatus))
      ) {
        this.leaveFromSocketRoom();
        this.fetchTestWebhooksResponse();
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

    setExecutionStatus(executionStatus) {
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

    /**
     * Join from execution room
     */
    joinToSocketRoom() {
      this.$socket
        .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
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
        .leave(this.getSocketRoomName(this.executionStatus._id))
        .removeListener(this.setExecutionStatus);
    },

    async runTestExecution() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        this.pending = true;
        this.clearResponses();

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

    async fetchTestWebhooksResponse() {
      const webhooksWithResponses = this.executionStatus.webhooks.filter(
        execution => isDeclareTicketExecutionSucceeded(execution) || isDeclareTicketExecutionFailed(execution),
      );

      this.webhooksResponses = await Promise.all(
        webhooksWithResponses.map(({ _id: id }) => this.fetchTestDeclareTicketExecutionWebhooksResponse({ id })),
      );
    },

    clearResponses() {
      this.webhooksResponses = [];
    },
  },
};
</script>
