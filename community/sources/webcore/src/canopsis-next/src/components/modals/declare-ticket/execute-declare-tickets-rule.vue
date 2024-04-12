<template>
  <modal-wrapper
    close
    minimize
  >
    <template #title="">
      <v-layout align-center>
        <span>{{ title }}</span>
        <declare-ticket-rule-execution-status
          v-if="modal.minimized"
          :running="isAllExecutionsRunning"
          :success="isAllExecutionsSucceeded"
          :fail-reason="failReason"
          class="ml-2 declare-ticket-rule-execute-status"
          color="white"
        />
      </v-layout>
    </template>
    <template #text="">
      <v-layout
        v-if="pending"
        justify-center
      >
        <v-progress-circular
          color="primary"
          indeterminate
        />
      </v-layout>
      <template v-else-if="config.singleMode">
        <v-layout
          class="declare-ticket-rule-execute-status__executions"
          column
        >
          <declare-ticket-rule-executions-group
            v-for="(executions, ruleName) of alarmExecutionsByTicketName"
            :key="ruleName"
            :executions="executions"
            :rule-name="ruleName"
            is-one-execution
            show-status
            show-rule-name
          />
        </v-layout>
      </template>
      <template v-else>
        <declare-ticket-rule-executions-group
          :executions="alarmExecutions"
          :is-one-execution="isOneTicket"
          :show-status="isOneTicket"
        />
      </template>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { groupBy, keyBy } from 'lodash';

import { SOCKET_ROOMS } from '@/config';
import { WEBHOOK_EXECUTION_STATUSES, MODALS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import {
  isWebhookExecutionFinished,
  isWebhookExecutionRunning,
  isWebhookExecutionSucceeded,
} from '@/helpers/entities/webhook-execution/entity';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import DeclareTicketRuleExecutionsGroup from '@/components/other/declare-ticket/partials/declare-ticket-rule-executions-group.vue';
import DeclareTicketRuleExecutionStatus from '@/components/other/alarm/partials/alarm-test-query-execution-status.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to execute declare tickets
 */
export default {
  name: MODALS.executeDeclareTickets,
  components: {
    DeclareTicketRuleExecutionStatus,
    DeclareTicketRuleExecutionsGroup,
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
            status: WEBHOOK_EXECUTION_STATUSES.running,
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

    alarmExecutionsByTicketName() {
      return groupBy(this.alarmExecutions, 'ruleName');
    },

    isAllExecutionsRunning() {
      return this.alarmExecutions.some(isWebhookExecutionRunning);
    },

    isAllExecutionsSucceeded() {
      return this.alarmExecutions.every(isWebhookExecutionSucceeded);
    },

    isAllExecutionsFinished() {
      return this.alarmExecutions.every(this.isExecutionFinished);
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
      if (value.length && this.isAllExecutionsFinished) {
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

    isExecutionFinished(execution) {
      return isWebhookExecutionFinished(execution);
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
        const execution = this.executionsStatusesById[executionId] ?? {};

        if (this.isExecutionFinished(execution)) {
          return;
        }

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
      return Promise.all(
        this.successExecutions.map(
          ({ executionId }) => this.fetchDeclareTicketExecutionWithoutStore({ id: executionId })
            .then(this.setExecutionStatus),
        ),
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.declare-ticket-rule-execute-status {
  display: flex;
  line-height: 24px !important;

  &__executions {
    gap: 24px;
  }
}
</style>
