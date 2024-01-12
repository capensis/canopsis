<template>
  <v-layout
    class="declare-ticket-test-query"
    column
  >
    <v-layout
      align-center
      justify-space-between
    >
      <v-flex xs10>
        <c-alarm-field
          v-model="alarm"
          :disabled="pending || isExecutionRunning"
          :params="alarmsParams"
          name="alarms"
        />
      </v-flex>
      <v-btn
        class="white--text"
        :disabled="hasErrors || !alarm"
        :loading="pending || isExecutionRunning"
        color="orange"
        @click="runTestExecution"
      >
        {{ $t('declareTicket.runTest') }}
      </v-btn>
    </v-layout>
    <v-expand-transition>
      <v-layout
        v-if="executionStatus"
        column
      >
        <v-layout
          class="mb-4"
          align-center
        >
          <span class="text-subtitle-1 mr-5">{{ $t('declareTicket.webhookStatus') }}:</span>
          <declare-ticket-rule-execution-status
            :running="isExecutionRunning"
            :success="isExecutionSucceeded"
            :fail-reason="executionStatus.fail_reason"
          />
          <c-action-btn
            v-if="isExecutionSucceeded || isExecutionFailed"
            type="delete"
            @click="clearWebhookStatus"
          />
        </v-layout>
        <v-card v-if="isSomeOneWebhookStarted">
          <v-card-text>
            <declare-ticket-rule-execution-webhooks-timeline :webhooks="executionStatus.webhooks" />
          </v-card-text>
        </v-card>
      </v-layout>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { SOCKET_ROOMS } from '@/config';

import Socket from '@/plugins/socket/services/socket';

import {
  formToDeclareTicketRule,
  isDeclareTicketExecutionFailed,
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
  isDeclareTicketExecutionWaiting,
} from '@/helpers/entities/declare-ticket/rule/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';

import { validationErrorsMixinCreator } from '@/mixins/form';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';
import DeclareTicketRuleExecutionWebhooksTimeline from './declare-ticket-rule-execution-webhooks-timeline.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleExecutionWebhooksTimeline, DeclareTicketRuleExecutionStatus },
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

    isExecutionRunning() {
      return this.executionStatus && isDeclareTicketExecutionRunning(this.executionStatus);
    },

    isExecutionSucceeded() {
      return isDeclareTicketExecutionSucceeded(this.executionStatus);
    },

    isExecutionFailed() {
      return isDeclareTicketExecutionFailed(this.executionStatus);
    },

    isSomeOneWebhookStarted() {
      return this.executionStatus?.webhooks.some(webhook => !isDeclareTicketExecutionWaiting(webhook));
    },

    alarmsPatternsParams() {
      return Object.entries(formFilterToPatterns(this.form.patterns))
        .reduce((acc, [key, value]) => {
          acc[key] = JSON.stringify(value);

          return acc;
        }, {});
    },

    alarmsParams() {
      return {
        opened: true,
        ...this.alarmsPatternsParams,
      };
    },
  },
  watch: {
    executionStatus(executionStatus) {
      if (
        executionStatus
        && (isDeclareTicketExecutionSucceeded(executionStatus) || isDeclareTicketExecutionFailed(executionStatus))
      ) {
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
