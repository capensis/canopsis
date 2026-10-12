<template>
  <c-layout column>
    <declare-ticket-event-alarms-tickets-field
      :value="form.alarms_by_tickets"
      :alarms="alarms"
      :tickets-by-alarms="ticketsByAlarms"
      :hide-tickets="hideTickets"
      :hide-row-select="hideRowSelect"
      @update:value="updateField('alarms_by_tickets', $event)"
    />
    <c-checkbox
      v-if="!hideTicketResource"
      :model-value="form.ticket_resources"
      :label="$t('declareTicket.ticketResources')"
      color="primary"
      @update:model-value="updateField('ticket_resources', $event)"
    />
    <c-description-field
      :model-value="form.comment"
      :label="$t('common.comment')"
      :autofocus="autofocus"
      class="mt-3"
      name="comment"
      @update:model-value="updateField('comment', $event)"
    />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import DeclareTicketEventAlarmsTicketsField from './fields/declare-ticket-event-alarms-tickets-field.vue';

export default {
  components: { DeclareTicketEventAlarmsTicketsField },
  mixins: [formMixin],
  emits: ['update:form'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    ticketsByAlarms: {
      type: Object,
      required: () => ({}),
    },
    alarms: {
      type: Array,
      default: () => [],
    },
    hideTickets: {
      type: Boolean,
      default: false,
    },
    hideRowSelect: {
      type: Boolean,
      default: false,
    },
    hideTicketResource: {
      type: Boolean,
      default: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
