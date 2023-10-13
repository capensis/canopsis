<template lang="pug">
  v-layout(column)
    v-data-table(:headers="headers", :items="alarms", hide-actions)
      template(#items="{ item, index }")
        td(v-if="!hideRowSelect")
          v-checkbox(
            :input-value="isEveryTicketsActive(item._id)",
            color="primary",
            hide-details,
            @change="updateAllTickets(item._id, $event)"
          )
        td.text-xs-left {{ item.v.connector_name }}
        td.text-xs-left {{ item.v.connector }}
        td.text-xs-left {{ item.v.component }}
        td.text-xs-left {{ item.v.resource }}
        td(v-if="!hideTickets")
          declare-ticket-event-tickets-field(
            :value="activeTicketsByAlarms[item._id]",
            :tickets="ticketsByAlarms[item._id]",
            @input="updateTickets(item._id, $event)"
          )
    v-divider
</template>

<script>
import { difference } from 'lodash';

import { filterValue, mapIds } from '@/helpers/array';
import { revertGroupBy } from '@/helpers/collection';

import { formMixin } from '@/mixins/form';

import DeclareTicketEventTicketsField from './declare-ticket-event-tickets-field.vue';

export default {
  components: { DeclareTicketEventTicketsField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    ticketsByAlarms: {
      type: Object,
      default: () => ({}),
    },
    alarms: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'alarms_by_tickets',
    },
    hideTickets: {
      type: Boolean,
      default: false,
    },
    hideRowSelect: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    activeTicketsByAlarms() {
      return revertGroupBy(this.value);
    },

    headers() {
      return [
        !this.hideRowSelect && {
          sortable: false,
          width: 80,
        },
        {
          text: this.$t('common.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('common.connector'),
          sortable: false,
        },
        {
          text: this.$t('common.component'),
          sortable: false,
        },
        {
          text: this.$t('common.resource'),
          sortable: false,
        },
        !this.hideTickets && {
          text: this.$tc('common.ticket', 2),
          sortable: false,
        },
      ].filter(Boolean);
    },
  },
  methods: {
    isEveryTicketsActive(alarmId) {
      return this.activeTicketsByAlarms[alarmId]?.length === this.ticketsByAlarms[alarmId]?.length;
    },

    updateTickets(alarmId, tickets) {
      const oldTickets = this.activeTicketsByAlarms[alarmId] ?? [];

      const removedTickets = difference(oldTickets, tickets);
      const addedTickets = difference(tickets, oldTickets);

      const newValue = { ...this.value };

      addedTickets.forEach((ticketId) => {
        newValue[ticketId] = [...this.value[ticketId], alarmId];
      });

      removedTickets.forEach((ticketId) => {
        newValue[ticketId] = filterValue(this.value[ticketId], alarmId);
      });

      this.updateModel(newValue);
    },

    updateAllTickets(alarmId, checked) {
      this.updateTickets(alarmId, checked ? mapIds(this.ticketsByAlarms[alarmId]) : []);
    },
  },
};
</script>
