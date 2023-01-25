<template lang="pug">
  modal-wrapper(close, minimize)
    template(#title="")
      v-layout(align-center)
        span {{ title }}
        declare-ticket-rule-execution-status.ml-2(
          v-if="modal.minimized",
          :running="isExecutionsRunning",
          :success="isExecutionsSucceeded",
          :fail-reason="failReason",
          color="white"
        )
    template(#text="")
      v-layout.mb-4(v-if="isOneExecution", row, align-center)
        span.subheading.mr-5 {{ $t('declareTicket.webhookStatus') }}:
        declare-ticket-rule-execution-status(
          :running="isExecutionsRunning",
          :success="isExecutionsSucceeded",
          :fail-reason="failReason"
        )
      declare-ticket-rule-execution(:alarm-executions="alarmExecutions", :is-one-execution="isOneExecution")
    template(#actions="")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.close') }}
</template>

<script>
import { SOCKET_ROOMS } from '@/config';

import { DECLARE_TICKET_EXECUTION_STATUSES, MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import DeclareTicketRuleExecution from '@/components/other/declare-ticket/partials/declare-ticket-rule-execution.vue';
import DeclareTicketRuleExecutionStatus from '@/components/other/declare-ticket/partials/declare-ticket-rule-execution-status.vue';

import ModalWrapper from '../modal-wrapper.vue';
import {
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/forms/declare-ticket-rule';
import Socket from '@/plugins/socket/services/socket';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.executeDeclareTicket,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DeclareTicketRuleExecutionStatus, DeclareTicketRuleExecution, ModalWrapper },
  mixins: [
    modalInnerMixin,
  ],
  data() {
    return {
      executionsStatusesById: {},
    };
  },
  computed: {
    isOneExecution() {
      return this.config.executions.length === 1;
    },

    alarmExecutions() {
      return this.config.executions.reduce((acc, { executionId, ruleName, alarms }) => {
        alarms.forEach((alarm) => {
          const status = this.executionsStatusesById[executionId] ?? {
            status: DECLARE_TICKET_EXECUTION_STATUSES.running,
          };

          acc.push({
            alarm,
            executionId,
            ruleName,
            ...status,
          });
        });

        return acc;
      }, []);
    },

    isExecutionsRunning() {
      return this.alarmExecutions.some(isDeclareTicketExecutionRunning);
    },

    isExecutionsSucceeded() {
      return this.alarmExecutions.every(isDeclareTicketExecutionSucceeded);
    },

    failReason() {
      return Object.values(this.executionsStatusesById).map(execution => execution.fail_reason).join('\n');
    },

    title() {
      if (this.config.title) {
        return this.config.title;
      }

      return this.isOneExecution
        ? this.config.executions[0].ruleName
        : this.$t('modals.executeDeclareTicket.title');
    },
  },
  mounted() {
    this.joinToSocketRooms();
  },
  beforeDestroy() {
    this.leaveFromSocketRooms();
  },
  methods: {
    getSocketRoomName(id) {
      return `${SOCKET_ROOMS.declareticket}/${id}`;
    },

    setExecutionStatus(executionStatus) {
      this.$set(this.executionsStatusesById, executionStatus._id, executionStatus);
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
    joinToSocketRooms() {
      this.config.executions.forEach(({ executionId }) => {
        this.$socket
          .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
          .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
          .join(this.getSocketRoomName(executionId))
          .addListener(this.setExecutionStatus);
      });
    },

    /**
     * Leave from execution room
     */
    leaveFromSocketRooms() {
      this.config.executions.forEach(({ executionId }) => {
        this.$socket
          .off(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
          .off(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
          .leave(this.getSocketRoomName(executionId))
          .removeListener(this.setExecutionStatus);
      });
    },
  },
};
</script>
