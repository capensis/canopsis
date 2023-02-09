<template lang="pug">
  modal-wrapper(close, minimize)
    template(#title="")
      v-layout(align-center)
        span {{ title }}
        declare-ticket-rule-execution-status.ml-2.declare-ticket-rule-execute-status(
          v-if="modal.minimized",
          :running="isExecutionsRunning",
          :success="isExecutionsSucceeded",
          :fail-reason="failReason",
          color="white"
        )
    template(#text="")
      v-layout(v-if="pending", justify-center)
        v-progress-circular(color="primary", indeterminate)
      template(v-else)
        v-layout.mb-4(v-if="isOneTicket", row, align-center)
          span.subheading.mr-5 {{ $t('declareTicket.webhookStatus') }}:
          declare-ticket-rule-execution-status(
            :running="isExecutionsRunning",
            :success="isExecutionsSucceeded",
            :fail-reason="failReason"
          )
        declare-ticket-rule-execution-alarms(:alarm-executions="alarmExecutions", :is-one-execution="isOneTicket")
    template(#actions="")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.close') }}
</template>

<script>
import { keyBy } from 'lodash';

import { SOCKET_ROOMS } from '@/config';

import { DECLARE_TICKET_EXECUTION_STATUSES, MODALS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import {
  isDeclareTicketExecutionFailed,
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/forms/declare-ticket-rule';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import DeclareTicketRuleExecutionStatus from '@/components/other/declare-ticket/partials/declare-ticket-rule-execution-status.vue';
import DeclareTicketRuleExecutionAlarms from '@/components/other/declare-ticket/partials/declare-ticket-rule-execution-alarms.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to execute declare tickets
 */
export default {
  name: MODALS.executeDeclareTickets,
  components: {
    DeclareTicketRuleExecutionAlarms,
    DeclareTicketRuleExecutionStatus,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    entitiesDeclareTicketRuleMixin,
  ],
  data() {
    return {
      pending: false,
      successExecutions: [],
      executionsStatusesById: {},
    };
  },
  computed: {
    isOneTicket() {
      return this.config.tickets.length === 1;
    },

    alarmExecutions() {
      return this.successExecutions.reduce((acc, { executionId, ruleName, alarms }) => {
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

    isExecutionsFailed() {
      return this.alarmExecutions.every(isDeclareTicketExecutionFailed);
    },

    failReason() {
      return Object.values(this.executionsStatusesById).map(execution => execution.fail_reason).join('\n');
    },

    title() {
      if (this.config.title) {
        return this.config.title;
      }

      return this.isOneTicket
        ? this.config.tickets[0].name
        : this.$t('modals.executeDeclareTickets.title');
    },
  },
  watch: {
    alarmExecutions(value) {
      if (value.length && (this.isExecutionsFailed || this.isExecutionsSucceeded)) {
        this.config.onExecute?.();
      }
    },
  },
  async mounted() {
    await this.createExecutions();
    await this.fetchExecutionsStatuses();
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
      this.successExecutions.forEach(({ executionId }) => {
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
      this.successExecutions.forEach(({ executionId }) => {
        this.$socket
          .off(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
          .off(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
          .leave(this.getSocketRoomName(executionId))
          .removeListener(this.setExecutionStatus);
      });
    },

    async createExecutions() {
      this.pending = true;

      try {
        const items = await this.bulkCreateDeclareTicketExecution({ data: this.config.executions });
        const successExecutions = items.filter(({ status }) => status >= 200 && status < 300);
        const alarmsById = keyBy(this.config.alarms, '_id');
        const ticketsById = keyBy(this.config.tickets, '_id');

        this.successExecutions = successExecutions.map(({ id, item }) => ({
          executionId: id,
          ruleName: ticketsById[item._id].name,
          alarms: item.alarms.map(alarmId => alarmsById[alarmId]),
        }));
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    fetchExecutionsStatuses() {
      this.successExecutions.forEach(({ executionId }) => {
        this.fetchDeclareTicketExecutionWithoutStore({ id: executionId })
          .then(this.setExecutionStatus);
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.declare-ticket-rule-execute-status {
  display: flex;
  line-height: 24px !important;
}
</style>
