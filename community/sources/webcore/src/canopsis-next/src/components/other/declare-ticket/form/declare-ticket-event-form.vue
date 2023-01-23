<template lang="pug">
  v-layout(column)
    declare-ticket-event-alarms-tickets-field(
      v-field="form.tickets_by_alarms",
      :alarms="alarms",
      :tickets-by-alarms="ticketsByAlarms"
    )
    c-description-field.mt-3(v-field="form.comment", :label="$tc('common.comment')", name="comment")
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
  },
};
</script>
