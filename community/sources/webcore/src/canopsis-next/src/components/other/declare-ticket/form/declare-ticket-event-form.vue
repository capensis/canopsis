<template>
  <v-layout column>
    <declare-ticket-event-alarms-tickets-field
      v-field="form.alarms_by_tickets"
      :alarms="alarms"
      :tickets-by-alarms="ticketsByAlarms"
      :hide-tickets="hideTickets"
      :hide-row-select="hideRowSelect"
    />
    <v-checkbox
      v-if="!hideTicketResource"
      v-field="form.ticket_resources"
      :label="$t('declareTicket.ticketResources')"
      color="primary"
    />
    <c-description-field
      class="mt-3"
      v-field="form.comment"
      :label="$tc('common.comment')"
      name="comment"
    />
  </v-layout>
</template>

<script>
import DeclareTicketEventAlarmsTicketsField from './fields/declare-ticket-event-alarms-tickets-field.vue';

export default {
  components: { DeclareTicketEventAlarmsTicketsField },
  model: {
    prop: 'form',
    event: 'input',
  },
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
  },
};
</script>
